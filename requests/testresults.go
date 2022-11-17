package requests

import (
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"github.com/fatih/color"
	"github.com/juju/ansiterm/tabwriter"
)

var (
	txtBold   = color.New(color.Bold)
	txtSucess = color.New(color.FgGreen)
	txtError  = color.New(color.FgRed)
)

type TestResult struct {
	Tests []*TestCase
}

func (t *TestResult) Run(path cue.Path, f func(t *TestCase)) {
	tc := TestCase{
		Path: path,
		Pass: true,
	}
	t.Tests = append(t.Tests, &tc)

	f(&tc)
}

func (t *TestResult) Print(verbose bool) {
	if len(t.Tests) == 0 {
		fmt.Println("No tests run.")
		return
	}
	if verbose {
		for _, test := range t.Tests {
			test.print()
		}
		fmt.Printf("\n\n")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\n", txtBold.Sprint("Test"), txtBold.Sprint("Result"), txtBold.Sprint("Msg"))
	for _, test := range t.Tests {
		fmt.Fprintf(w, "%s\t%s\t%s\n", test.Path, test.status(), test.Msg)
	}
	w.Flush()
}

type TestCase struct {
	Path   cue.Path
	Msg    string
	Pass   bool
	Output cue.Value
}

func (tc *TestCase) Fail(msg string) {
	tc.Msg = msg
	tc.Pass = false
}

func (tc *TestCase) Failf(format string, args ...any) {
	tc.Msg = fmt.Sprintf(format, args...)
	tc.Pass = false
}

func (tc *TestCase) status() string {
	if tc.Pass {
		return txtSucess.Sprint("PASS")
	}
	return txtError.Sprint("FAIL")
}

func (tc *TestCase) print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "===\tRUN:\t", tc.Path)
	fmt.Fprintln(w, "\tPath:\t", tc.Path)
	fmt.Fprintln(w, "\tStatus:\t", tc.status())
	fmt.Fprintln(w, "\tMsg:\t", tc.Msg)
	fmt.Fprintln(w, "\tOutput:\t")
	fmt.Fprintln(w, tc.Output)

	w.Flush()
}
