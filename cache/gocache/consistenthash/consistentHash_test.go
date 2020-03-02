package consistenthash

import (
	"strconv"
	"testing"
)

func calHash(t *testing.T, h *ConsistentHash, val string, expect string) {
	res := h.Get(val)

	if res != expect {
		t.Errorf("error: res %s, expect %s", res, expect)
	}
}

func TestConsistentHash(t *testing.T) {
	h := New(2, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	//2 3  5 6  8  9
	h.Add("2", "5", "9")

	calHash(t, h, "2", "2")
	calHash(t, h, "4", "5")
	calHash(t, h, "7", "9")
	calHash(t, h, "9", "9")

	h.Add("7")
	calHash(t, h, "7", "7")

}
