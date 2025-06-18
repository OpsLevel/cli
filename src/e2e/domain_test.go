package e2e

import (
	"strings"
	"testing"
)

func TestDomainHappyPath(t *testing.T) {
	tc := TestCase{
		Steps: []Step{
			Create("create domain -f -", `
name: "Integration Test Domain"
description: "Created by integration test"
`),
			List("list domain", func(u *Utility, out string) {
				if !strings.Contains(out, "Integration Test Domain") {
					u.Fatalf("list missing domain: %s", out)
				}
			}),
			Update("update domain", `
name: "Integration Test Domain Updated"
description: "Updated by integration test"
`, func(u *Utility, out string) {
				if !strings.Contains(out, "Integration Test Domain Updated") || !strings.Contains(out, "Updated by integration test") {
					u.Fatalf("get after update1 failed\nout: %s", out)
				}
			}),
			Get("get domain", func(u *Utility, out string) {
				if !strings.Contains(out, "Integration Test Domain Updated") || !strings.Contains(out, "Updated by integration test") {
					u.Fatalf("get after update1 failed: %s", out)
				}
			}),
			Update("update domain", `
name: "Integration Test Domain Updated Again"
description: null
`, func(u *Utility, out string) {
				if !strings.Contains(out, "Integration Test Domain Updated Again") || strings.Contains(out, "Updated by integration test") {
					u.Fatalf("get after update2 failed (description should be unset)\nout: %s", out)
				}
			}),
			Get("get domain", func(u *Utility, out string) {
				if !strings.Contains(out, "Integration Test Domain Updated Again") || strings.Contains(out, "Updated by integration test") {
					u.Fatalf("get after update2 failed (description should be unset): %s", out)
				}
			}),
			Delete("delete domain"),
		},
	}
	tc.Run(t)
}
