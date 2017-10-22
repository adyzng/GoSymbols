package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	clog "gopkg.in/clog.v1"
)

// IndexHandle for index page
//
func IndexHandle(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err == nil {
		tmpl.Execute(w, nil)
		//w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
	}
}

// StaticHandler serve public files, exclude folder
//
func StaticHandler(folder string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if request folder, return not found
		if strings.TrimRight(r.RequestURI, "/") != r.RequestURI {
			http.NotFound(w, r)
		} else {
			http.FileServer(http.Dir(folder)).ServeHTTP(w, r)
		}
	})
}

// LogHandler print request trace log
//
func LogHandler(h http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)

		// "GET / HTTP/1.1" 200 2552 UserAgent
		clog.Info("%s - %s %s %s - %s",
			r.RemoteAddr,
			r.Proto,
			r.Method,
			r.RequestURI,
			time.Since(start))
	})
}
