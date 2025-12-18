package lru

import (
	"sync"
)

type idNode struct {
	id   uint32
	prev *idNode
	next *idNode
}

type LRUIds struct {
	capacity int
	size     int
	items    map[uint32]*idNode
	head     *idNode
	tail     *idNode
	mu       sync.Mutex
}

func NewLRUIds(capacity int) *LRUIds {
	if capacity <= 0 {
		return nil
	}
	return &LRUIds{
		capacity: capacity,
		items:    make(map[uint32]*idNode),
	}
}

func (l *LRUIds) Put(id uint32) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if nd, ok := l.items[id]; ok {
		l.moveToHeadLocked(nd)
		return
	}
	nd := &idNode{id: id}
	l.items[id] = nd
	l.addToHeadLocked(nd)
	l.size++
	if l.size > l.capacity {
		l.evictTailLocked()
	}
}

func (l *LRUIds) PushBack(id uint32) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.items[id]; ok {
		l.moveToHeadLocked(l.items[id])
		return
	}

	nd := &idNode{id: id}
	l.items[id] = nd

	nd.prev = l.tail
	nd.next = nil
	if l.tail != nil {
		l.tail.next = nd
	}
	l.tail = nd
	if l.head == nil {
		l.head = nd
	}
	l.size++
	for l.size > l.capacity {
		l.evictHeadLocked()
	}
}

func (l *LRUIds) LeastRecent() (uint32, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.tail == nil {
		return 0, false
	}
	return l.tail.id, true
}

func (l *LRUIds) MoveTailToHead() (uint32, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.tail == nil {
		return 0, false
	}
	nd := l.tail
	l.moveToHeadLocked(nd)
	return nd.id, true
}

func (l *LRUIds) addToHeadLocked(nd *idNode) {
	nd.prev = nil
	nd.next = l.head
	if l.head != nil {
		l.head.prev = nd
	}
	l.head = nd
	if l.tail == nil {
		l.tail = nd
	}
}

func (l *LRUIds) removeLocked(nd *idNode) {
	if nd.prev != nil {
		nd.prev.next = nd.next
	}
	if nd.next != nil {
		nd.next.prev = nd.prev
	}
	if nd == l.head {
		l.head = nd.next
	}
	if nd == l.tail {
		l.tail = nd.prev
	}
	nd.prev = nil
	nd.next = nil
}

func (l *LRUIds) moveToHeadLocked(nd *idNode) {
	if nd == l.head {
		return
	}
	l.removeLocked(nd)
	l.addToHeadLocked(nd)
}

func (l *LRUIds) evictTailLocked() {
	if l.tail == nil {
		return
	}
	ev := l.tail
	l.removeLocked(ev)
	delete(l.items, ev.id)
	l.size--
}

func (l *LRUIds) evictHeadLocked() {
	if l.head == nil {
		return
	}
	ev := l.head
	l.removeLocked(ev)
	delete(l.items, ev.id)
	l.size--
}
