package llib

import (
	"bytes"
	"io/ioutil"

	"fmt"
	"net/http"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

//HTTPHost is set server host
var HTTPHost string

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

func doRequest(method string, L *lua.LState) int {
	method = strings.ToUpper(method)
	path := L.CheckString(1)
	ud := L.CheckUserData(2)

	header, ok := ud.Value.(http.Header)
	if !ok {
		L.ArgError(1, "header expected")
		return 0
	}
	body := L.CheckTable(3)
	var req *http.Request
	if method != "GET" {
		jsonBody, err := valueToJSON(body)
		if err != nil {
			panic(err.Error)
		}
		req, _ = http.NewRequest(method, HTTPHost+path, bytes.NewReader(jsonBody))
	} else {
		req, _ = http.NewRequest(method, HTTPHost+path, nil)
	}
	req.Header = header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	//return status code
	L.Push(lua.LNumber(resp.StatusCode))

	//return response header
	ud = L.NewUserData()
	ud.Value = header
	L.SetMetatable(ud, L.GetTypeMetatable(luaHeaderTypeName))
	L.Push(ud)

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
