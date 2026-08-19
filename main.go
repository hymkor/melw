package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Handler struct {
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	thePath, err := url.QueryUnescape(req.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	thePath = filepath.Join(".", filepath.FromSlash(thePath))

	stat, err := os.Stat(thePath)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if stat.IsDir() {
		handler.listIndex(w, req, thePath)
	} else if strings.HasSuffix(thePath, ".md") {
		if err := catAsMarkdown(thePath, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		handler.catFile(w, req, thePath)
	}
}

func (handler *Handler) catFile(w http.ResponseWriter, req *http.Request, thePath string) {
	fd, err := os.Open(thePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer fd.Close()
	io.Copy(w, fd)
}

func (handler *Handler) listIndex(w http.ResponseWriter, req *http.Request, dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "<html><body><h1>Markdown Editor Like Wiki</h1><ul>\n")
	defer io.WriteString(w, "</ul></body></html>\n")

	for _, entry := range files {
		name := entry.Name()
		path := filepath.ToSlash(filepath.Join(dir, name))
		fmt.Fprintf(w, "<li><a href=\"%s\">%s",
			url.PathEscape(path),
			name)
		if entry.IsDir() {
			w.Write([]byte{'/'})
		}
		io.WriteString(w, "</a></li>\n")
	}
}

func mains() error {
	handler := new(Handler)
	service := &http.Server{
		Addr:           ":8000",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := service.ListenAndServe()
	closeErr := service.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
