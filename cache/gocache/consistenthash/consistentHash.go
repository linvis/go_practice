package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func([]byte) uint32

type Hash struct {
	hashFunc HashFunc
	replicas int
	keys     []int
	dict     map[int]string
}

func New(replicas int, hashFunc HashFunc) *Hash {
	h := &Hash{
		hashFunc: hashFunc,
		replicas: replicas,
		keys:     []int{},
		dict:     make(map[int]string),
	}

	if hashFunc == nil {
		h.hashFunc = crc32.ChecksumIEEE
	}

	return h
}

func (h *Hash) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < h.replicas; i++ {
			v := int(h.hashFunc([]byte(strconv.Itoa(i) + key)))
			h.keys = append(h.keys, v)
			h.dict[v] = key
		}
	}

	sort.Ints(h.keys)
}

func (h *Hash) Get(key string) string {
	hash := int(h.hashFunc([]byte(key)))

	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	})

	return h.dict[h.keys[idx%len(h.keys)]]
}
