package main

import (
	"context"
	"errors"
	"html/template"
	"net"
	"net/http"
	"os"

	"gobasics.dev/app"
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
</html>
	`
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

type provider struct {
	*http.Server
}

func (p provider) Serve(l net.Listener) error {
	return p.Serve(l)
}

func (p provider) GracefulStop() {
	if err := p.Shutdown(context.Background()); err != nil {
		log.Error(err.Error())
	}
}

func run() error {
	dirCache, ok := env.Get("DIR_CACHE")
	if !ok {
		return errors.New("DIR_CACHE is not set in the environment")
	}

	hostnames, ok := env.Get("HOSTNAMES")
	if !ok {
		return errors.New("HOSTNAMES is not set in the environment")
	}

	t, err := template.New("webpage").Parse(TPL)
	if err != nil {
		return err
	}

	app := &app.Server{
		Config: app.Config{
			Letsencrypt: true,
			DirCache:    dirCache.String(),
			HostNames:   hostnames.StringSlice(","),
			Port:        443,
		},
		Provider: provider{&http.Server{
			Handler: handler{t},
		}},
	}

	return app.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
