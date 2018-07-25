package llib

import (
	"encoding/base64"
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

const luaHeaderTypeName = "header"

var headerMethods = map[string]lua.LGFunction{
	"add":            headerAdd,
	"del":            headerDel,
	"set":            headerSet,
	"get":            headerGet,
	"set_basic_auth": basicAuthSet,
}

func newHeader(L *lua.LState) int {
	header := make(http.Header)
	hk := L.CheckString(1)
	hv := L.CheckString(2)
	header.Add(hk, hv)

	ud := L.NewUserData()
	ud.Value = header
	L.SetMetatable(ud, L.GetTypeMetatable(luaHeaderTypeName))
	L.Push(ud)
	return 1
}

func headerAdd(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(http.Header); ok {
		hk := L.CheckString(2)
		hv := L.CheckString(3)
		v.Add(hk, hv)
	}
	return 0
}

func headerDel(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(http.Header); ok {
		hk := L.CheckString(2)
		v.Del(hk)
	}
	return 0
}

func headerSet(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(http.Header); ok {
		hk := L.CheckString(2)
		hv := L.CheckString(3)
		v.Set(hk, hv)
	}
	return 0
}

func headerGet(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(http.Header); ok {
		hk := L.CheckString(2)
		L.Push(lua.LString(v.Get(hk)))
	} else {
		L.Push(nil)
	}
	return 1
}

func basicAuthSet(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(http.Header); ok {
		email := L.CheckString(2)
		pwd := L.CheckString(3)
		auth := email + ":" + pwd
		v.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	}
	return 0
}
