package gocache

import pb "gocache/gocachepb"

type PeerPicker interface {
	PeerPick(key string) (PeerGetter, error)
}

type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
