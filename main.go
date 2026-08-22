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
			http.NotFound(w, req)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Println(err.Error())
		return
	}
	if stat.IsDir() {
		if err := handler.listIndex(w, req, thePath); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
		}
	} else if strings.HasSuffix(thePath, ".md") {
		if err := catAsMarkdown(thePath, w, req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
		}
	} else {
		http.ServeFile(w, req, thePath)
	}
}

func (handler *Handler) listIndex(w http.ResponseWriter, req *http.Request, dir string) error {
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
		path := filepath.ToSlash(filepath.Join(dir, name))
		fmt.Fprintf(w, "<li><a href=\"%s\">%s",
			url.PathEscape(path),
			html.EscapeString(name))
		if entry.IsDir() {
			w.Write([]byte{'/'})
		}
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
