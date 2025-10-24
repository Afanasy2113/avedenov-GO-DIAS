package hw04lrucache

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
	front *ListItem
	back  *ListItem
	len   int
}

// Back implements List.
func (l *list) Back() *ListItem {
	return l.back
}

// Front implements List.
func (l *list) Front() *ListItem {
	return l.front
}

// Len implements List.
func (l *list) Len() int {
	return l.len
}

// MoveToFront implements List.
func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}

	// Вставляем в начало
	i.Next = l.front
	i.Prev = nil
	if l.front != nil {
		l.front.Prev = i
	} else {
		l.back = i // На случай, если список был пуст (маловероятно)
	}
	l.front = i
	// len не меняется при перемещении
}

// PushBack implements List.
func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v, Prev: l.back}
	if l.back != nil {
		l.back.Next = newItem
	} else {
		l.front = newItem // Список был пуст
	}
	l.back = newItem
	l.len++
	return newItem
}

// PushFront implements List.
func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v, Next: l.front}
	if l.front != nil {
		l.front.Prev = newItem
	} else {
		l.back = newItem // Список был пуст
	}
	l.front = newItem
	l.len++
	return newItem
}

// Remove implements List.
func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next // i был первым
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev // i был последним
	}

	// Очищаем указатели у удалённого элемента
	i.Next = nil
	i.Prev = nil
	l.len--
}

func NewList() List {
	return new(list)
}
