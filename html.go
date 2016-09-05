package main

//go:generate go-bindata-assetfs static/... templates/...

/*
 * html.go
 * Serve up webpages
 * By MagisterQuis
 * Created 20160822
 * Last Modified 20160822
 */

import (
	"bytes"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"strings"
)

var templates = template.New("")

func init() {
	for _, path := range AssetNames() {
		if !strings.HasSuffix(path, ".tmpl") {
			continue
		}
		bytes, err := Asset(path)
		if err != nil {
			log.Fatalf(
				"Unable to parse template %v: %v",
				path,
				err,
			)
		}
		templates.New(path).Parse(string(bytes))
	}
}

/* ServeHTML serves up info via HTTP(s).  If http is true,addr is taken to be a
port, and http requests will be services on the loopback interface on that
port. */
func ServeFCGI(addr string, serveHTTP bool) {
	/* Register handlers */
	/* TODO: Finish this */
	http.Handle("/", http.FileServer(assetFS())) /* DEBUG */
	http.HandleFunc("/hello", hello)

	/* Listen and serve http */
	if serveHTTP {
		addr = "127.0.0.1:" + addr
		l, err := net.Listen("tcp4", addr)
		if nil != err {
			log.Fatalf("Unable to listen on %v: %v", addr, err)
		}
		log.Printf("Listening for HTTP connections on %v", l.Addr())
		log.Fatalf(
			"Error serving HTTP on %v: %v",
			l.Addr(),
			http.Serve(l, nil),
		)
	}
	/* Make a unix domain socket */
	l, err := net.Listen("unix", addr)
	if nil != err {
		log.Fatalf(
			"Cannot listen on unix socket %v: %v",
			addr,
			err,
		)
	}
	defer l.Close()
	log.Printf("Listening on %v for FastCGI requests", l.Addr())
	log.Fatalf(
		"Error serving FastCGI on %v: %v",
		l.Addr(),
		fcgi.Serve(l, nil),
	)
}

/* renderTemplate renders the template tmpl and sends it to w, given a model
p */
func renderTemplateBG(w http.ResponseWriter, tmpl string, p interface{}) {
	go func() {
		b := &bytes.Buffer{}
		err := templates.ExecuteTemplate(b, tmpl, p)
		log.Printf("E: %v", err)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}
		w.Write(b.Bytes())
	}()
}

func hello(w http.ResponseWriter, r *http.Request) {
	renderTemplateBG(w, "templates/index.tmpl", map[string]string{"Title": "worky"})
}

/* TODO: auto-register templates */
