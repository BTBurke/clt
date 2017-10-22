package snapshot

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	tt := []struct {
		Name string
		Exp  string
		Recv string
		Diff string
	}{
		{Name: "no diff", Exp: "test\nstring", Recv: "test\nstring", Diff: ""},
		{Name: "diff", Exp: "test\nstring", Recv: "test\nstring2", Diff: "--- Expected\n+++ Received\n@@ -1,2 +1,2 @@\n test\n-string\n+string2\n"},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			diff, err := getDiff([]byte(tc.Exp), []byte(tc.Recv))
			assert.NoError(t, err)
			assert.Equal(t, tc.Diff, diff, "Expected diff:\n%s\nGot:\n%s\n", tc.Diff, diff)
		})
	}
}

func TestSnapFilename(t *testing.T) {
	tt := []struct {
		Name string
		In   string
		Out  string
	}{
		{Name: "no spaces", In: "TestRun", Out: "testrun.snap"},
		{Name: "with spaces", In: "Test Run", Out: "test-run.snap"},
		{Name: "i'm with stupid ✓", In: "i'm with stupid ✓", Out: "i-m-with-stupid-✓.snap"},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Out, getSnapFilename(tc.In))
		})
	}
}

func TestCreateSnaps(t *testing.T) {
	tt := []struct {
		Name string
		Out  string
	}{
		{Name: "single line output", Out: "this is a single line of output with UTF-8 ✓"},
		{Name: "multi line output", Out: "this is\nmulti line output with newlines\n\n"},
		{Name: "multi line with special chars", Out: "this is\nmulti\t{line}&\r\n"},
	}

	tmpdir, err := ioutil.TempDir("", "snap")
	if err != nil {
		t.Fatalf("Unexpected error creating temp dir: %s", err)
	}
	defer os.RemoveAll(tmpdir)

	c, _ := New(SnapDirectory(tmpdir))

	if err := os.Setenv("UPDATE_SNAPSHOTS", "true"); err != nil {
		t.Fatalf("Unexpected error setting environment: %s", err)
	}
	defer os.Unsetenv("UPDATE_SNAPSHOTS")

	// test creating snaps
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {

			c.Assert(t, []byte(tc.Out))

			recv, err := ioutil.ReadFile(path.Join(tmpdir, getSnapFilename(t.Name())))
			if err != nil {
				t.Fatalf("Expected snap file %s to be created but does not exist", getSnapFilename(t.Name()))
			}
			assert.True(t, bytes.Equal([]byte(tc.Out), recv))
		})
	}
}

func TestSnaps(t *testing.T) {

	tt := []struct {
		Name string
		Out  string
	}{
		{Name: "single line output", Out: "this is a single line of output with UTF-8 ✓"},
		{Name: "multi line output", Out: "this is\nmulti line output with newlines\n\n"},
		{Name: "multi line with special chars", Out: "this is\nmulti\t{line}&\r\n"},
	}

	tmpdir, err := ioutil.TempDir("", "snap")
	if err != nil {
		t.Fatalf("Unexpected error creating temp dir: %s", err)
	}
	defer os.RemoveAll(tmpdir)

	c, _ := New(SnapDirectory(tmpdir))

	// Check all created snaps and expect no errors, then update
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := os.Unsetenv("UPDATE_SNAPSHOTS"); err != nil {
				t.Fatalf("Unexpected error unsetting environment: %s", err)
			}
			if err := createSnapshot(t.Name(), []byte(tc.Out), c.Directory); err != nil {
				t.Fatalf("Unexpected error writing snap file: %s", err)
			}
			c.Assert(t, []byte(tc.Out))

			if err := os.Setenv("UPDATE_SNAPSHOTS", "true"); err != nil {
				t.Fatalf("Unexpected error setting environment: %s", err)
			}
			c.Assert(t, []byte(tc.Out+"\nupdated"))
			recv, err := ioutil.ReadFile(path.Join(c.Directory, getSnapFilename(t.Name())))
			if err != nil {
				t.Fatalf("Unexpected error reading snap file: %s", err)
			}
			assert.True(t, bytes.Equal([]byte(tc.Out+"\nupdated"), recv))
		})
	}

	os.Unsetenv("UPDATE_SNAPSHOTS")

}

func TestConfig(t *testing.T) {
	c, err := New()
	assert.NoError(t, err)
	wd, _ := os.Getwd()
	assert.Equal(t, c.Directory, path.Join(wd, "__snapshots__"))
	assert.Equal(t, c.Context, 10)

	c2, err := New(SnapDirectory("/test"), ContextLines(20))
	assert.NoError(t, err)
	assert.Equal(t, c2.Directory, "/test")
	assert.Equal(t, c2.Context, 20)
}
