package main

/*
 */
import "C"

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	queueSize     = 10
	nMiners       = 100
	sliceSize     = uint64(1_000_000)
	lowerBitsMask = uint64(0x0fffffff)
)

type queue struct {
	items    [queueSize]uint64
	head     int
	tail     int
	length   int
	closed   bool
	mu       sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
}

func newQueue() *queue {
	q := &queue{}
	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)
	return q
}

func (q *queue) close() {
	q.mu.Lock()
	if !q.closed {
		q.closed = true
		q.notEmpty.Broadcast()
		q.notFull.Broadcast()
	}
	q.mu.Unlock()
}

func (q *queue) add(item uint64) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.length == queueSize && !q.closed {
		q.notFull.Wait()
	}
	if q.closed {
		return false
	}

	q.items[q.tail] = item
	q.tail = (q.tail + 1) % queueSize
	q.length++
	q.notEmpty.Signal()
	return true
}

func (q *queue) pop() (uint64, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.length == 0 && !q.closed {
		q.notEmpty.Wait()
	}
	if q.length == 0 && q.closed {
		return 0, false
	}

	item := q.items[q.head]
	q.head = (q.head + 1) % queueSize
	q.length--
	q.notFull.Signal()
	return item, true
}

func hash(x uint64) uint64 {
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	x = x ^ (x >> 31)
	return x
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := rng.Uint64()

	q := newQueue()
	var solutionMu sync.Mutex
	var solution uint64
	var wg sync.WaitGroup

	setSolution := func(value uint64) bool {
		solutionMu.Lock()
		defer solutionMu.Unlock()
		if solution != 0 {
			return false
		}
		solution = value
		return true
	}

	getSolution := func() uint64 {
		solutionMu.Lock()
		defer solutionMu.Unlock()
		return solution
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		sliceBase := sliceSize
		for {
			if getSolution() != 0 {
				q.close()
				return
			}
			if !q.add(sliceBase) {
				return
			}
			fmt.Printf("sent %d\n", sliceBase)
			sliceBase += sliceSize
		}
	}()

	for i := 0; i < nMiners; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				if getSolution() != 0 {
					return
				}

				sliceBase, ok := q.pop()
				if !ok {
					return
				}

				for candidate := sliceBase; candidate < sliceBase+sliceSize; candidate++ {
					hashed := candidate ^ seed
					for j := 0; j < 10; j++ {
						hashed = hash(hashed)
					}

					if hashed&lowerBitsMask == 0 {
						if setSolution(candidate) {
							fmt.Printf("miner found solution %d\n", candidate)
							q.close()
						}
						return
					}
				}
			}
		}()
	}

	wg.Wait()
	fmt.Println(getSolution())
}
