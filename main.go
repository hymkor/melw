package main

import (
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var jumpTable = map[string]func(h *Handler, w http.ResponseWriter, req *http.Request) error{
	"Edit":    (*Handler).actionEdit,
	"Preview": (*Handler).actionPreview,
	"Save":    (*Handler).actionSave,
	"New":     (*Handler).actionNew,
	"Cancel":  (*Handler).actionCancel,
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s \"%s\"\n", req.RemoteAddr, req.Method, req.URL.Path)

	err := handler.serveHTTP(w, req)
	if err != nil {
		log.Println(err.Error())
		if os.IsNotExist(err) {
			http.NotFound(w, req)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *Handler) serveHTTP(w http.ResponseWriter, req *http.Request) error {
	if f, ok := jumpTable[req.FormValue("a")]; ok {
		return f(h, w, req)
	}
	stat, err := h.Stat(req.URL.Path)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return h.listIndex(w, req)
	} else if strings.HasSuffix(req.URL.Path, ".md") {
		return h.catAsMarkdown(w, req)
	} else {
		return h.catFile(w, req)
	}
}

func (h *Handler) catFile(w http.ResponseWriter, req *http.Request) error {
	fd, err := h.Open(req.URL.Path)
	if err != nil {
		return err
	}
	defer fd.Close()
	io.Copy(w, fd)
	return nil
}

func (h *Handler) listIndex(w http.ResponseWriter, req *http.Request) error {
	if !strings.HasSuffix(req.URL.Path, "/") {
		http.Redirect(w, req, req.URL.Path+"/", http.StatusMovedPermanently)
		return nil
	}
	files, err := h.ReadDir(req.URL.Path)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	dir_ := html.EscapeString(req.URL.Path)
	fmt.Fprintf(w, "<html><head><title>Index: %s/</title></head>\n", dir_)
	fmt.Fprintf(w, "<body><h1>Index: %s</h1>\n", dir_)
	io.WriteString(w, "<ul>\n")
	for _, entry := range files {
		name := entry.Name()
		slash := ""
		if entry.IsDir() {
			slash = "/"
		}
		fmt.Fprintf(w, "<li><a href=\"%s%s\">%s%s",
			url.PathEscape(name),
			slash,
			html.EscapeString(name),
			slash)
		io.WriteString(w, "</a></li>\n")
	}
	io.WriteString(w, "</ul>\n")
	fmt.Fprintf(w, "<form action=\"%s\" method=\"POST\">\n", dir_)
	fmt.Fprintf(w, "<input type=\"hidden\" name=\"dir\" value=\"%s\" />\n", dir_)
	io.WriteString(w, "<input type=\"text\" name=\"p\" /><tt>.md</tt>\n")
	io.WriteString(w, "<input type=\"submit\" name=\"a\" value=\"New\" />\n")
	io.WriteString(w, "</form>\n")
	io.WriteString(w, "</body></html>\n")
	return nil
}

var (
	flagC = flag.String("C", "", "Change working directory")
	flagP = flag.Uint("P", 8000, "Port")
)

func mains() error {
	if *flagC != "" {
		if err := os.Chdir(*flagC); err != nil {
			return err
		}
	}
	root, err := os.OpenRoot(".")
	if err != nil {
		return err
	}
	handler := &Handler{root: root}
	service := &http.Server{
		Addr:           fmt.Sprintf(":%d", *flagP),
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = service.ListenAndServe()
	closeErr := service.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func main() {
	flag.Parse()
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
