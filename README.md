# Worker Pool

A flexible Go worker pool model library using reflect.

Just toss all your jobs at it and it manages the concurrency using a pool of
 workers.

The only limitation is that the function that processes the job should only
 have one argument that receives the job parameter. If you need more
 parameters, consider a composite type.

## Look and Feel
```go
// the function that does the job
func Print(num int) {
	time.Sleep(time.Duration(6-num) * time.Second)
	fmt.Println(num)
}

// fn and concurrency
wp, err := New(Print, 10)
if err != nil {
	panic("can't create workerqueue " + err.Error())
}
for i := 0; i <= 5; i++ {
	// send a job
	wp.Work(i)
}
// close the job queue and start waiting all the job to be finished
wp.Wait()
```

For more examples, see `main_test.go`