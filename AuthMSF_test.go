package rpc_test

import (
	"testing"

	. "github.com/cSploit/daemon"
)

var r = RPC{Host: "192.168.0.23", Port: 55552}

func TestAuthLogin(t *testing.T) {
	t.Log("Testing MSFRPC authentication")
	rpcAuth := r.AuthLogin("msf", "Xh3BtmUc")
	if rpcAuth.Error == true {
		t.Errorf("rpcAuth authentication failed")
	}
	t.Logf("token received: ", r.Token)
}

/*
func TestAuthTokenList(t *testing.T) {
	t.Log("Testing MSFRPC token list")
	rpcAuthTokenList := r.AuthTokenList()
	t.Logf("token list: ", rpcAuthTokenList)
}

func TestAuthTokenAdd(t *testing.T) {
	t.Log("Testing MSFRPC token add")
	token := "newtoken"
	rpcAuthTokenAdd := r.AuthTokenAdd(token)
	if rpcAuthTokenAdd.MsfResponse == "success" {
		t.Logf("Token: %s successfully added", token)
	}
	t.Errorf("rpcAuthTokenAdd failed")
}

func TestAuthTokenRemove(t *testing.T) {
	t.Log("Testing MSFRPC token list")
	rpcAuthTokenRemove := r.AuthTokenRemove(r.Token)
	if rpcAuthTokenRemove.MsfResponse == "success" {
		t.Logf("Token: %s successfully removed", r.Token)
	}
	t.Errorf("rpcAuthTokenRemove failed")
}
*/
