package singleflight

import "sync"

// call代表正在进行中或已经结束的请求，使用sync.WaitGroup避免重入发生缓存雪崩问题
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group 用于管理不同的请求call
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized
}

// 实际调用的方法，虽然可能会有N个请求，但倘若缓存中无数据的话只会访问数据库一次，其他的请求直接返回缓存值，不会访问数据库
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil { //懒加载，直到使用的时候才进行初始化
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
