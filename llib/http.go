package llib

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"net/http"
	"net/url"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func httpLoader(L *lua.LState) int {
	//http exports are static functions
	mod := L.SetFuncs(L.NewTable(), httpExports)

	//The header module is belongs http module
	headerModule := L.NewTypeMetatable(luaHeaderTypeName)
	L.SetField(headerModule, "new", L.NewFunction(newHeader))
	L.SetField(headerModule, "__index", L.SetFuncs(L.NewTable(), headerMethods))
	L.SetField(mod, luaHeaderTypeName, headerModule)

	urlValuesModule := L.NewTypeMetatable(luaURLValuesTypeName)
	// static methods
	L.SetField(urlValuesModule, "new", L.NewFunction(newURLValues))
	L.SetField(urlValuesModule, "parse_query", L.NewFunction(parseQuery))
	// methods
	L.SetField(urlValuesModule, "__index", L.SetFuncs(L.NewTable(), urlValueMethods))
	L.SetField(mod, luaURLValuesTypeName, urlValuesModule)

	urlModule := L.NewTypeMetatable(luaURLTypeName)
	// static methods
	L.SetField(urlModule, "new", L.NewFunction(newURL))
	L.SetField(urlModule, "parse_rawURL", L.NewFunction(parseRawURL))
	L.SetField(urlModule, "parse_rquestURL", L.NewFunction(parseRequestURI))
	// methods
	L.SetField(urlModule, "__index", L.SetFuncs(L.NewTable(), urlMethods))
	L.SetField(mod, luaURLTypeName, urlModule)

	L.Push(mod)
	return 1
}

var httpExports = map[string]lua.LGFunction{
	"request": httpRequest,
	"get":     httpGet,
	"post":    httpPost,
	"put":     httpPut,
	"head":    httpHead,
	"delete":  httpDelete,
}

func getBody(n int, L *lua.LState) (body []byte) {
	bodyVal := L.CheckAny(n)
	switch bodyVal.Type() {
	case lua.LTString:
		lstr, _ := bodyVal.(lua.LString)
		body = []byte(string(lstr))
		break
	case lua.LTTable:
		jsonBody, err := valueToJSON(bodyVal)
		if err != nil {
			L.ArgError(3, "Marshal json error, msg is: "+err.Error())
			break
		}
		body = jsonBody
	default:
		break
	}
	return nil
}

func doRequest(method string, L *lua.LState) int {
	method = strings.ToUpper(method)

	urlUD := L.CheckUserData(1)
	reqURL, ok := urlUD.Value.(*url.URL)
	if !ok {
		L.ArgError(1, "url.URL expected")
		return 0
	}

	headerUD := L.CheckUserData(2)
	header, ok := headerUD.Value.(http.Header)
	if !ok {
		L.ArgError(2, "http.header expected")
		return 0
	}

	argumentLen := L.GetTop()
	var body []byte
	if argumentLen == 3 {
		body = getBody(3, L)
	}

	urlStr := reqURL.String()
	fmt.Printf("Request URL is: %s, Method is %s\n", urlStr, method)
	req, err := http.NewRequest(method, reqURL.String(), bytes.NewReader(body))
	if err != nil {
		L.RaiseError("Call http.NewRequest error, msg is %s", err.Error())
		return 0
	}
	req.Header = header

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
	responseUD := L.NewUserData()
	responseUD.Value = header
	L.SetMetatable(responseUD, L.GetTypeMetatable(luaHeaderTypeName))
	L.Push(responseUD)

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
