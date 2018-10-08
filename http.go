package purr

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"net/http"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

var interlHost string

//SetInternalHost set http request host
func SetInternalHost(host string) {
	interlHost = host
}

var httpExports = map[string]lua.LGFunction{
	"request": httpRequest,
	"get":     httpGet,
	"post":    httpPost,
	"put":     httpPut,
	"head":    httpHead,
	"delete":  httpDelete,
}

func getBody(n int, L *lua.LState) io.Reader {
	bodyVal := L.CheckAny(n)
	switch bodyVal.Type() {
	case lua.LTString:
		lstr, _ := bodyVal.(lua.LString)
		return strings.NewReader(string(lstr))
	case lua.LTTable:
		jsonBody, err := valueToJSON(bodyVal)
		if err != nil {
			L.ArgError(3, "Marshal json error, msg is: "+err.Error())
			break
		}
		return bytes.NewReader(jsonBody)
	case lua.LTUserData:
		ud, ok := bodyVal.(*lua.LUserData)
		if !ok {
			L.ArgError(3, "File object expected. ")
		}
		file, ok := ud.Value.(*os.File)
		if !ok {
			L.ArgError(3, "File object expected. ")
		}
		b, err := ioutil.ReadAll(file)
		if err != nil {
			L.ArgError(3, "File can not read")
		}
		return bytes.NewBuffer(b)
	}
	return nil
}

func doRequest(method string, L *lua.LState) int {
	method = strings.ToUpper(method)

	urlStr := strings.TrimSpace(L.CheckString(1))
	urlStr = interlHost + urlStr

	fmt.Printf("Request URL is: %s, Method is %s\n", urlStr, method)
	var req *http.Request
	var err error
	if L.GetTop() >= 2 {
		tb := L.CheckTable(2)
		tb.ForEach(func(k lua.LValue, val lua.LValue) {
			req.Header.Add(k.String(), val.String())
		})
	}

	if L.GetTop() == 3 {
		body := getBody(3, L)
		req, err = http.NewRequest(method, urlStr, body)
	} else {
		req, err = http.NewRequest(method, urlStr, nil)
	}

	if err != nil {
		L.RaiseError("Call http.NewRequest error, msg is %s", err.Error())
		return 0
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		L.RaiseError("Call Do request error, msg is %s", err.Error())
		return 0
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		L.RaiseError("Call ioutil.ReadAll error, msg is %s", err.Error())
		return 0
	}

	//return status code
	L.Push(lua.LNumber(resp.StatusCode))

	//return response header
	var responseHeader = new(lua.LTable)
	for k, vs := range resp.Header {
		str := k + ": "
		for i, v := range vs {
			if i != 0 {
				str += "; " + v
			} else {
				str += v
			}
		}
		responseHeader.RawSetString(k, lua.LString(str))
	}
	L.Push(responseHeader)

	//resturn response body
	L.Push(lua.LString(string(respBody)))
	return 3
}

func httpRequest(L *lua.LState) int {
	method := L.CheckString(1)
	L.Remove(1)
	return doRequest(method, L)
}

func httpGet(L *lua.LState) int {
	return doRequest("GET", L)
}

func httpPost(L *lua.LState) int {
	return doRequest("POST", L)
}

func httpPut(L *lua.LState) int {
	return doRequest("PUT", L)
}

func httpHead(L *lua.LState) int {
	return doRequest("HEAD", L)
}

func httpDelete(L *lua.LState) int {
	return doRequest("DELETE", L)
}
