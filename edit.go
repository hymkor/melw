package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"path"
	"strings"
)

func drawForm(w http.ResponseWriter, req *http.Request, netPath, source string) error {
	path_ := html.EscapeString(netPath)
	fmt.Fprintf(w, "<h1>Edit: %s</h1>\n", path_)
	fmt.Fprintf(w, "<form action=\"%s\" method=\"POST\">\n", path_)
	fmt.Fprintf(w, "<textarea name=\"text\" style=\"width:100%%\" cols=\"80\" rows=\"20\" hx-post=\"?a=RawPreview\" hx-trigger=\"input changed delay:500ms\" hx-target=\"#preview\" hx-swap=\"innerHTML\">%s</textarea>\n",
		html.EscapeString(source))
	io.WriteString(w, "<input type=\"submit\" name=\"a\" value=\"Save\" />\n")
	io.WriteString(w, "<input type=\"submit\" name=\"a\" value=\"Cancel\" style=\"float:right\" />\n")
	io.WriteString(w, "</form>\n")
	htmls, err := markdownToHtml([]byte(source))
	if err != nil {
		return err
	}
	io.WriteString(w, "<div id=\"preview\"></div>\n")
	htmls.WriteTo(w)
	io.WriteString(w, "</div>\n")
	return nil
}

func (h *Handler) actionNew(w http.ResponseWriter, req *http.Request) error {
	dir := req.FormValue("dir")
	page := req.FormValue("p") + ".md"
	netPath := path.Join(dir, page)
	return h.doEdit(w, req, netPath)
}

func (h *Handler) actionEdit(w http.ResponseWriter, req *http.Request) error {
	return h.doEdit(w, req, req.URL.Path)
}

func showOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
}

func (h *Handler) doEdit(w http.ResponseWriter, req *http.Request, netPath string) error {
	source, err := h.ReadFile(netPath)
	if err != nil {
		source = []byte{}
	}
	showOK(w)
	fmt.Fprintf(w, htmlHeader, gitHubCss)
	err = drawForm(w, req, netPath, string(source))
	fmt.Fprintln(w, htmlFooter)
	return err
}

func (h *Handler) actionRawPreview(w http.ResponseWriter, req *http.Request) error {
	source := req.FormValue("text")
	htmls, err := markdownToHtml([]byte(source))
	if err != nil {
		return err
	}
	showOK(w)
	htmls.WriteTo(w)
	return nil
}

func transfer(w http.ResponseWriter, newUrl string) {
	newUrl = html.EscapeString(newUrl)

	showOK(w)

	io.WriteString(w, "<html><head>\n")
	fmt.Fprintf(w, "<meta http-equiv=\"refresh\" content=\"1;URL=%s\">\n", newUrl)
	io.WriteString(w, "</head><body>\n")
	fmt.Fprintf(w, "<p><a href=\"%s\">Wait or Click Here</a></p>\n", newUrl)
	io.WriteString(w, "</body></html>\n")
}

func (h *Handler) actionSave(w http.ResponseWriter, req *http.Request) error {
	source := req.FormValue("text")
	newUrl := req.URL.Path
	var err error
	if strings.TrimSpace(source) == "" {
		err = h.Remove(req.URL.Path)
		newUrl = path.Dir(newUrl)
	} else {
		err = h.WriteFile(req.URL.Path, []byte(source), 0644)
	}
	if err != nil {
		return err
	}
	transfer(w, newUrl)
	return nil
}

func (h *Handler) actionCancel(w http.ResponseWriter, req *http.Request) error {
	_, err := h.Stat(req.URL.Path)
	if err != nil {
		transfer(w, path.Dir(req.URL.Path))
		return nil
	}
	return h.serveMarkdown(w, req)
}
