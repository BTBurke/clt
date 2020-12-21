package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/BTBurke/clt"
)

func main() {

	// This is a basic loading indicator that disappears after loading is complete.
	// Unlike progress bars or spinners, there is no indication of success or failure.
	// This is useful when making short server calls.  The delay parameter
	// prevents flashing the loading symbol.  If your call completes within
	// this delay parameter, the loading status will never be shown.
	fmt.Println("\nShowing a loading symbol while we make a remote call:")
	pL := clt.NewLoadingMessage("Loading...", clt.Dots, 100*time.Millisecond)
	pL.Start()
	time.Sleep(4 * time.Second)
	pL.Success()

	// An example of a progress spinner that succeeds.  Calling Start()
	// starts a new go routine to render the spinner and returns control
	// to the calling function.  You must then call Success() to terminate
	// the go routine and show the user the OK.
	fmt.Println("\nDoing something that succeeds after 3 seconds:")
	p := clt.NewProgressSpinner("Testing a successful result")
	p.Start()
	time.Sleep(3 * time.Second)
	p.Success()

	// An example of a progress spinner that fails.  Calling Fail() will
	// let the user know the action failed.
	fmt.Println("\nDoing something that fails after 3 seconds:")
	pF := clt.NewProgressSpinner("Testing a failed result")
	pF.Start()
	time.Sleep(3 * time.Second)
	pF.Fail()

	// An example of a progress bar that succeeds.  You must call
	// Update(<pct>) with the completion percentage (float64 between
	// 0.0 and 1.0).  Finally, call Success() or Fail() to terminate
	// the go routine.
	fmt.Println("\nDoing something that eventually succeeds:")
	pB := clt.NewProgressBar("Implement progress bar")
	pB.Start()
	for i := 0; i < 50; i++ {
		pB.Update(float64(i) / 50.0)
		time.Sleep(time.Duration(50) * time.Millisecond)
	}
	pB.Success()

	// An example of a progress bar that fails.  You must call
	// Update(<pct>) with the completion percentage (float64 between
	// 0.0 and 1.0).
	fmt.Println("\nDoing something that eventually fails:")
	pB2 := clt.NewProgressBar("Implement progress bar")
	pB2.Start()
	for i := 0; i < 20; i++ {
		pB2.Update(float64(i) / 50.0)
		time.Sleep(time.Duration(50) * time.Millisecond)
	}
	pB2.Fail()

	// An example of an incremental progress bar.  See incremental_example.go.
	fmt.Println("\nDoing incremental progress from multiple go routimes:")
	incremental()

}

func incremental() {
	// This is an incremental progress example.  It starts a progress bar with 10 total steps then some go routines
	// to simulate doing work and then updates the progress as each finishes

	p := clt.NewIncrementalProgressBar(10, "Doing work")

	ch := make(chan int, 1)
	var wg sync.WaitGroup

	// start 3 go routines to do some work.  Pass them a copy of the progress bar so they can
	// call increment after each task is done
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(ch chan int, p *clt.Progress) {
			defer wg.Done()
			for range ch {
				time.Sleep(time.Duration(rand.Intn(1000) * 1000000))
				p.Increment()
			}
		}(ch, p)
	}

	// start the bar then pass it some work to do
	p.Start()
	for i := 0; i < 10; i++ {
		ch <- i
	}
	// wait until all the work is done
	wg.Done()

	// call Success to close the progress channel and update to 100%
	p.Success()

}
