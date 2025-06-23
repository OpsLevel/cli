package e2e

import (
	"strings"
	"testing"
)

func TestTeamHappyPath(t *testing.T) {
	tc := CLITest{
		Steps: []Step{
			Example{
				Cmd: "example team",
				Yaml: `
managerEmail: example_manager_email
name: example_name
parentTeam:
  alias: example_parent_team
responsibilities: example_responsibilities
`,
			},
			Create{
				Cmd: "create team",
				Input: `
name: "TestTeam"
responsibilities: "Created by integration test"
`,
			},
			Get{
				Cmd: "get team",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "TestTeam") {
						u.Fatalf("get after create failed: %s", out)
					}
				},
			},
			List{
				Cmd: "list team",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "TestTeam") {
						u.Fatalf("list missing team: %s", out)
					}
				},
			},
			Update{
				Cmd: "update team",
				Input: `
name: "TestTeam Updated"
responsibilities: "Updated by integration test"
`,
			},
			Get{
				Cmd: "get team",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "TestTeam Updated") || !strings.Contains(out, "Updated by integration test") {
						u.Fatalf("update1 failed\nout: %s", out)
					}
				},
			},
			// TODO: responsibilities cannot be unset yet
			//			Update{
			//				Cmd: "update team",
			//				Input: `
			//name: "Integration Test Team Updated Again"
			//responsibilities: null
			//`,
			//			},
			//			Get{
			//				Cmd: "get team",
			//				Validate: func(u *Utility, out string) {
			//					if !strings.Contains(out, "Integration Test Team Updated Again") || strings.Contains(out, "Updated by integration test") {
			//						u.Fatalf("update2 failed (responsibilities should be unset)\nout: %s", out)
			//					}
			//				},
			//			},
			Delete{
				Cmd: "delete team",
			},
			Missing{
				Cmd: "get team",
			},
		},
	}
	tc.Run(t)
}
