package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var expect = `package example

const (
	UserId       = "Id"
	UserAge      = "Age"
	UserName     = "UserName"
	UserNickName = "NickName"
)

const (
	UserDetailId      = "Id"
	UserDetailUserId  = "UserId"
	UserDetailAddress = "Address"
	UserDetailPicture = "Picture"
)
`

func TestMain_Gen(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := gen(buffer, "example/user.go")
	require.NoError(t, err)

	assert.Equal(t, expect, buffer.String())
}
