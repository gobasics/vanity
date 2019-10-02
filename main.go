package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"gobasics.dev/env"
	"gobasics.dev/log"
)

const (
	TPL = `
<html>
	<head>
		<title>{{.Src}}</title>
		<meta name="go-import" content="{{.Src}} {{.Dst}}"/>
	</head>
	<body></body>
</html>`
)

type handler struct {
	template *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Type", "charset=utf-8")

	type R struct {
		Src, Dst string
	}

	var data = R{
		Src: os.Getenv("SRC"),
		Dst: os.Getenv("DST"),
	}

	if err := h.template.Execute(os.Stdout, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func run() error {
	t, err := template.New("webpage").Parse(TPL)
	if err != nil {
		return err
	}

	port, err := env.Get("PORT").Int()
	if err != nil {
		port = 80
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler{t})
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
