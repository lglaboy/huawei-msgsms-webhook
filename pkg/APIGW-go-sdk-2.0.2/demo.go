package main

import (
	"bytes"
	"fmt"
	core "huaweicloud.com/apig/signer"
	"io/ioutil"
	"net/http"
)

func main() {
	demoAppApigw()
}

func demoAppApigw() {
	s := core.Signer{
		Key:    "apigateway_sdk_demo_key",
		Secret: "apigateway_sdk_demo_secret",
	}
	r, err := http.NewRequest("POST", "https://30030113-3657-4fb6-a7ef-90764239b038.apigw.cn-north-1.huaweicloud.com/app1?a=1&b=2",
		ioutil.NopCloser(bytes.NewBuffer([]byte("foo=bar"))))
	if err != nil {
		fmt.Println(err)
		return
	}

	r.Header.Add("content-type", "application/json; charset=utf-8")
	r.Header.Add("x-stage", "RELEASE")
	s.Sign(r)
	fmt.Println(r.Header)
	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
}
