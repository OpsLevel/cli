package e2e

import (
	"strings"
	"testing"
)

func TestDomainHappyPath(t *testing.T) {
	test := CLITest{
		Create: Create("create domain -f -", `
name: "Integration Test Domain"
description: "Created by integration test"
`),
		Get:    Get("get domain"),
		Delete: Delete("delete domain"),
		Steps: []Step{
			func(u *Utility) {
				out, err := u.Run("list domain")
				if err != nil || !strings.Contains(out, "Integration Test Domain") {
					u.Fatalf("list failed: %v\nout: %s", err, out)
				}
			},
			func(u *Utility) {
				updateInput1 := `
name: "Integration Test Domain Updated"
description: "Updated by integration test"
`
				out, err := u.Run("update domain "+u.ID+" -f -", updateInput1)
				if err != nil {
					u.Fatalf("update1 failed: %v\nout: %s", err, out)
				}
				out, err = u.Run("get domain " + u.ID)
				if err != nil || !strings.Contains(out, "Integration Test Domain Updated") || !strings.Contains(out, "Updated by integration test") {
					u.Fatalf("get after update1 failed: %v\nout: %s", err, out)
				}
			},
			func(u *Utility) {
				updateInput2 := `
name: "Integration Test Domain Updated Again"
description: null
`
				out, err := u.Run("update domain "+u.ID+" -f -", updateInput2)
				if err != nil {
					u.Fatalf("update2 (unset) failed: %v\nout: %s", err, out)
				}
				out, err = u.Run("get domain " + u.ID)
				if err != nil || !strings.Contains(out, "Integration Test Domain Updated Again") || strings.Contains(out, "Updated by integration test") {
					u.Fatalf("get after update2 failed (description should be unset): %v\nout: %s", err, out)
				}
			},
		},
	}
	test.Run(t)
}
