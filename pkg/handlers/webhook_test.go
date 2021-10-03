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

func TestRouteParse(t *testing.T) {
	envVars := []string{
		"RULE_TEST=event|address|jq rule",
	}
	routes := getRoutes(envVars)
	assert.Len(t, routes, 1)
	assert.Len(t, routes[0].events, 1)
	assert.Equal(t, "event", routes[0].events[0])
	assert.Equal(t, "address", routes[0].route)
	assert.Equal(t, "RULE_TEST", routes[0].name)
	assert.Equal(t, "jq rule", routes[0].query)
}

func TestRouteParseMultiple(t *testing.T) {
	envVars := []string{
		"RULE_TEST=event|address|jq rule",
		"RULE_FOO=event2|address2|jq rule2",
	}
	routes := getRoutes(envVars)
	assert.Len(t, routes, 2)
	assert.Equal(t, "event2", routes[1].events[0])
}

func TestRouteParseInvaid(t *testing.T) {
	envVars := []string{
		"RULE_TEST=event|address|jq rule",
		"RULE_FOO=event|address",
	}
	routes := getRoutes(envVars)
	assert.Len(t, routes, 1)
}

func TestRouteMultipleEvents(t *testing.T) {
	envVars := []string{
		"RULE_TEST=event1,event2|address|jq rule",
	}
	routes := getRoutes(envVars)
	assert.Len(t, routes, 1)
	assert.Len(t, routes[0].events, 2)
}
