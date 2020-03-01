package gocache

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"strings"
)

const basePath = "/gocache"

type HTTPServer struct {
	basePath string
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.URL.Path, s.basePath) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// /basePath/groupName/key
	pathes := strings.SplitN(r.URL.Path, "/", -1)
	if len(pathes) < 3 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	g := GetGroup(pathes[2])
	if g == nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	v, err := g.Get(pathes[3])
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// b, _ := getBytes("say hello")

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(v.(string)))
}

func NewServer() *HTTPServer {
	return &HTTPServer{basePath}

}

func (s *HTTPServer) Run(path string) {
	http.ListenAndServe(path, s)
}
