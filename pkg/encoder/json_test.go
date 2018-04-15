package encoder

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonEncoder(t *testing.T) {
	var buffer bytes.Buffer

	var encoder Encoder
	encoder = JsonEncoder{}
	testData := map[string]string{"key": "value"}
	expectedData := []byte(`{"key":"value"}`)

	err := encoder.Encode(testData, &buffer)

	actualData := buffer.Bytes()

	assert.Nil(t, err, "Error in encoder.Encode")
	assert.Equal(t, expectedData, actualData)
}
