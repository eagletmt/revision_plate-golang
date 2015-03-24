package revision_plate

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

type Handler struct {
	filePath string
	revision []byte
}

func New(filePath string) *Handler {
	h := &Handler{filePath: filePath}
	h.readCurrentRevision()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	switch r.Method {
	case "GET", "HEAD":
		rev, err := h.getCurrentRevision()
		if err == nil {
			w.WriteHeader(http.StatusOK)
			if r.Method == "GET" {
				w.Write(rev)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			if r.Method == "GET" {
				w.Write([]byte(err.Error()))
			}
		}
	}
}

func (h *Handler) readCurrentRevision() {
	file, err := os.Open(h.revisionFilePath())
	if err != nil {
		return
	}
	defer file.Close()
	revision, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	h.revision = revision
}

func (h *Handler) getCurrentRevision() ([]byte, error) {
	if h.revision == nil {
		return nil, errors.New("REVISION_FILE_NOT_FOUND")
	} else {
		_, err := os.Stat(h.revisionFilePath())
		if err == nil {
			return h.revision, nil
		} else {
			return nil, errors.New("REVISION_FILE_REMOVED")
		}
	}
}

func (h *Handler) revisionFilePath() string {
	if h.filePath == "" {
		return "REVISION"
	} else {
		return h.filePath
	}
}
