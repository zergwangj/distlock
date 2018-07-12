package distlock

import (
	"testing"
	"sync"
	"time"
	"fmt"
)

func testLock(p *sync.WaitGroup, i int) {
	dl, err := NewDistLock(":8500", "test")
	if err != nil {
		fmt.Println("Create distributed lock failed")
		p.Done()
		return
	}
	fmt.Println("Create distributed lock success")
	defer dl.Destroy()
	stopCh := make(chan struct{})
	lockCh, err := dl.AquireLock(stopCh)
	if err != nil {
		fmt.Println("Aquire distributed lock failed")
		p.Done()
		return
	}
	fmt.Println("Aquire distributed lock success")
	time.Sleep(time.Second * 1)
	dl.ReleaseLock()

	fmt.Println("Distributed unlock success")

	<-lockCh
	p.Done()
}

func TestDistLock(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go testLock(&wg, i)

	}
	wg.Wait()
}

