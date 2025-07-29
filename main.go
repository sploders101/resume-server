package main

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"log"
	"strings"
	"time"

	"html/template"
	"net/http"
	"os"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
	"github.com/sploders101/resume-server/migrations"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

const DB_PATH = "./data.db"

func main() {
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		log.Fatal(err)
	}
	migrations.InitDb(db)

	resumePath, exists := os.LookupEnv("RESUME_PATH")
	if !exists {
		log.Fatalln("RESUME_PATH unspecified.")
	}

	router := mux.NewRouter()

	router.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/static/style.css")
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile(resumePath)
		if err != nil {
			log.Fatalln(err)
		}

		var resume Resume
		err = yaml.Unmarshal(file, &resume)
		if err != nil {
			log.Fatalln(err)
		}
		tmpl := template.Must(template.New("resume.html").Funcs(template.FuncMap{
			"md":       markdown,
			"initials": initials,
		}).ParseFiles("./assets/resume.html"))
		tmpl.Execute(w, resume)
	})
	router.HandleFunc("/pdf", func(w http.ResponseWriter, r *http.Request) {
		result, err := pdfGrabber("http://127.0.0.1:8080/")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(result)
	})

	log.Println("Listening for connections...")
	http.ListenAndServe(":8080", router)
}

func initials(name string) string {
	segments := strings.Fields(name)
	initials := ""
	for _, segment := range segments {
		if len(segment) == 0 {
			continue
		}
		for _, initial := range segment {
			initials += string(initial)
			break
		}
	}

	return initials
}

func markdown(markdown string) any {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return "[ERROR]" + err.Error()
	}
	return template.HTML(buf.String())
}

func pdfGrabber(url string) ([]byte, error) {
	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()
	var pdfBuffer []byte
	tasks := chromedp.Tasks{
		emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(url),
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.Sleep(100 * time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			printer := page.PrintToPDF()
			printer = printer.WithPrintBackground(true)
			printer = printer.WithMarginTop(0.4)
			printer = printer.WithMarginLeft(0.4)
			printer = printer.WithMarginRight(0.4)
			printer = printer.WithMarginBottom(0.4)
			pdfBuffer, _, err = printer.Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
	if err := chromedp.Run(taskCtx, tasks); err != nil {
		log.Fatal(err)
	}
	return pdfBuffer, nil
}
