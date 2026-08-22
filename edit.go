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
)

func drawForm(w http.ResponseWriter, req *http.Request, source string) {
	path_ := html.EscapeString(path.Clean(req.URL.Path))
	fmt.Fprintf(w, "<h1>Edit: %s</h1>\n", path_)
	fmt.Fprintf(w, "<form action=\"%s\" method=\"POST\">\n", path_)
	fmt.Fprintf(w, "<textarea name=\"text\" style=\"width:100%%\" cols=\"80\" rows=\"20\">%s</textarea>\n",
		html.EscapeString(source))
	fmt.Fprintf(w, "<input type=\"submit\" name=\"a\" value=\"Preview\" />\n")
	fmt.Fprintf(w, "<input type=\"submit\" name=\"a\" value=\"Save\" />\n")
	fmt.Fprintf(w, "<input type=\"submit\" name=\"a\" value=\"Cancel\" style=\"float:right\" />\n")
	io.WriteString(w, "</form>\n")

}

func actionEdit(w http.ResponseWriter, req *http.Request) error {
	thePath, err := url.QueryUnescape(req.URL.Path)
	thePath = filepath.Join(".", filepath.FromSlash(thePath))
	source, err := os.ReadFile(thePath)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, htmlHeader, gitHubCss)
	drawForm(w, req, string(source))
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
	drawForm(w, req, source)
	htmls.WriteTo(w)

	fmt.Fprintln(w, htmlFooter)
	return nil
}

func actionSave(w http.ResponseWriter, req *http.Request) error {
	source := req.FormValue("text")
	thePath, err := url.QueryUnescape(req.URL.Path)
	thePath = filepath.Join(".", filepath.FromSlash(thePath))
	err = os.WriteFile(thePath, []byte(source), 0644)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)

	newUrl := html.EscapeString(path.Clean(req.URL.Path))

	fmt.Fprint(w, "<html><head>\n")
	fmt.Fprintf(w, "<meta http-equiv=\"refresh\" content=\"1;URL=%s\">\n", newUrl)
	fmt.Fprint(w, "</head><body>\n")
	fmt.Fprintf(w, "<p><a href=\"%s\">Wait or Click Here</a></p>\n", newUrl)
	fmt.Fprintf(w, "</body></html>\n")
	return nil
}
