# Resume Server

I'll come back and flesh out this description later, but this is essentially a server made to render my resume from a YAML file that can be injected using a Kubernetes ConfigMap.

This project is my demo project for learning golang, chosen for its practical use, and coverage of API routing, database usage, template rendering, and 3rd-party module usage. It is very simple in concept, but uses a little bit of everything, which makes it a good candidate for a learning project.

I got tired of copying HTML/CSS back and forth, and not being able to obtain PDFs from my phone (mobile browsers seem to format the PDF different from a PC for some reason).

This server can render out an HTML template, or use that template with a headless Dockerized Chrome instance to get a PDF.
