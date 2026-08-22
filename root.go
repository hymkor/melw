package main

import (
	"io/fs"
	"os"
)

type Handler struct {
	root *os.Root
}

func normPath(path string) string {
	for {
		if len(path) <= 0 || path[0] != '/' {
			break
		}
		path = path[1:]
	}
	for {
		if len(path) <= 0 || path[len(path)-1] != '/' {
			break
		}
		path = path[:len(path)-1]
	}
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
