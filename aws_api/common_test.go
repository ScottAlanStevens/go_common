package aws_api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestRequestV1 struct {
	Description string `json:"description" validate:"required"`
}

func TestCommon(t *testing.T) {

	t.Run("UnmarshalRequest", func(t *testing.T) {
		// Arrange
		testReq := TestRequestV1{
			Description: "Hi there!",
		}
		b, _ := json.Marshal(testReq)

		// Act
		var result TestRequestV1
		ok, err := UnmarshalRequest(string(b), &result)

		// Assert
		assert.Nil(t, err)
		assert.True(t, ok)
		assert.Equal(t, testReq, result)
	})

	t.Run("ValidateRequest", func(t *testing.T) {
		// Arrange
		testReq := TestRequestV1{
			Description: "Hi there!",
		}
		b, _ := json.Marshal(testReq)

		// Act
		result, errres := ValidateRequest[TestRequestV1](string(b))

		// Assert
		assert.Nil(t, errres)
		assert.Equal(t, testReq, *result)
	})

	t.Run("ValidateRequest Error", func(t *testing.T) {
		// Arrange
		testReq := TestRequestV1{}
		b, _ := json.Marshal(testReq)

		// Act
		result, errres := ValidateRequest[TestRequestV1](string(b))

		// Assert
		assert.NotNil(t, errres)
		assert.Nil(t, result)
	})
}
