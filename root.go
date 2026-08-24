package main

import (
	"fmt"
	"html"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
)

type Handler struct {
	root *os.Root
}

func trimHeadRoot(path string) string {
	for {
		if len(path) <= 0 || path[0] != '/' {
			return path
		}
		path = path[1:]
	}
}

func trimTailRoot(path string) string {
	for {
		if len(path) <= 0 || path[len(path)-1] != '/' {
			return path
		}
		path = path[:len(path)-1]
	}
}

func normPath(path string) string {
	path = trimHeadRoot(trimTailRoot(path))

	if path == "" {
		return "."
	}
	return path
}

func (h *Handler) Open(name string) (*os.File, error) {
	return h.root.Open(normPath(name))
}

func (h *Handler) Stat(name string) (os.FileInfo, error) {
	return h.root.Stat(normPath(name))
}

func (h *Handler) ReadFile(name string) ([]byte, error) {
	return h.root.ReadFile(normPath(name))
}

func (h *Handler) Remove(name string) error {
	return h.root.Remove(normPath(name))
}

func (h *Handler) WriteFile(name string, bin []byte, perm os.FileMode) error {
	return h.root.WriteFile(normPath(name), bin, perm)
}

func (h *Handler) ReadDir(name string) ([]fs.DirEntry, error) {
	return h.root.FS().(fs.ReadDirFS).ReadDir(normPath(name))
}

func printNestPath(w http.ResponseWriter, netPath string) {
	if netPath == "" || netPath == "/" {
		io.WriteString(w, "<a href=\"/\">/ </a>")
		return
	}
	netPath = trimTailRoot(netPath)
	parent := path.Dir(netPath)
	child := path.Base(netPath)
	if parent == "" || parent == "/" {
		io.WriteString(w, "<a href=\"/\">/ </a>")
	} else {
		printNestPath(w, parent)
		io.WriteString(w, " / ")
	}

	url1 := url.URL{Path: netPath}
	fmt.Fprintf(w, "<a href=\"%s\">%s</a>",
		url1.EscapedPath(),
		html.EscapeString(child))
}
