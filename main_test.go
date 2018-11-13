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
}

func TestNew(t *testing.T) {
	wp, err := New(PrintNum, 10)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i <= 5; i++ {
		wp.Work(i)
	}
	wp.Close()
	fmt.Println("waiting...")
	wp.Wait()

	wp1, err := New(PrintStruct, 3)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i <= 5; i++ {
		t := T{strconv.Itoa(i)}
		wp1.Work(t)
	}
	time.Sleep(time.Duration(10) * time.Second)

	wp2, err := New(func(a, b string) {}, 3)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i <= 5; i++ {
		t := strconv.Itoa(i)
		wp2.Work(t)
	}
	time.Sleep(time.Duration(10) * time.Second)
}
