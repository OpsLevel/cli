package e2e

import (
	"strings"
	"testing"
)

func TestInfrastructureHappyPath(t *testing.T) {
	tc := CLITest{
		Steps: []Step{
			Example{
				Cmd: "example infra",
				Yaml: `
schema: example_schema
owner: Z2lkOi8vc2VydmljZS8xMjM0NTY3ODk
provider:
  account: example_account
  name: example_provider_name
  type: example_provider_type
  url: example_external_url
data:
  endpoint: https://google.com
  engine: BigQuery
  name: my-big-query
  replica: false
`,
			},
			Create{
				Cmd: "create infra",
				Input: `
schema: "Database"
provider:
  account: "Dev - 123456789"
  name: "GCP"
  type: "BigQuery"
  url: "https://google.com"
data:
  name: "my-big-query"
  endpoint: "https://google.com"
  engine: "BigQuery"
  replica: false
`,
			},
			Get{
				Cmd: "get infra",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "my-big-query") {
						u.Fatalf("get after create failed: %s", out)
					}
				},
			},
			List{
				Cmd: "list infra",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "my-big-query") {
						u.Fatalf("list missing infra: %s", out)
					}
				},
			},
			Update{
				Cmd: "update infra",
				Input: `
schema: "Database"
data:
  name: "my-big-query-updated"
`,
			},
			Get{
				Cmd: "get infra",
				Validate: func(u *Utility, out string) {
					if !strings.Contains(out, "my-big-query-updated") {
						u.Fatalf("update1 failed\nout: %s", out)
					}
				},
			},
			Delete{
				Cmd: "delete infra",
			},
			Missing{
				Cmd: "get infra",
			},
		},
	}
	tc.Run(t)
}
