package main

import (
	"flag"
	"fmt"
	"gocache"
	"net/http"
)

var testMap = map[string]string{
	"aaa": "111",
	"bbb": "222",
	"ccc": "333",
}

func createGroup() *gocache.Group {
	return gocache.NewGroup("test", 200, func(key string) ([]byte, error) {
		fmt.Println("hit")
		return []byte(testMap[key]), nil
	})
}

func startServer(addr string, peers []string, group *gocache.Group) {
	pool := gocache.NewHTTPPool()
	pool.Set(peers...)
	group.RegisterPeer(pool)
	pool.Run(addr)
}

func startAPIServer(apiAddr string, group *gocache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			fmt.Println(key)
			view, err := group.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))

	fmt.Println("fontend server is running at", apiAddr)
	http.ListenAndServe(apiAddr[7:], nil)
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	fmt.Println(port, api)

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	g := createGroup()

	if api {
		go startAPIServer(apiAddr, g)
	}
	startServer(addrMap[port], []string(addrs), g)
}
