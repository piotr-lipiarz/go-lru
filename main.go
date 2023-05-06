package main

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

type Key comparable
type Value any

type LRUCache[K Key, T Value] struct {
	length   int
	capacity int

	// list
	head *Node[K, T]
	tail *Node[K, T]

	// map
	cache map[K]*Node[K, T]
}

type Node[K Key, T Value] struct {
	key   K
	value T
	prev  *Node[K, T]
	next  *Node[K, T]
}

func Constructor[K Key, T Value](capacity int) LRUCache[K, T] {
	head := &Node[K, T]{}
	tail := &Node[K, T]{}
	head.next = tail
	tail.prev = head
	return LRUCache[K, T]{
		cache:    make(map[K]*Node[K, T], capacity),
		length:   0,
		capacity: capacity,
		head:     head,
		tail:     tail,
	}
}

func (l *LRUCache[K, T]) moveToHead(node *Node[K, T]) {
	l.remove(node)
	l.addToFront(node)
}

func (l *LRUCache[K, T]) remove(node *Node[K, T]) {
	// remove node
	prev := node.prev
	next := node.next

	prev.next = next
	next.prev = prev

	node.next = nil
	node.prev = nil
}

func (l *LRUCache[K, T]) addToFront(node *Node[K, T]) {
	next := l.head.next
	l.head.next = node
	next.prev = node
	node.next = next
	node.prev = l.head
}

func (l *LRUCache[K, T]) removeLast() *Node[K, T] {
	x := l.tail.prev
	l.remove(x)
	return x
}

func (l *LRUCache[K, T]) Get(key K) (T, error) {
	node, found := l.cache[key]
	if !found {
		var zeroValue T
		return zeroValue, ErrNotFound
	}
	l.moveToHead(node)
	return node.value, nil
}

func (l *LRUCache[K, T]) Put(key K, value T) {
	node, found := l.cache[key]
	if found {
		// update
		node.value = value
		l.moveToHead(node)
		return
	}
	// not found
	// add to map & add to DLL
	n := &Node[K, T]{
		value: value,
		key:   key,
	}
	l.cache[key] = n
	l.addToFront(n)
	// evict DLL tail & remove from map
	if l.length < l.capacity {
		l.length++
	} else {
		last := l.removeLast()
		delete(l.cache, last.key)
	}
}
