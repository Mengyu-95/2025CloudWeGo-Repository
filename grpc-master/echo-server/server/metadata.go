package server

import "google.golang.org/grpc/metadata"

func getMetadataByMap(mp map[string]string) metadata.MD {
	//通过map 初始化metadata
	md := metadata.New(mp)
	return md
}

func getMetadataByKV(kv ...string) metadata.MD {
	//通过键值对的方式初始化metadata
	md := metadata.Pairs(kv...)
	return md
}
