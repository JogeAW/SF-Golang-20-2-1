package ringbuffer

import "sync"

type RingBuffer struct {
	array  []int
	cursor int
	len    int
	mu     sync.Mutex
}

func NewRingBuffer(len int) *RingBuffer {
	return &RingBuffer{
		array:  make([]int, len),
		cursor: -1,
		len:    len,
		mu:     sync.Mutex{},
	}
}

func (r *RingBuffer) Append(el int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cursor == r.len-1 {
		for i := 1; i <= r.len-1; i++ {
			r.array[i-1] = r.array[i]
		}
		r.array[r.cursor] = el
	} else {
		r.cursor++
		r.array[r.cursor] = el
	}
}

func (r *RingBuffer) Read() []int {
	if r.cursor < 0 {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	var data []int = r.array[:r.cursor+1]
	r.cursor = -1

	return data
}
