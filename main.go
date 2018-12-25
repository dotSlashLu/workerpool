package workerpool

import (
	"errors"
	"reflect"
	"sync"
)

type Pool struct {
	c      chan []interface{}
	wg     *sync.WaitGroup
	closed bool
}

// Create a new worker pool.
func New(fn interface{}, concurrency int) (wp *Pool, err error) {
	wp = new(Pool)
	wp.wg = new(sync.WaitGroup)

	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		err = errors.New("fn should be a function")
		return
	}
	wp.c = make(chan []interface{})

	// start the pool
	fnV := reflect.ValueOf(fn)
	wp.wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wp.wg.Done()
			for {
				m, ok := <-wp.c
				if !ok {
					return
				}
				args := make([]reflect.Value, len(m))
				for i, argv := range m {
					args[i] = reflect.ValueOf(argv)
				}
				fnV.Call(args)
			}
		}()
	}
	return
}

// Add a job to the job queue for processing.
func (wp *Pool) Work(args ...interface{}) {
	msg := make([]interface{}, len(args))
	for i, m := range args {
		msg[i] = m
	}
	wp.c <- args
}

// Close the job queue channel and wait for all the jobs to be done.
func (wp *Pool) Wait() {
	wp.Close()
	wp.wg.Wait()
}

// Close the job queue channel.
func (wp *Pool) Close() {
	if wp.closed {
		return
	}
	close(wp.c)
	wp.closed = true
}
