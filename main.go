package workerpool

import (
	"errors"
	"reflect"
	"sync"
)

type Pool struct {
	C  reflect.Value
	wg sync.WaitGroup
}

// Create a new worker pool
func New(fn interface{}, concurrency int) (wp *Pool, err error) {
	wp = new(Pool)

	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func || fnType.NumIn() != 1 {
		err = errors.New("fn should be a function with one parameter " +
			"that receives job arguments")
		return
	}
	// make a chan of the first arg's type
	msgT := fnType.In(0)
	chanT := reflect.ChanOf(reflect.BothDir, msgT)
	wp.C = reflect.MakeChan(chanT, concurrency)

	// start the pool
	fnV := reflect.ValueOf(fn)
	for i := 0; i < concurrency; i++ {
		go func() {
			wp.wg.Add(1)
			for {
				m, ok := wp.C.Recv()
				if !ok {
					wp.wg.Done()
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
	wp.C.Send(reflect.ValueOf(msg))
}

// Wait for all the jobs to be done
// Should be executed after Close()
func (wp *Pool) Wait() {
	wp.wg.Wait()
}

// Close the work queue
func (wp *Pool) Close() {
	wp.C.Close()
}
