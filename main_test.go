package workerpool

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func PrintNum(num int) {
	time.Sleep(time.Duration(6-num) * time.Second)
	fmt.Println("PrintNum", num)
}

type T struct {
	Str string
}

func PrintStruct(t T) {
	fmt.Println("PrintStruct", t.Str)
	time.Sleep(time.Duration(1) * time.Second)
}

func TestNew(t *testing.T) {
	fmt.Println("should print consecutively in an interval of 1s")
	wp, err := New(PrintNum, 10)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i <= 5; i++ {
		wp.Work(i)
	}
	fmt.Println("waiting...")
	wp.Wait()

	fmt.Println("should print 3 items each time in an interval of 1s")
	wp1, err := New(PrintStruct, 3)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i <= 10; i++ {
		t := T{strconv.Itoa(i)}
		wp1.Work(t)
	}
	time.Sleep(time.Duration(1) * time.Second)

	// func with multiple args
	wp2, err := New(func(a, b string) {
		fmt.Println("func with two args", a, b)
	}, 3)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i <= 3; i++ {
		t := strconv.Itoa(i)
		t1 := strconv.Itoa(i + 1)
		wp2.Work(t, t1)
	}
	wp2.Close()

	// func with no args
	fmt.Println("func with no args")
	wp3, err := New(func() { fmt.Println("no args") }, 3)
	for i := 0; i <= 3; i++ {
		wp3.Work()
	}

	// variadic func
	fmt.Println("variadic func")
	wp4, err := New(func(args ...int) { fmt.Println(args) }, 3)
	wp4.Work()
	wp4.Work(2)
	wp4.Work(2, 3, 4)
}
