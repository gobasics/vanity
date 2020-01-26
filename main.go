package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"gobasics.dev/env"
	"gobasics.dev/log"
)

const (
	TPL = `
<html>
	<head>
		<title>{{.Src}}</title>
		<meta name="go-import" content="{{.Src}} git https://github.com/{{.Dst}}"/>
	</head>
	<body>{{.Time}}</body>
</html>`
)

type handler struct {
	template *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Type", "charset=utf-8")

	type Redirect struct {
		Src, Dst string
		Time     time.Time
	}

	var path = r.URL.Path
	if len(path) > 0 && string(path[0]) == "/" {
		path = string(path[1:])
	}

	var data = Redirect{
		Src:  fmt.Sprintf("%s/%s", r.Host, path),
		Dst:  fmt.Sprintf("%s/%s.git", os.Getenv("GITHUB_HANDLE"), path),
		Time: time.Now(),
	}

	if err := h.template.Execute(w, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func run() error {
	t, err := template.New("webpage").Parse(TPL)
	if err != nil {
		return err
	}

	port, err := env.Get("PORT").Int()
	if err != nil {
		port = 8080
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler{t})
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
