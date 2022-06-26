package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/Romern/redactionschemes"
	"strings"
	"sync"
	"syscall/js"
)

func SplitIntoWords(s string) *redactionschemes.PartitionedData {
	var out redactionschemes.PartitionedData
	for _, v := range strings.Split(s, " ") {
		out = append(out, []byte(string(v)))
	}
	return &out
}

var wg sync.WaitGroup // 1
type fn func(this js.Value, args []js.Value) (any, error)

var (
	jsErr     js.Value = js.Global().Get("Error")
	jsPromise js.Value = js.Global().Get("Promise")
)

func asyncFunc(innerFunc fn) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		handler := js.FuncOf(func(_ js.Value, promFn []js.Value) any {
			resolve, reject := promFn[0], promFn[1]

			go func() {
				defer func() {
					if r := recover(); r != nil {
						reject.Invoke(jsErr.New(fmt.Sprint("panic:", r)))
					}
				}()

				res, err := innerFunc(this, args)
				if err != nil {
					reject.Invoke(jsErr.New(err.Error()))
				} else {
					resolve.Invoke(res)
				}
			}()

			return nil
		})

		return jsPromise.New(handler)
	})
}

var nodejs = true

func main() { //js.FuncOf

	document := js.Global().Get("document")
	if document.IsUndefined() {
		nodejs = true
	} else {
		nodejs = false
	}

	js.Global().Set("test3", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		wg.Add(1)
		defer wg.Done()
		return js.ValueOf("ok")
	}))
	js.Global().Set("test2", asyncFunc(func(this js.Value, args []js.Value) (any, error) {
		wg.Add(1)
		defer wg.Done()
		return "ok", nil
	}))
	js.Global().Set("test", asyncFunc(func(this js.Value, args []js.Value) (any, error) {
		wg.Add(1)
		defer wg.Done()

		return test(args[0].String()), nil
	}))
	wg.Wait()
	if !nodejs { // for browser hold run with await "await run() " but not for nodejs
		//select {}
		ch := make(chan struct{})
		<-ch
	}
	//important! for NODEJS print a msg
	fmt.Println("main finished")
}

func test(input_string string) string {

	var c string
	//output:= map[string]string{}
	//args := os.Args
	if len(input_string) == 0 {
		input_string = "This is a test"
	}

	//argCount := len(args[1:])

	//if argCount > 0 {
	//	input_string = string(args[1])
	//}
	//Curve := "{\"Curve\":{\"P\":115792089210356248762697446949407573530086143415290314195533631308867097853951,\"N\":115792089210356248762697446949407573529996955224135760342422259061068512044369,\"B\":41058363725152142129326129780047268409114441015993725554835256314039467401291,\"Gx\":48439561293906451759052585252797914202762949526041747995844080717082404635286,\"Gy\":36134250956749795798585127919587881956611106672985015071877198253568414405109,\"BitSize\":256,\"Name\":\"P-256\"},\"X\":67388459212115845247207701922531519802641110201111896421354254139385779119286,\"Y\":44896789494410232385848926874344592323430414432693281886513413348049879198856,\"D\":52236789956085535229552633414674581936149017638830956144882376308958343958264}"
	//var private_key ecdsa.PrivateKey
	////var private_key crypto.PrivateKey
	////private_key := new(ecdsa.PrivateKey)
	//json.Unmarshal([]byte(Curve), &private_key)

	private_key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//data, _ := json.Marshal(private_key)
	//fmt.Println(string(data))
	//return ""
	//private_key := 63279814311596505439298825162025065669816518537380806320723222594145037816235

	var sig redactionschemes.NaiveSignature

	var priv crypto.PrivateKey = private_key
	//output["Private key: "]=string(private_key.D)

	//fmt.Printf("Private key: %s\n", private_key.D)
	//fmt.Printf("Public key: %s, %s\n", private_key.X, private_key.Y)

	c += fmt.Sprintln("Private key: %s\n", private_key.D)
	c += fmt.Sprintln("Public key: %s, %s\n", private_key.X, private_key.Y)

	//fmt.Printf("Private key: %s\n", private_key.D)
	data_input := SplitIntoWords(input_string)

	incorrect_data := SplitIntoWords(input_string + "incorrect!")

	redacted_data := []int{0, 1} // Remove first and second word

	//fmt.Printf("\nData input: %s\n", input_string)
	//fmt.Printf("Data input: %v\n", data_input)
	//fmt.Printf("Redacted data: %v\n", redacted_data)

	c += fmt.Sprintln("\nData input: %s\n", input_string)
	c += fmt.Sprintln("Data input: %v\n", data_input)
	c += fmt.Sprintln("Redacted data: %v\n", redacted_data)

	err := sig.Sign(data_input, &priv)

	if err == nil {

		//sig_marshaled, _ := sig.Marshal()
		//fmt.Printf("\nSuccessful signing %x\n", sig.BaseSignature)
		c += fmt.Sprintln("\nSuccessful signing %x\n", sig.BaseSignature)

		err = sig.Verify(data_input)

		if err == nil {

			//fmt.Printf("Successful Verification\n")
			c += fmt.Sprintln("Successful Verification\n")
		}

	}

	newSig, _ := sig.Redact(redacted_data, data_input)

	if err == nil {

		r, _ := newSig.Marshal()
		var v map[string]string
		_ = json.Unmarshal([]byte(r), &v)

		//fmt.Printf("\nSuccessful redaction %x\n", v["BaseSignature"])
		c += fmt.Sprintln("\nSuccessful redaction %x\n", v["BaseSignature"])

		redacted_strings, _ := data_input.Redact(redacted_data)
		if redacted_strings != nil {

			err = newSig.Verify(redacted_strings)

			if err == nil {

				//fmt.Printf("Successful Redaction\n")
				c += fmt.Sprintln("Successful Redaction\n")
			}

			err = sig.Verify(incorrect_data)

			if err != nil {

				//fmt.Printf("Successful Checking against incorrect data\n")
				c += fmt.Sprintln("Successful Checking against incorrect data\n")
			}
		} else {
			//fmt.Println("redacted_strings is nil")
			c += fmt.Sprintln("redacted_strings is nil")

		}

	}
	return c

}
