package main

import (
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
)

// ------------------------------------
//   Worker pool
// ------------------------------------

type Result[T any] struct {
	Idx int64
	Val T
	Err error
}

type Task[T any] struct {
	Idx int64
	Fnc func() (T, error)
}

type WorkerPool[T any] struct {
	results    []*Result[T]
	isDone     bool
	taskCount  int64
	taskChan   chan *Task[T]
	resultChan chan *Result[T]
	waitTask   sync.WaitGroup
	waitResult sync.WaitGroup
}

func NewWorkerPool[T any](count int) *WorkerPool[T] {
	if count == 0 {
		count = runtime.NumCPU()
	}
	if count < 0 {
		return nil
	}

	pool := &WorkerPool[T]{
		taskChan:   make(chan *Task[T]),
		resultChan: make(chan *Result[T]),
		results:    make([]*Result[T], 0),
		waitTask:   sync.WaitGroup{},
	}

	for i := 0; i < count; i++ {
		pool.waitTask.Add(1)
		go pool.doTask(i)
	}

	pool.waitResult.Add(1)
	go pool.collectResult()

	return pool
}

func (pool *WorkerPool[T]) doTask(wn int) {
	for task := range pool.taskChan {
		value, err := task.Fnc()
		pool.resultChan <- &Result[T]{task.Idx, value, err}
	}
	pool.waitTask.Done()
}

func (pool *WorkerPool[T]) collectResult() {
	for res := range pool.resultChan {
		pool.results = append(pool.results, res)
	}
	sort.Slice(pool.results, func(i, j int) bool {
		return pool.results[i].Idx < pool.results[j].Idx
	})
	pool.isDone = true
	pool.waitResult.Done()
}

// Submit a new task for execution
func (pool *WorkerPool[T]) Submit(f func() (T, error)) {
	task := &Task[T]{
		Idx: atomic.AddInt64(&pool.taskCount, 1),
		Fnc: f,
	}
	pool.taskChan <- task
}

// Close the worker pool
// Tasks can no longer be submitted
// This function must be explicitly called before Wait
func (pool *WorkerPool[T]) Close() {
	close(pool.taskChan)
}

// Wait blocks until all tasks are finished and the results are collected
func (pool *WorkerPool[T]) Wait() []*Result[T] {
	pool.waitTask.Wait()
	close(pool.resultChan)
	pool.waitResult.Wait()
	return pool.results
}

// Results prevents accessing to the underlying result array while it still not completed
func (pool WorkerPool[T]) Results() []*Result[T] {
	if pool.isDone {
		return pool.results
	}
	return nil
}
