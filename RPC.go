package rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"gopkg.in/vmihailenco/msgpack.v2"
)

type ServerException struct {
	Error        bool
	ErrorCode    uint64
	ErrorMessage string
}

type Result struct {
	MsfResponse string
}

type RPC struct {
	Host, Token string
	LogginState bool
	Port        int
}

func (r *RPC) Call(args ...interface{}) map[string]interface{} {
	if args[0] != "auth.login" {
		args = append(args, 0)
		copy(args[1:], args[1:])
		args[1] = r.Token
	}
	encoded, err := msgpack.Marshal(args)
	if err != nil {
		log.Fatalln("error marshalling:", err)
	}
	httpClient := http.Client{}
	httpData := bytes.NewReader(encoded)
	httpReq, err := http.NewRequest("POST", "http://"+r.Host+":"+strconv.Itoa(r.Port)+"/api/1.0", httpData)
	if err != nil {
		log.Fatalln("error with httpreq:", err)
	}
	httpReq.Header.Add("Content-Type", "binary/message-pack")
	httpReq.Header.Add("Content-Length", strconv.Itoa(len(encoded)))
	// fmt.Printf("httpReq: %v\n", httpReq)
	response, err := httpClient.Do(httpReq)
	if err != nil {
		log.Fatalf("Unable to connect to the MSFPRC. Error: %v\n", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln("reading body failed:", err)
	}
	var out map[string]interface{}
	err = msgpack.Unmarshal(body, &out)
	if err != nil {
		log.Fatalln("unmarshal failed:", err)
	}
	fmt.Printf("return from auth calling: %v\n", out)
	return out
}
