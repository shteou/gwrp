package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidSignature(t *testing.T) {
	// echo -n "data" | openssl dgst -sha256 -hmac secret
	signature := "sha256=1b2c16b75bd2a870c114153ccda5bcfca63314bc722fa160d690de133ccbb9db"
	result, err := isValidSignature("secret", signature, []byte("data"))

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestValidSignatureWithLeadingsha256Chars(t *testing.T) {
	// TrimLeft was previously used to strip sha256= from the Signature header
	// This caused a bug because all characters from sha256= were being stripped in the
	// prefix, i.e. sigs starting with a would have that stripped
	// echo -n "data1" | openssl dgst -sha256 -hmac secret
	signature := "sha256=adbf386f8df2776192df3c30026d3bd19f01bff25dd0bfa10852caea9bc4759c"
	result, err := isValidSignature("secret", signature, []byte("data1"))

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

func TestRouteParseInvalid(t *testing.T) {
	envVars := []string{
		"RULE_TEST=event|address|jq rule",
		"RULE_FOO=event|address",
	}
	routes := getRoutes(envVars)
	assert.Len(t, routes, 1)
}

func TestRouteParseMultipleEvents(t *testing.T) {
	envVars := []string{
		"RULE_TEST=event1,event2|address|jq rule",
	}
	routes := getRoutes(envVars)
	assert.Len(t, routes, 1)
	assert.Len(t, routes[0].events, 2)
}

func TestRouteMatchesEvent(t *testing.T) {
	route := route{events: []string{"ping"}, name: "foo", query: ".", route: "somewhere"}
	matches, err := routeMatches([]byte("{}"), "ping", route)

	assert.Nil(t, err)
	assert.True(t, matches)
}

func TestRouteDoesntMatchEvent(t *testing.T) {
	route := route{events: []string{"ping"}, name: "foo", query: ".", route: "somewhere"}
	matches, err := routeMatches([]byte("{}"), "pull_request", route)

	assert.Nil(t, err)
	assert.False(t, matches)
}
