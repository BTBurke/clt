# Snapshot

Snapshot is a simple library for snapshot testing in Go.  I find it useful for testing complex command-line applications when I'm interested in finding regressions to the UI.  When you run a test for the first time, the output is written to a snapshot file.  Subsequent tests ensure that test output matches your snapshot.

```go
package example

import (
	"testing"

	"github.com/BTBurke/snapshot"
)

func TestOutput(t *testing.T) {

	// This output matches the testoutput.snap in the __snapshots__ directory
	output := []byte("This is my UI output\nAnd it can be complex")
	snapshot.Assert(t, output)

	// If I assert something different, I'll get a test failure and a diff of what changed
	output = []byte("This is my UI output\nAnd it can be *very* complex")
	snapshot.Assert(t, output)

	// Output:
	// Snapshot test failed for: TestOutput.  Diff:
	//
	// --- Expected
	// +++ Received
	// @@ -1,2 +1,2 @@
	// This is my UI output
	// -And it can be complex
	// +And it can be *very* complex
}
```

## Updating snapshots

When you make UI changes and want to create new snapshots, set the environment variable `UPDATE_SNAPSHOTS` before running your tests.

```
UPDATE_SNAPSHOTS=true go test -v -cover
```

## Setting a custom snapshot directory and context

By default, snapshots are saved in the `__snapshots__` directory relative to the working directory where you run `go test`.  To change this behavior, you can create your own configuration prior to running your snapshot assertions.  The snapshot directory should be the full path.  You can also change the number of lines of context that are shown around your diffs.

```go
import (
	"testing"

	"github.com/BTBurke/snapshot"
)

func TestOutput(t *testing.T) {

    mycfg, _ := snapshot.New(SnapDirectory("/snaps"), ContextLines(0))

	// This output matches the testoutput.snap in the /snaps directory
	output := []byte("This is my UI output\nAnd it can be complex")
	mycfg.Assert(t, output)

	// If I assert something different, I'll get a test failure and a diff of what changed
	output = []byte("This is my UI output\nAnd it can be *very* complex")
	mycfg.Assert(t, output)

	// Output:
	// Snapshot test failed for: TestOutput.  Diff:
	//
	// --- Expected
	// +++ Received
	// @@ -1,2 +1,2 @@
	// -And it can be complex
	// +And it can be *very* complex
}
```