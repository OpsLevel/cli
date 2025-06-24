package e2e

import (
	"strings"
	"testing"
)

func TestSystemHappyPath(t *testing.T) {
	tc := CLITest{
		Steps: []Step{
			Example{
				Cmd: "example system",
				Yaml: `
description: example_description
name: example_name
note: example_note
ownerId: Z2lkOi8vc2VydmljZS8xMjM0NTY3ODk
parent:
  alias: domain-alias
`,
			},
			Create{
				Cmd: "create system",
				Input: `
name: "Integration Test System"
description: "Created by integration test"
`,
			},
			Get{
				Cmd: "get system",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "Integration Test System") {
						u.Fatalf("get after create failed: %s", out)
					}
				},
			},
			List{
				Cmd: "list system",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "Integration Test System") {
						u.Fatalf("list missing system: %s", out)
					}
				},
			},
			Update{
				Cmd: "update system",
				Input: `
name: "Integration Test System Updated"
description: "Updated by integration test"
`,
			},
			Get{
				Cmd: "get system",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "Integration Test System Updated") || !strings.Contains(out, "Updated by integration test") {
						u.Fatalf("update1 failed\nout: %s", out)
					}
				},
			},
			// TODO: description cannot be unset yet
			//			Update{
			//				Cmd: "update system",
			//				Input: `
			//name: "Integration Test System Updated Again"
			//description: null
			//`,
			//			},
			//			Get{
			//				Cmd: "get system",
			//				Validate: func(u *Utility, out string) {
			//					if !strings.Contains(out, "Integration Test System Updated Again") || strings.Contains(out, "Updated by integration test") {
			//						u.Fatalf("update2 failed (description should be unset)\nout: %s", out)
			//					}
			//				},
			//			},
			Delete{
				Cmd: "delete system",
			},
			Missing{
				Cmd: "get system",
			},
		},
	}
	tc.Run(t)
}
