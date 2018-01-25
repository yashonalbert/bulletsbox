package server

// ReadyHeap struct
type ReadyHeap []*Item

func (rh ReadyHeap) Len() int {
	return len(rh)
}

func (rh ReadyHeap) Less(i, j int) bool {
	return rh[i].score < rh[j].score
}

func (rh ReadyHeap) Swap(i, j int) {
	rh[i], rh[j] = rh[j], rh[i]
	rh[i].index = i
	rh[j].index = j
}

// Push func
func (rh *ReadyHeap) Push(x interface{}) {
	n := len(*rh)
	item := x.(*Item)
	item.index = n
	*rh = append(*rh, item)
}

// Pop func
func (rh *ReadyHeap) Pop() interface{} {
	old := *rh
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*rh = old[0 : n-1]
	return item
}

