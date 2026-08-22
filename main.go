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
	"New":     actionNew,
	"Cancel":  actionCancel,
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

func (handler *Handler) serveHTTP(w http.ResponseWriter, req *http.Request) error {
	if f, ok := jumpTable[req.FormValue("a")]; ok {
		return f(w, req)
	}

	thePath, err := url.QueryUnescape(req.URL.Path)
	if err != nil {
		return err
	}
	thePath = filepath.Join(".", filepath.FromSlash(thePath))

	stat, err := os.Stat(thePath)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return handler.listIndex(w, req, thePath)
	} else if strings.HasSuffix(thePath, ".md") {
		return catAsMarkdown(thePath, w, req)
	} else {
		http.ServeFile(w, req, thePath)
		return nil
	}
}

func (handler *Handler) listIndex(w http.ResponseWriter, req *http.Request, dir string) error {
	if !strings.HasSuffix(req.URL.Path, "/") {
		http.Redirect(w, req, req.URL.Path+"/", http.StatusMovedPermanently)
		return nil
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	dir_ := html.EscapeString(dir)
	fmt.Fprintf(w, "<html><head><title>Index: %s/</title></head>\n", dir_)
	fmt.Fprintf(w, "<body><h1>Index: %s/</h1>\n", dir_)
	fmt.Fprint(w, "<ul>\n")
	for _, entry := range files {
		name := entry.Name()
		slash := ""
		if entry.IsDir() {
			slash = "/"
		}
		// path := filepath.ToSlash(filepath.Join(dir, name))
		fmt.Fprintf(w, "<li><a href=\"%s%s\">%s%s",
			url.PathEscape(name),
			slash,
			html.EscapeString(name),
			slash)
		io.WriteString(w, "</a></li>\n")
	}
	fmt.Fprint(w, "</ul>\n")
	fmt.Fprintf(w, "<form action=\"%s\" method=\"POST\">\n", dir_)
	fmt.Fprintf(w, "<input type=\"hidden\" name=\"dir\" value=\"%s\" />\n", dir_)
	fmt.Fprint(w, "<input type=\"text\" name=\"p\" /><tt>.md</tt>\n")
	fmt.Fprint(w, "<input type=\"submit\" name=\"a\" value=\"New\" />\n")
	fmt.Fprint(w, "</form>\n")
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
	handler := new(Handler)
	service := &http.Server{
		Addr:           fmt.Sprintf(":%d", *flagP),
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
	flag.Parse()
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
