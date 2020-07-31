package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestGetGroupsFromAuthenticatedRequestNoGroups(t *testing.T) {
	r := &http.Request{}
	ret := GetGroupsFromAuthenticatedRequest(r)
	assert.Equal(t, []string{}, ret)
}

func TestGetGroupsFromAuthenticatedRequestWrongDataType(t *testing.T) {
	r := &http.Request{}

	claims := jwt.MapClaims{
		"groups": "garbage"}

	ctx := context.WithValue(r.Context(), interface{}("claims"), claims)

	ret := GetGroupsFromAuthenticatedRequest(r.WithContext(ctx))
	assert.Equal(t, []string{}, ret)
}

func TestGetGroupsFromAuthenticatedRequest(t *testing.T) {
	r := &http.Request{}

	claims := jwt.MapClaims{
		"groups": []interface{}{"groupA", "groupB"}}

	ctx := context.WithValue(r.Context(), interface{}("claims"), claims)

	ret := GetGroupsFromAuthenticatedRequest(r.WithContext(ctx))
	assert.Equal(t, []string{"groupA", "groupB"}, ret)
}
