# Golang分布式缓存

## 算法流程
1. 获取key值 --> 检查key值对应的val是否被缓存 -是-> 返回缓存值
2. 缓存值不存在 --> 使用一致性哈希选择远程节点 --> 通过http协议与远程节点进行交互 -->返回缓存值
## LRU 算法简介
最近最少使用，相对于仅考虑时间因素的 FIFO 和仅考虑访问频率的 LFU，LRU 算法可以认为是相对平衡的一种淘汰算法。LRU 认为，如果数据最近被访问过，那么将来被访问的概率也会更高。LRU 算法的实现非常简单，维护一个队列，如果某条记录被访问了，则移动到队尾，那么队首则是最近最少访问的数据，淘汰该条记录即可。
### LRU 算法实现
* 数据结构:双向链表，map
* 实现思路:使用golang标准库所提供的"container/list"双向链表，来作为lru算法底层的数据结构，当请求一条缓存数据的时候，会通过map尝试获取数据，
倘若获取成功则证明该缓存已存在，此时只需将该数据移到队首即可，倘若数据不存在，则需要将一条新数据插入到链表中，并移到队首，需要注意的是为了减少内存的浪费，每隔一段时间就对链表进行扫描，移除链表尾部的数据

## 一致性哈希算法简介
一致性哈希算法主要是用于解决如何选择远程节点的问题，将key值hash成一个数字，这个数字将被映射到2^32空间上并通过取余的方式形成一个哈希环，每次需要访问远程节点的时候，计算key值的hash值，然后沿着hash环
寻找第一个节点