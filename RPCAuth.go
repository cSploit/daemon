package rpc

import "fmt"

func (rpc *RPC) AuthLogin(username, password string) (ret ServerException) {
	res := rpc.Call("auth.login", username, password)
	fmt.Printf("return from auth calling: %v\n", res)
	if res["error"] != nil && res["error"].(bool) == true {
		ret.Error = true
		fmt.Print("AuthLogin error\n")
		ret.ErrorCode = res["error_code"].(uint64)
		fmt.Printf("ErrorCode: %v\n", ret.ErrorCode)
		ret.ErrorMessage = res["error_message"].(string)
		rpc.LogginState = false
	} else {
		rpc.Token = res["token"].(string)
		rpc.LogginState = true
	}
	return
}

func (rpc *RPC) AuthLogout() (ret ServerException) {
	res := rpc.Call("auth.logout", rpc.Token)
	if res["response"] == "success" {
		fmt.Print("AuthLogout error\n")
		fmt.Printf("ErrorCode: %v\n", ret.ErrorCode)
	} else {
		rpc.LogginState = false
	}
	return
}

func (rpc *RPC) AuthTokenAdd(newToken string) (ret Result) {
	res := rpc.Call("auth.token_add", newToken)
	ret.MsfResponse = res["result"].(string)
	return
}

func (rpc *RPC) AuthTokenGenerate() (ret Result) {
	res := rpc.Call("auth.token_generate")
	rpc.Token = res["token"].(string)
	ret.MsfResponse = res["result"].(string)
	return
}

func (rpc *RPC) AuthTokenList() (s []string) {
	res := rpc.Call("auth.token_list")["tokens"].([]interface{})
	for i := 0; i < len(res); i++ {
		s = append(s, res[i].(string))
	}
	return
}

func (rpc *RPC) AuthTokenRemove(token string) (ret Result) {
	res := rpc.Call("auth.token_remove", token)
	ret.MsfResponse = res["result"].(string)
	if ret.MsfResponse == "success" {
		rpc.Token = ""
	}
	return
}
