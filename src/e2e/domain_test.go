package e2e

import (
	"strings"
	"testing"
)

func TestDomainHappyPath(t *testing.T) {
	createInput := `
name: "Integration Test Domain"
description: "Created by integration test"
`
	updateInput1 := `
name: "Integration Test Domain Updated"
description: "Updated by integration test"
`
	updateInput2 := `
name: "Integration Test Domain Updated Again"
description: null
`
	domainName := "Integration Test Domain"
	updatedDomainName := "Integration Test Domain Updated"
	updatedAgainDomainName := "Integration Test Domain Updated Again"

	test := CLITest{
		Create: Create("create domain -f -", createInput),
		Get:    Get("get domain"),
		Delete: Delete("delete domain"),
		Steps: []Step{
			func(u *Utility) {
				out, err := u.Run("list domain")
				if err != nil || !strings.Contains(out, domainName) {
					u.Fatalf("list failed: %v\nout: %s", err, out)
				}
			},
			func(u *Utility) {
				out, err := u.Run("update domain "+u.ID+" -f -", updateInput1)
				if err != nil {
					u.Fatalf("update1 failed: %v\nout: %s", err, out)
				}
				out, err = u.Run("get domain " + u.ID)
				if err != nil || !strings.Contains(out, updatedDomainName) || !strings.Contains(out, "Updated by integration test") {
					u.Fatalf("get after update1 failed: %v\nout: %s", err, out)
				}
			},
			func(u *Utility) {
				out, err := u.Run("update domain "+u.ID+" -f -", updateInput2)
				if err != nil {
					u.Fatalf("update2 (unset) failed: %v\nout: %s", err, out)
				}
				out, err = u.Run("get domain " + u.ID)
				if err != nil || !strings.Contains(out, updatedAgainDomainName) || strings.Contains(out, "Updated by integration test") {
					u.Fatalf("get after update2 failed (description should be unset): %v\nout: %s", err, out)
				}
			},
		},
	}
	test.Run(t)
}
