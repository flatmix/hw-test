package hw04lrucache

func insertAfter(list *list, item *ListItem, newItem *ListItem) {
	newItem.Prev = item
	if item.Next == nil {
		newItem.Next = nil
		list.lastItem = newItem
	} else {
		newItem.Next = item.Next
		item.Next.Prev = newItem
	}
	item.Next = newItem
}

func insertBefore(list *list, item *ListItem, newItem *ListItem) {
	newItem.Next = item
	if item.Prev == nil {
		newItem.Prev = nil
		list.firstItem = newItem
	} else {
		newItem.Prev = item.Prev
		item.Prev.Next = newItem
	}
	item.Prev = newItem
}

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	firstItem *ListItem
	lastItem *ListItem
	len      int
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	insertBefore(l, l.firstItem, i)
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.firstItem = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		l.lastItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := ListItem{
		Value: v,
	}
	if l.firstItem == nil {
		l.firstItem = &newItem
		l.lastItem = &newItem
		newItem.Prev = nil
		newItem.Next = nil
	} else {
		insertBefore(l, l.firstItem, &newItem)
	}
	l.len++
	return &newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.len++
	newItem := ListItem{
		Value: v,
	}
	if l.lastItem == nil {
		l.PushFront(v)
	} else {
		insertAfter(l, l.lastItem, &newItem)
	}
	return &newItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Len() int {
	return l.len
}

func NewList() List {
	return new(list)
}
