package server

// DelayHeap struct
type DelayHeap []*Item

func (rh DelayHeap) Len() int {
	return len(rh)
}

func (rh DelayHeap) Less(i, j int) bool {
	if rh[i].delay < rh[j].delay {
		return true
	}
	if rh[i].delay > rh[j].delay {
		return false
	}
	return rh[i].score < rh[j].score
}

func (rh DelayHeap) Swap(i, j int) {
	rh[i], rh[j] = rh[j], rh[i]
	rh[i].index = i
	rh[j].index = j
}

// Push func
func (rh *DelayHeap) Push(x interface{}) {
	n := len(*rh)
	item := x.(*Item)
	item.index = n
	*rh = append(*rh, item)
}

// Pop func
func (rh *DelayHeap) Pop() interface{} {
	old := *rh
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*rh = old[0 : n-1]
	return item
}


