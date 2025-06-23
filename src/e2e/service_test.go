package e2e

import (
	"strings"
	"testing"
)

func TestServiceHappyPath(t *testing.T) {
	tc := CLITest{
		Steps: []Step{
			Example{
				Cmd: "example service",
				Yaml: `
name: example_name
product: example_product
description: example_description
language: example_language
framework: example_framework
tier: example_alias
lifecycle: example_alias
skipAliasesValidation: false
`,
			},
			Create{
				Cmd: "create service",
				Input: `
name: "TestService"
description: "Created by integration test"
`,
			},
			Get{
				Cmd: "get service",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "TestService") {
						u.Fatalf("get after create failed: %s", out)
					}
				},
			},
			List{
				Cmd: "list service",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "TestService") {
						u.Fatalf("list missing service: %s", out)
					}
				},
			},
			Update{
				Cmd: "update service",
				Input: `
name: "TestServiceUpdated"
description: "Updated by integration test"
`,
			},
			Get{
				Cmd: "get service",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "TestServiceUpdated") || !strings.Contains(out, "Updated by integration test") {
						u.Fatalf("update1 failed\nout: %s", out)
					}
				},
			},
			// TODO: description cannot be unset yet
			//			Update{
			//				Cmd: "update service",
			//				Input: `
			//name: "TestServiceUpdatedAgain"
			//description: null
			//`,
			//			},
			//			Get{
			//				Cmd: "get service",
			//				Validate: func(u *Utility, out string) {
			//					if !strings.Contains(out, "TestServiceUpdatedAgain") || strings.Contains(out, "Updated by integration test") {
			//						u.Fatalf("update2 failed (description should be unset)\nout: %s", out)
			//					}
			//				},
			//			},
			Delete{
				Cmd: "delete service",
			},
			Missing{
				Cmd: "get service",
			},
		},
	}
	tc.Run(t)
}
