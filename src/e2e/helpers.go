package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

type Utility struct {
	*testing.T
	ID string // The resource ID created by a step, if needed
}

type Step interface {
	Run(u *Utility)
	Name() string
	Deferred() bool
}

// CLITest uses []Step interface now
type CLITest struct {
	Steps []Step
}

func (tc *CLITest) Run(t *testing.T) {
	util := &Utility{T: t}
	// Run non-deferred steps
	for i, step := range tc.Steps {
		if step.Deferred() {
			continue
		}
		t.Run(step.Name()+"_"+strconv.Itoa(i), func(t *testing.T) {
			util.T = t
			step.Run(util)
		})
	}
	// Run deferred steps
	for i, step := range tc.Steps {
		if !step.Deferred() {
			continue
		}
		t.Run(step.Name()+"_"+strconv.Itoa(i), func(t *testing.T) {
			util.T = t
			step.Run(util)
		})
	}
}

// Run executes the CLI using 'go run main.go' from the ./src directory with the given arguments and optional stdin, returning combined output and error.
func (u *Utility) Run(args string, stdin ...string) (string, error) {
	// TODO: need to allow using the pre-built binary
	// cmd := exec.Command("opslevel", strings.Split(args, " ")...)

	cmd := exec.Command("go", append([]string{"run", "main.go"}, strings.Split(args, " ")...)...)
	cmd.Dir = ".."
	cmd.Env = os.Environ()
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if len(stdin) > 0 {
		cmd.Stdin = strings.NewReader(stdin[0])
	}
	err := cmd.Run()
	return out.String() + errBuf.String(), err
}

// Create step
type Create struct {
	Cmd   string
	Input string
}

func (s Create) Run(u *Utility) {
	out, err := u.Run(s.Cmd+" -f -", s.Input)
	if err != nil {
		panic("create failed: " + err.Error() + "\nout: " + out)
	}
	u.ID = strings.TrimRight(strings.TrimLeft(strings.TrimSpace(out), "\""), "\"")
	if u.ID == "" {
		panic("expected ID, got: " + out)
	}
}

func (s Create) Name() string   { return "Create" }
func (s Create) Deferred() bool { return false }

// Get step
type Get struct {
	Cmd      string
	Validate func(u *Utility, out string)
}

func (s Get) Run(u *Utility) {
	out, err := u.Run(s.Cmd + " " + u.ID)
	if err != nil {
		u.Fatalf("get failed: %v\nout: %s", err, out)
	}
	s.Validate(u, out)
}

func (s Get) Name() string   { return "Get" }
func (s Get) Deferred() bool { return false }

// List step
type List struct {
	Cmd      string
	Validate func(u *Utility, out string)
}

func (s List) Run(u *Utility) {
	out, err := u.Run(s.Cmd)
	if err != nil {
		u.Fatalf("list failed: %v\nout: %s", err, out)
	}
	if s.Validate != nil {
		s.Validate(u, out)
	}
}

func (s List) Name() string   { return "List" }
func (s List) Deferred() bool { return false }

// Update step
type Update struct {
	Cmd      string
	Input    string
	Validate func(u *Utility, out string)
}

func (s Update) Run(u *Utility) {
	out, err := u.Run(s.Cmd+" -f - "+u.ID, s.Input)
	if err != nil {
		u.Fatalf("update failed: %v\nout: %s", err, out)
	}
	if s.Validate != nil {
		s.Validate(u, out)
	}
}

func (s Update) Name() string   { return "Update" }
func (s Update) Deferred() bool { return false }

type Delete struct {
	Cmd string
}

func (s Delete) Run(u *Utility) {
	out, err := u.Run(s.Cmd + " " + u.ID)
	if err != nil {
		u.Fatalf("delete failed: %v\nout: %s", err, out)
	}
}

func (s Delete) Name() string   { return "Delete" }
func (s Delete) Deferred() bool { return true }

type Missing struct {
	Cmd string
}

func (s Missing) Run(u *Utility) {
	out, err := u.Run(s.Cmd + " " + u.ID)
	lower := strings.ToLower(out)
	if err == nil || !(strings.Contains(lower, "not found") || strings.Contains(lower, "missing") || strings.Contains(lower, "does not exist on this account")) {
		u.Fatalf("expected get after delete to fail with not found, got: %v\nout: %s", err, out)
	}
}

func (s Missing) Name() string   { return "Missing" }
func (s Missing) Deferred() bool { return true }

type Example struct {
	Cmd  string
	Yaml string
}

func (s Example) Run(u *Utility) {
	out, err := u.Run(s.Cmd + " --yaml")
	if err != nil {
		panic("example failed: " + err.Error() + "\nout: " + out)
	}
	wip := strings.TrimSpace(out)
	expected := strings.TrimSpace(s.Yaml)
	if wip != expected {
		u.Fatalf("example mismatch for '%s'\nExpected:\n%s\nWIP:\n%s", s.Cmd, expected, wip)
	}
}

func (s Example) Name() string   { return "Example" }
func (s Example) Deferred() bool { return false }
