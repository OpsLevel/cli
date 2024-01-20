package common_test

import (
	"fmt"
	"github.com/opslevel/cli/common"
	"reflect"
	"slices"
	"testing"
)

func TestGetAliases(t *testing.T) {
	type TestCase struct {
		s   string
		exp []string
	}
	tests := []TestCase{
		{
			s: "PropertyDefinition",
			exp: []string{
				"propertydefinition", "propertydefinitions",
				"property_definition", "property_definitions",
				"property-definition", "property-definitions",
				"pd", "pds",
			},
		},
		{
			s: "FooFooBar",
			exp: []string{
				"foofoobar", "foofoobars",
				"foo-foo-bar", "foo-foo-bars",
				"foo_foo_bar", "foo_foo_bars",
			},
		},
		{
			s: "IcedTea",
			exp: []string{
				"icedtea", "icedteas",
				"iced-tea", "iced-teas",
				"iced_tea", "iced_teas",
			},
		},
		{
			s:   "Soda",
			exp: []string{"soda", "sodas"},
		},
		{
			s:   "A",
			exp: []string{"a", "as"},
		},
		{
			s:   "AB",
			exp: []string{"ab", "abs", "a-b", "a-bs", "a_b", "a_bs"},
		},
	}
	for _, tt := range tests {
		slices.Sort(tt.exp)
		testName := fmt.Sprintf("GetAliases_%s", tt.s)
		t.Run(testName, func(t *testing.T) {
			got := common.GetAliases(tt.s)
			if len(tt.exp) != len(got) {
				t.Errorf("expected len of %d got %d", len(tt.exp), len(got))
			}
			if !reflect.DeepEqual(tt.exp, got) {
				t.Errorf("expected and got are not equal:\n%s\n%s\n", tt.exp, got)
			}
		})
	}
}
