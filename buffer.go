package main

// RingGrowing is a growing ring buffer.
// Not thread safe.
type RingGrowing struct {
	data     []interface{}
	n        int // Size of Data
	beg      int // First available element
	readable int // Number of data items available
}

// NewRingGrowing constructs a new RingGrowing instance with provided parameters.
func NewRingGrowing(initialSize int) *RingGrowing {
	return &RingGrowing{
		data: make([]interface{}, initialSize),
		n:    initialSize,
	}
}

// ReadOne reads (consumes) first item from the buffer if it is available, otherwise returns false.
func (r *RingGrowing) ReadOne() (data interface{}, ok bool) {
	if r.readable == 0 {
		return nil, false
	}
	r.readable--
	element := r.data[r.beg]
	r.data[r.beg] = nil // Remove reference to the object to help GC
	if r.beg == r.n-1 {
		// Was the last element
		r.beg = 0
	} else {
		r.beg++
	}
	return element, true
}

// WriteOne adds an item to the end of the buffer, growing it if it is full.
func (r *RingGrowing) WriteOne(data interface{}) {
	if r.readable == r.n {
		// Time to grow
		newN := r.n * 2
		newData := make([]interface{}, newN)
		to := r.beg + r.readable
		if to <= r.n {
			copy(newData, r.data[r.beg:to])
		} else {
			copied := copy(newData, r.data[r.beg:])
			copy(newData[copied:], r.data[:(to%r.n)])
		}
		r.beg = 0
		r.data = newData
		r.n = newN
	}
	r.data[(r.readable+r.beg)%r.n] = data
	r.readable++
}
