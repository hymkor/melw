package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func drawForm(w http.ResponseWriter, req *http.Request, thePath, source string) {
	path_ := html.EscapeString(path.Clean(thePath))
	fmt.Fprintf(w, "<h1>Edit: %s</h1>\n", path_)
	fmt.Fprintf(w, "<form action=\"%s\" method=\"POST\">\n", path_)
	fmt.Fprintf(w, "<textarea name=\"text\" style=\"width:100%%\" cols=\"80\" rows=\"20\">%s</textarea>\n",
		html.EscapeString(source))
	fmt.Fprintf(w, "<input type=\"submit\" name=\"a\" value=\"Preview\" />\n")
	fmt.Fprintf(w, "<input type=\"submit\" name=\"a\" value=\"Save\" />\n")
	fmt.Fprintf(w, "<input type=\"submit\" name=\"a\" value=\"Cancel\" style=\"float:right\" />\n")
	io.WriteString(w, "</form>\n")

}

func actionNew(w http.ResponseWriter, req *http.Request) error {
	dir := req.FormValue("dir")
	page := req.FormValue("p") + ".md"
	thePath := filepath.Join(dir, page)
	return doEdit(w, req, thePath)
}

func actionEdit(w http.ResponseWriter, req *http.Request) error {
	thePath, err := url.QueryUnescape(req.URL.Path)
	if err != nil {
		return err
	}
	thePath = filepath.Join(".", filepath.FromSlash(thePath))
	return doEdit(w, req, thePath)
}

func doEdit(w http.ResponseWriter, req *http.Request, thePath string) error {
	source, err := os.ReadFile(thePath)
	if err != nil {
		source = []byte{}
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, htmlHeader, gitHubCss)
	drawForm(w, req, thePath, string(source))
	fmt.Fprintln(w, htmlFooter)
	return nil
}

func actionPreview(w http.ResponseWriter, req *http.Request) error {
	source := req.FormValue("text")
	htmls, err := markdownToHtml([]byte(source))
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, htmlHeader, gitHubCss)
	drawForm(w, req, req.URL.Path, source)
	htmls.WriteTo(w)

	fmt.Fprintln(w, htmlFooter)
	return nil
}

func transfer(w io.Writer, newUrl string) {
	newUrl = html.EscapeString(path.Clean(newUrl))

	fmt.Fprint(w, "<html><head>\n")
	fmt.Fprintf(w, "<meta http-equiv=\"refresh\" content=\"1;URL=%s\">\n", newUrl)
	fmt.Fprint(w, "</head><body>\n")
	fmt.Fprintf(w, "<p><a href=\"%s\">Wait or Click Here</a></p>\n", newUrl)
	fmt.Fprintf(w, "</body></html>\n")
}

func actionSave(w http.ResponseWriter, req *http.Request) error {
	source := req.FormValue("text")
	thePath, err := url.QueryUnescape(req.URL.Path)
	thePath = filepath.Join(".", filepath.FromSlash(thePath))
	newUrl := req.URL.Path
	if strings.TrimSpace(source) == "" {
		err = os.Remove(thePath)
		newUrl = path.Dir(newUrl)
	} else {
		err = os.WriteFile(thePath, []byte(source), 0644)
	}
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	transfer(w, newUrl)
	return nil
}

func actionCancel(w http.ResponseWriter, req *http.Request) error {
	thePath, err := url.QueryUnescape(req.URL.Path)
	if err != nil {
		return err
	}
	thePath = filepath.Join(".", filepath.FromSlash(thePath))
	_, err = os.Stat(thePath)
	if err != nil {
		transfer(w, path.Dir(req.URL.Path))
		return nil
	}
	return catAsMarkdown(thePath, w, req)
}
