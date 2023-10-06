package geecache

// PeerPicker 选择节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 节点需要实现的一个获取值的一个方法
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
