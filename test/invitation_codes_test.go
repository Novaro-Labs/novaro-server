package test

import (
	"novaro-server/model"
	"testing"
)

func TestMakeInvitationCode(t *testing.T) {
	n := 8
	code, err := model.MakeInvitationCode(n)
	if err != nil {
		t.Fatalf("Make invitation code error: %v", err)
	}
	if len(code) != n {
		t.Fatalf("Invitation code length is %v, not %v: %v", len(code), n, code)
	}
}
