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
	ID string // The resource ID created by Create
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

type Step func(*Utility)

type CLITest struct {
	Create Step
	Get    [2]Step
	Delete Step
	Steps  []Step
}

func (ct *CLITest) Run(t *testing.T) {
	util := &Utility{T: t}
	defer func() {
		if util.ID != "" {
			ct.Delete(util)
		}
	}()

	t.Run("Create", func(t *testing.T) {
		util.T = t
		ct.Create(util)
	})
	t.Run("Get", func(t *testing.T) {
		util.T = t
		ct.Get[0](util) // Should exist after create
	})
	for i, step := range ct.Steps {
		t.Run("Step "+strconv.Itoa(i), func(t *testing.T) {
			util.T = t
			step(util)
		})
	}
	t.Run("Delete", func(t *testing.T) {
		util.T = t
		ct.Delete(util)
		ct.Get[1](util) // Should not exist after delete
		util.ID = ""    // Mark as deleted so defer doesn't try again
	})
}

func Create(cmd string, input string) Step {
	return func(u *Utility) {
		out, err := u.Run(cmd, input)
		if err != nil {
			u.Fatalf("create failed: %v\nout: %s", err, out)
		}
		u.ID = strings.TrimSpace(out)
		if u.ID == "" {
			u.Fatalf("expected ID, got: %q", out)
		}
	}
}

func Get(cmd string) [2]Step {
	return [2]Step{func(u *Utility) {
		out, err := u.Run(cmd + " " + u.ID)
		if err != nil {
			u.Fatalf("get failed: %v\nout: %s", err, out)
		}
	}, func(u *Utility) {
		out, err := u.Run(cmd + " " + u.ID)
		lower := strings.ToLower(out)
		if err == nil || !(strings.Contains(lower, "not found") || strings.Contains(lower, "missing") || strings.Contains(lower, "does not exist on this account")) {
			u.Fatalf("expected get after delete to fail with not found, got: %v\nout: %s", err, out)
		}
	}}
}

func Delete(cmd string) Step {
	return func(u *Utility) {
		out, err := u.Run(cmd + " " + u.ID)
		if err != nil {
			u.Fatalf("delete failed: %v\nout: %s", err, out)
		}
	}
}
