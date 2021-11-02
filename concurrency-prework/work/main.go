package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	// s1 := Queue1{items: []uint64{1, 2, 3, 4, 5}}
	s2 := Queue2{items: []uint64{1, 2, 3, 4, 5}}
	s3 := Queue3{items: []uint64{1, 2, 3, 4, 5}}
	s4 := Queue4{items: []uint64{1, 2, 3, 4, 5}}
	var wg sync.WaitGroup
	n := 5
	num_q := 3
	wg.Add((5 * num_q) + num_q)
	// go func() {
	// 	defer wg.Done()
	// 	fmt.Println("Queue1")
	// 	printIDConcurrently(&s1, n, "Q1", &wg)
	// }()
	go func() {
		defer wg.Done()
		fmt.Println("Queue2")
		printIDConcurrently(&s2, n, "Q2", &wg)
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Queue3")
		printIDConcurrently(&s3, n, "Q3", &wg)
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Queue3")
		printIDConcurrently(&s4, n, "Q4", &wg)
	}()
	wg.Wait()

}

type idService interface {
	// Returns values in ascending order; it should be safe to call
	// getNext() concurrently without any additional synchronization.
	getNext() uint64
}

// Race conditions
type Queue1 struct {
	currentIdx int64    // index of next item to return
	items      []uint64 // items in ascending order
}

// Fix RC w/ Atomic
type Queue2 struct {
	currentIdx int64    // index of next item to return
	items      []uint64 // items in ascending order
}

// Fix RC w/ Mutex
type Queue3 struct {
	currentIdx int64    // index of next item to return
	items      []uint64 // items in ascending order
	mutex      sync.Mutex
}

// Fix RC w/ Channels
type Queue4 struct {
	currentIdx int64    // index of next item to return
	items      []uint64 // items in ascending order
}

func (q *Queue1) getNext() uint64 {
	item := q.items[q.currentIdx]
	q.currentIdx += 1
	return item
}

func (q *Queue2) getNext() uint64 {
	item := q.items[q.currentIdx]
	atomic.AddInt64(&q.currentIdx, 1)
	return item
}

func (q *Queue3) getNext() uint64 {
	q.mutex.Lock()
	item := q.items[q.currentIdx]
	q.currentIdx += 1
	q.mutex.Unlock()
	return item
}

type channelService struct {
	requests  chan struct{}
	responses chan uint64
	idx       uint64
}

func makeChannelService() *channelService {
	service := &channelService{
		requests:  make(chan struct{}),
		responses: make(chan uint64),
	}
	service.Start()
	return service
}

func (c *channelService) Start() {
	go func() {
		for range c.requests {
			c.responses <- c.idx
			c.idx += 1
		}
	}()
}

func (q *Queue4) getNext() uint64 {
	service := makeChannelService()
	service.requests <- struct{}{}
	idx := <-service.responses
	return q.items[idx]
}

func printIDConcurrently(s idService, n int, label string, wg *sync.WaitGroup) {
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			res := s.getNext()
			fmt.Printf("%s: %v\n", label, res)
		}()
	}
}
