package gocache

type PeerPicker interface {
	PeerPick(key string) (PeerGetter, error)
}

type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
