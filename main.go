package main

import (
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/browser"
)

var jumpTable = map[string]func(h *Handler, w http.ResponseWriter, req *http.Request) error{
	"Edit":       (*Handler).actionEdit,
	"Save":       (*Handler).actionSave,
	"New":        (*Handler).actionNew,
	"Cancel":     (*Handler).actionCancel,
	"RawPreview": (*Handler).actionRawPreview,
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
	} else if ext := path.Ext(req.URL.Path); strings.EqualFold(ext, ".md") || strings.EqualFold(ext, ".mkd") {
		return h.serveMarkdown(w, req)
	} else {
		return h.serveFile(w, req)
	}
}

func (h *Handler) serveFile(w http.ResponseWriter, req *http.Request) error {
	http.ServeFileFS(w, req, h.root.FS(), req.URL.Path)
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
	showOK(w)
	fmt.Fprintf(w, htmlHeader1, gitHubCss)
	dir_ := html.EscapeString(req.URL.Path)
	io.WriteString(w, "<h1>Index: ")
	printNestPath(w, req.URL.Path)
	io.WriteString(w, "</h1><ul>\n")
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
	fmt.Fprintln(w, htmlFooter)
	return nil
}

var (
	flagC     = flag.String("C", "", "Change working directory")
	flagP     = flag.Uint("P", 8000, "Port number")
	flagStart = flag.Bool("start", false, "Open web browser at startup")
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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *flagP))
	if err != nil {
		return err
	}
	if *flagStart {
		url := "http://" + listener.Addr().String() + "/"
		go func() {
			if err := browser.OpenURL(url); err != nil {
				log.Printf("cannot open browser: %v", err)
			}
		}()
	}
	handler := &Handler{root: root}
	service := &http.Server{
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err = service.Serve(listener)
	closeErr := service.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func main() {
	fmt.Fprintf(os.Stderr, "%s %s-%s-%s by %s\n",
		filepath.Base(os.Args[0]),
		version,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.Version())
	flag.Parse()
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
