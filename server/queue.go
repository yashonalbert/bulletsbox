package server

import (
	"container/heap"
	"time"
)

// Item struct
// score 优先级 0-2^32 默认1024
// delay 延迟 ready队列必须为0 delay队列为0时状态为buried
type Item struct {
	score uint32
	delay uint32
	body  []byte
	index int
}

// NewItem return a Item
func NewItem(score uint32, delay uint32, body []byte) *Item {
	i := new(Item)
	i.score = score
	i.delay = delay
	i.body = body
	return i
}

// Reserve struct
// lock false: Item未占用 | true: Item占用
type Reserve struct {
	item *Item
	lock bool
}

// NewReserve return a Reserve
func NewReserve(item *Item, lock bool) *Reserve {
	r := new(Reserve)
	r.item = item
	r.lock = lock
	return r
}

// Queue struct
type Queue struct {
	ready    ReadyHeap
	delay	 DelayHeap
	reserve  *Reserve
	list     map[string]bool
	n        int
}

// NewQueue return a Queue
func NewQueue() *Queue {
	q := new(Queue)
	q.ready = make(ReadyHeap, 0)
	heap.Init(&q.ready)
	q.delay = make(DelayHeap, 0)
	heap.Init(&q.delay)
	q.reserve = NewReserve(new(Item), false)
	q.list = make(map[string]bool)
	go q.watchReady()
	go q.watchDelay()
	return q
}

func (q *Queue) enQueue(item *Item) {
	if item.delay == 0 {
		heap.Push(&q.ready, item)
	} else {
		item.delay = item.delay + uint32(time.Now().Unix())
		heap.Push(&q.delay, item)
	}
}

func (q *Queue) deQueue() {
	if len(q.ready) > 0 {
		q.reserve = NewReserve(q.ready.Pop().(*Item), true)
	}
}

func (q *Queue) watchReady() {
	for {
		if !q.reserve.lock {
			q.deQueue()
		}
	}
}

func (q *Queue) watchDelay() {
	for {
		if len(q.delay) > 0 && uint32(time.Now().Unix()) == q.delay[0].delay {
			heap.Push(&q.ready, q.delay.Pop().(*Item))
		}
	}
}
