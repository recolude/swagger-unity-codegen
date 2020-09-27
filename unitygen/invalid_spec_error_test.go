package unitygen_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/stretchr/testify/assert"
)

func Test_InvalidSpecError_GeneratesAppropriateErrorMessage(t *testing.T) {
	// ******************************** ARRANGE *******************************
	specErr := unitygen.InvalidSpecError{Path: []string{"some", "path"}, Reason: "Cause I said so"}

	// ********************************** ACT *********************************
	message := specErr.Error()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "Invalid spec at some.path: Cause I said so", message)
}

func Test_InvalidSpecError_GeneratesAppropriateErrorMessageWithNoPath(t *testing.T) {
	// ******************************** ARRANGE *******************************
	specErr := unitygen.InvalidSpecError{Path: nil, Reason: "Cause I said so"}

	// ********************************** ACT *********************************
	message := specErr.Error()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "Invalid spec: Cause I said so", message)
}
