package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidSignature(t *testing.T) {
	// echo -n "data" | openssl dgst -sha256 -hmac secret
	signature := "1b2c16b75bd2a870c114153ccda5bcfca63314bc722fa160d690de133ccbb9db"
	result, err := isValidSignature("secret", signature, []byte("data"))

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestInvalidSignature(t *testing.T) {
	signature := "invalid"
	result, err := isValidSignature("secret", signature, []byte("data"))

	assert.Nil(t, err)
	assert.False(t, result)
}
