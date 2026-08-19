package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Handler struct {
}

var jumpTable = map[string]func(w http.ResponseWriter, req *http.Request) error{
	"Edit":    actionEdit,
	"Preview": actionPreview,
	"Save":    actionSave,
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s \"%s\"\n", req.RemoteAddr, req.Method, req.URL.Path)

	if f, ok := jumpTable[req.FormValue("a")]; ok {
		if err := f(w, req); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	thePath, err := url.QueryUnescape(req.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
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
		log.Println(err.Error())
		return
	}
	if stat.IsDir() {
		if err := handler.listIndex(w, req, thePath); err != nil {
			log.Println(err.Error())
		}
	} else if strings.HasSuffix(thePath, ".md") {
		if err := catAsMarkdown(thePath, w, req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
		}
	} else {
		if err := handler.catFile(w, req, thePath); err != nil {
			log.Println(err.Error())
		}
	}
}

func (handler *Handler) catFile(w http.ResponseWriter, req *http.Request, thePath string) error {
	fd, err := os.Open(thePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	defer fd.Close()
	io.Copy(w, fd)
	return nil
}

func (handler *Handler) listIndex(w http.ResponseWriter, req *http.Request, dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
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
	return nil
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
