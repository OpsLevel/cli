package e2e

import (
	"strings"
	"testing"
)

func TestDomainHappyPath(t *testing.T) {
	tc := CLITest{
		Steps: []Step{
			Create{
				Cmd: "create domain",
				Input: `
name: "Integration Test Domain"
description: "Created by integration test"
`,
			},
			Get{
				Cmd: "get domain",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "Integration Test Domain") {
						u.Fatalf("get after create failed: %s", out)
					}
				},
			},
			List{
				Cmd: "list domain",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "Integration Test Domain") {
						u.Fatalf("list missing domain: %s", out)
					}
				},
			},
			Update{
				Cmd: "update domain",
				Input: `
name: "Integration Test Domain Updated"
description: "Updated by integration test"
`,
			},
			Get{
				Cmd: "get domain",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "Integration Test Domain Updated") || !strings.Contains(out, "Updated by integration test") {
						u.Fatalf("update1 failed\nout: %s", out)
					}
				},
			},
			// TODO: description cannot be unset yest
			//			Update{
			//				Cmd: "update domain",
			//				Input: `
			//name: "Integration Test Domain Updated Again"
			//description: null
			//`,
			//			},
			//			Get{
			//				Cmd: "get domain",
			//				Validate: func(u *Utility, out string) {
			//					if !strings.Contains(out, "Integration Test Domain Updated Again") || strings.Contains(out, "Updated by integration test") {
			//						u.Fatalf("update2 failed (description should be unset)\nout: %s", out)
			//					}
			//				},
			//			},
			Delete{
				Cmd: "delete domain",
			},
			Missing{
				Cmd: "get domain",
			},
		},
	}
	tc.Run(t)
}
