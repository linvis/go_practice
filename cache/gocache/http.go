package gocache

import (
	"errors"
	"fmt"
	"gocache/consistenthash"
	pb "gocache/gocachepb"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
)

const basePath = "/gocache"
const replicas = 50

type HTTPPool struct {
	basePath    string
	self        string
	mu          sync.Mutex
	peers       *consistenthash.Hash
	httpGetters map[string]*httpGetter
}

func (h *HTTPPool) Set(peers ...string) {
	h.mu.Lock()

	defer h.mu.Unlock()

	h.peers = consistenthash.New(replicas, nil)

	h.peers.Add(peers...)

	h.httpGetters = make(map[string]*httpGetter)

	for _, peer := range peers {
		h.httpGetters[peer] = &httpGetter{baseURL: peer + h.basePath}
	}
}

func (h *HTTPPool) PeerPick(key string) (PeerGetter, error) {
	peer := h.peers.Get(key)

	if peer == h.self {
		return nil, errors.New("request self")
	}

	getter := h.httpGetters[peer]

	return getter, nil
}

func (s *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	view, err := g.Get(pathes[3])
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func NewHTTPPool() *HTTPPool {
	return &HTTPPool{
		basePath: basePath,
	}

}

func (s *HTTPPool) Run(path string) {
	s.self = path
	http.ListenAndServe(path[7:], s)
}

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	url := fmt.Sprintf(
		"%v/%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)

	fmt.Println("request peer  ", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad request %d", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return err
	}

	return nil
}
