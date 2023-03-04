package main

type LinkedList[T any] interface {
	PushBack(node T) *Node[T]
	Remove(node T) *Node[T]
	Next() *Node[T]
	Previous() *Node[T]
	First() *Node[T]
	Last() *Node[T]
	Len() int
}

type List[T any] struct {
	head *Node[T]
}

func (list *List[T]) Len() int {
	if list.head != nil {
		var l = 0
		for e := list.First(); e != nil; e = e.Next() {
			l += 1
		}
		return l
	} else {
		return 0
	}
}

func (list *List[T]) First() *Node[T] {
	if list.head == nil {
		return nil
	}
	return list.head
}

func (list *List[T]) Last() *Node[T] {
	if list.head == nil {
		return nil
	}
	for e := list.First(); e != nil; e = e.Next() {
		if e.Next() == nil {
			return e
		}
	}
	return nil
}

func (list *List[T]) PushBack(val T) *Node[T] {
	if list.head == nil {
		list.head = NewNode(nil, val, nil)
		return list.head
	}
	for e := list.First(); e != nil; e = e.Next() {
		if e.Next() == nil {
			e.next = NewNode(e, val, nil)
			return e.next
		}
	}
	return nil
}

func (list *List[T]) Remove(node *Node[T]) {
	for e := list.First(); e != nil; e = e.Next() {
		if e == node {
			if list.Len() > 1 {
				e.Remove()
			} else {
				list.head = nil
			}
		}
	}
}

type Node[T any] struct {
	next     *Node[T]
	val      T
	previous *Node[T]
}

func NewNode[T any](previous *Node[T], val T, next *Node[T]) *Node[T] {
	return &Node[T]{
		next:     next,
		val:      val,
		previous: previous,
	}
}

func (node *Node[T]) Next() *Node[T] {
	return node.next
}

func (node *Node[T]) Previous() *Node[T] {
	return node.previous
}

func (node *Node[T]) First() *Node[T] {
	if node.previous == nil {
		return node
	}
	return node.previous
}

func (node *Node[T]) Last() *Node[T] {
	if node.next == nil {
		return node
	}
	return node.next
}

func (node *Node[T]) PushBack(val T) *Node[T] {
	// get last element
	var last = node.Last()
	var newLast = NewNode(last, val, nil)
	newLast.previous = last
	last.next = newLast
	return last.next
}

func (node *Node[T]) Remove() {
	var oldPrevious = node.previous
	node.previous.next = node.next
	node.next.previous = oldPrevious
}

func (node *Node[T]) Len(l int) int {
	return node.Len(l + 1)
}
