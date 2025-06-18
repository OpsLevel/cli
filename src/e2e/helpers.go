package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type Utility struct {
	*testing.T
	ID string // The resource ID created by a step, if needed
}

type Step struct {
	Name     string
	Run      func(*Utility)
	Deferred bool // If true, run only at the end as cleanup
}

type TestCase struct {
	Steps []Step
}

// Run executes the CLI using 'go run main.go' from the ./src directory with the given arguments and optional stdin, returning combined output and error.
func (u *Utility) Run(args string, stdin ...string) (string, error) {
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

func (tc *TestCase) Run(t *testing.T) {
	util := &Utility{T: t}
	// Run non-deferred steps
	for _, step := range tc.Steps {
		if step.Deferred {
			continue
		}
		t.Run(step.Name, func(t *testing.T) {
			util.T = t
			step.Run(util)
		})
	}
	// Run deferred steps
	for _, step := range tc.Steps {
		if !step.Deferred {
			continue
		}
		t.Run(step.Name, func(t *testing.T) {
			util.T = t
			step.Run(util)
		})
	}
}

func Create(cmd string, input string) Step {
	return Step{
		Name: "Create",
		Run: func(u *Utility) {
			out, err := u.Run(cmd, input)
			if err != nil {
				u.Fatalf("create failed: %v\nout: %s", err, out)
			}
			u.ID = strings.TrimSpace(out)
			if u.ID == "" {
				u.Fatalf("expected ID, got: %q", out)
			}
		},
	}
}

func Delete(cmd string) Step {
	return Step{
		Name:     "Delete",
		Deferred: true,
		Run: func(u *Utility) {
			out, err := u.Run(cmd + " " + u.ID)
			if err != nil {
				u.Fatalf("delete failed: %v\nout: %s", err, out)
			}
		},
	}
}

// Get returns a Step that runs the get command and validates the output using the provided function.
func Get(cmd string, validate func(u *Utility, stdout string)) Step {
	return Step{
		Name: "Get",
		Run: func(u *Utility) {
			out, err := u.Run(cmd + " " + u.ID)
			if err != nil {
				u.Fatalf("get failed: %v\nout: %s", err, out)
			}
			validate(u, out)
		},
	}
}

// List returns a Step that runs the list command and validates the output using the provided function.
func List(cmd string, validate func(u *Utility, stdout string)) Step {
	return Step{
		Name: "List",
		Run: func(u *Utility) {
			out, err := u.Run(cmd)
			if err != nil {
				u.Fatalf("list failed: %v\nout: %s", err, out)
			}
			validate(u, out)
		},
	}
}

// Update returns a Step that runs the update command with input and validates the output using the provided function.
func Update(cmd string, input string, validate func(u *Utility, stdout string)) Step {
	return Step{
		Name: "Update",
		Run: func(u *Utility) {
			out, err := u.Run(cmd+" "+u.ID+" -f -", input)
			if err != nil {
				u.Fatalf("update failed: %v\nout: %s", err, out)
			}
			validate(u, out)
		},
	}
}
