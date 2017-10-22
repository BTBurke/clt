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
