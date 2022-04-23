package common_test

import (
	"testing"

	common "github.com/opslevel/cli/common"

	"github.com/rocktavious/autopilot"
)

func TestMinInt(t *testing.T) {
	// Arrange
	// Act
	// Assert
	autopilot.Equals(t, 4, common.MinInt(4, 5, 6, 7))
	autopilot.Equals(t, 4, common.MinInt(7, 6, 5, 4))
	autopilot.Equals(t, 1, common.MinInt(10, 1, 9, 7))
}
