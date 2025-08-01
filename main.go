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

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
	"github.com/sploders101/resume-server/migrations"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

func main() {
	resumePath, exists := os.LookupEnv("RESUME_PATH")
	if !exists {
		log.Fatalln("RESUME_PATH unspecified.")
	}

	dbPath, exists := os.LookupEnv("DB_PATH")
	if !exists {
		log.Fatalln("DB_PATH unspecified")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	migrations.InitDb(db)

	router := mux.NewRouter()

	router.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/static/style.css")
	})
	router.HandleFunc("/", htmlResumeHandler(resumePath, "./assets/resume.html"))
	router.HandleFunc("/pdf", func(w http.ResponseWriter, r *http.Request) {
		result, err := pdfGrabber("http://127.0.0.1:8080/")
		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(500)
			w.Write([]byte("An internal server error occurred."))
			return
		}
		w.Header().Add("Content-Type", "application/pdf")
		w.Header().Add("Content-Disposition", `attachment; filename="resume.pdf"`)
		w.Write(result)
	})

	log.Println("Listening for connections...")
	http.ListenAndServe(":8080", router)
}

// Gets the resume data from a yaml file, given a path
func getResumeData(resumePath string) (Resume, error) {
	file, err := os.ReadFile(resumePath)
	if err != nil {
		return Resume{}, err
	}

	var resume Resume
	err = yaml.Unmarshal(file, &resume)
	if err != nil {
		return Resume{}, err
	}

	return resume, nil
}

// Gets the resume template, given a path
func getResumeTemplate(templatePath string) (*template.Template, error) {
	tmpl, err := template.New("resume.html").Funcs(
		template.FuncMap{
			"md":       markdown,
			"initials": initials,
		},
	).ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// Generates a handler function for serving resumes via HTTP, given a resume and template path.
func htmlResumeHandler(resumePath string, templatePath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resume, err := getResumeData(resumePath)
		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(500)
			w.Write([]byte("An error occurred while reading the resume.\n"))
			return
		}

		template, err := getResumeTemplate(templatePath)
		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(500)
			w.Write([]byte("An error occurred while reading the resume template.\n"))
			return
		}

		w.Header().Add("Content-Type", "text/html")
		template.Execute(w, resume)
	}
}

// Generates initials from a name. Intended for use in the template.
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

// Generates HTML from a markdown string. Intended for use in the template.
func markdown(markdown string) any {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return "[ERROR]" + err.Error()
	}
	return template.HTML(buf.String())
}

// Creates a PDF from the given URL.
// The intended use is for generating PDFs from a resume template.
func pdfGrabber(url string) ([]byte, error) {
	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()
	var pdfBuffer []byte
	tasks := chromedp.Tasks{
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
