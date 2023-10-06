package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           //采用依赖注入的哈希函数
	replicas int            //虚拟节点的个数
	keys     []int          // Sorted 哈希环
	hashMap  map[int]string //虚拟节点于真实节点之间的映射
}

// 延迟初始化
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 往哈希环上添加真实节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) //哈希函数需要传入字节类型的数据
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key //建立真实节点与虚拟节点的映射关系
		}
	}
	sort.Ints(m.keys)
}

// 根据key值获取虚拟节点并根据map映射为真实节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash //如果m.keys[i]的值比key对应的哈希值要大则证明m.keys[i]是其存放的节点
	})

	return m.hashMap[m.keys[idx%len(m.keys)]] //先获取虚拟节点然后再通过map映射至真实节点
}
