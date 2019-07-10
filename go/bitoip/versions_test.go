package bitoip

import (
	"gotest.tools/assert"
	"testing"
)

func TestVersion_Bytes(t *testing.T) {
	v := Version{
		1,
		2,
		3,
		POC,
	}
	assert.Equal(t, v.String(), "1.2.3-pre-alpha")
	assert.DeepEqual(t, v.Bytes(), []byte{01, 02, 03, byte(POC)})

}
