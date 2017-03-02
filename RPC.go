package rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/ugorji/go/codec"
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

	fmt.Printf("Open content: \n", args)

	var msgpackH codec.MsgpackHandle
	msgpackH.RawToString = true
	var encodedBody []byte
	var encoder *codec.Encoder
	encoder = codec.NewEncoderBytes(&encodedBody, &msgpackH)
	err := encoder.Encode(args)
	if err != nil {
		log.Fatalln("Error encoding arguments:", err)
	}
	fmt.Printf("Encoding content: \n", encodedBody)

	/** http requests handler */
	httpClient := http.Client{}
	httpData := bytes.NewReader(encodedBody)
	httpReq, err := http.NewRequest("POST", "http://"+r.Host+":"+strconv.Itoa(r.Port)+"/api/1.0", httpData)
	if err != nil {
		log.Fatalln("error with httpreq:", err)
	}
	httpReq.Header.Add("Content-Type", "binary/message-pack")
	httpReq.Header.Add("Content-Length", strconv.Itoa(len(encodedBody)))
	fmt.Printf("httpReq: %v\n", httpReq)
	response, err := httpClient.Do(httpReq)
	if err != nil {
		log.Fatalf("Unable to connect to the MSFPRC. Error: %v\n", err)
	}

	/** ... var body contains the data to decode from */
	decodedBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln("reading body failed:", err)
	}
	var out map[string]interface{}
	var h codec.Handle = new(codec.MsgpackHandle)
	var decoder *codec.Decoder = codec.NewDecoderBytes(decodedBody, h)
	err = decoder.Decode(out)
	if err != nil {
		log.Fatalln("unmarshal failed:", err)
	}
	fmt.Printf("return from RPC calling: %v\n", out)
	return out
}
