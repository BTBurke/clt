package main

import (
	"fmt"
	"github.com/BTBurke/clt"
	"time"
)

func main() {

	// An example of a progress spinner that succeeds.  Calling Start()
	// starts a new go routine to render the spinner and returns control
	// to the calling function.  You must then call Success() to terminate
	// the go routine and show the user the OK.
	fmt.Println("\nDoing something that succeeds after 3 seconds:")
	p := clt.NewProgressSpinner("Testing a successful result")
	p.Start()
	time.Sleep(time.Duration(3) * time.Second)
	p.Success()

	// An example of a progress spinner that fails.  Calling Fail() will
	// let the user know the action failed.
	fmt.Println("\nDoing something that fails after 3 seconds:")
	pF := clt.NewProgressSpinner("Testing a failed result")
	pF.Start()
	time.Sleep(time.Duration(3) * time.Second)
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

}
