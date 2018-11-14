package workerpool

import (
	"errors"
	"reflect"
	"sync"
)

type Pool struct {
	c  reflect.Value
	wg *sync.WaitGroup
}

// Create a new worker pool
func New(fn interface{}, concurrency int) (wp *Pool, err error) {
	wp = new(Pool)
	wp.wg = new(sync.WaitGroup)

	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func || fnType.NumIn() != 1 {
		err = errors.New("fn should be a function with one parameter " +
			"that receives job arguments")
		return
	}
	// make a chan of the first arg's type
	msgT := fnType.In(0)
	chanT := reflect.ChanOf(reflect.BothDir, msgT)
	wp.c = reflect.MakeChan(chanT, concurrency)

	// start the pool
	fnV := reflect.ValueOf(fn)
	wp.wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wp.wg.Done()
			for {
				m, ok := wp.c.Recv()
				if !ok {
					return
				}
				args := [1]reflect.Value{m}
				fnV.Call(args[:])
			}
		}()
	}
	return
}

// Add job to the queue for processing
func (wp *Pool) Work(msg interface{}) {
	wp.c.Send(reflect.ValueOf(msg))
}

// Wait for all the jobs to be done
// Should be executed after Close()
func (wp *Pool) Wait() {
	wp.wg.Wait()
}

// Close the work queue
func (wp *Pool) Close() {
	wp.c.Close()
}
