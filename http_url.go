package purr

import (
	"net/url"

	lua "github.com/yuin/gopher-lua"
)

type purrURL struct {
	*url.URL
	values url.Values
}

var purrURLMethods = map[string]lua.LGFunction{
	"query_add": queryAdd,
	"query_del": queryDel,
	"query_set": querySet,
	"query_get": queryGet,
	"encode":    urlEncode,
}

const purrURLTypeName = "url"

func constructURLUD(L *lua.LState, url *purrURL) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = url
	L.SetMetatable(ud, L.GetTypeMetatable(purrURLTypeName))
	return ud
}

func checkURL(L *lua.LState) *purrURL {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*purrURL); ok {
		return v
	}
	L.ArgError(1, purrURLTypeName+" expected")
	return nil
}

//static methods
func newURL(L *lua.LState) int {
	path := L.CheckString(1)

	pu := new(purrURL)
	pu.values = url.Values{}
	pu.URL = new(url.URL)
	pu.Path = path

	ud := constructURLUD(L, pu)
	L.Push(ud)
	return 1
}

func parseRawURL(L *lua.LState) int {
	urlStr := L.CheckString(1)
	rawURL, err := url.Parse(urlStr)
	if err != nil {
		L.ArgError(1, err.Error())
		return 0
	}
	pu := purrURL{URL: rawURL, values: url.Values{}}
	ud := constructURLUD(L, &pu)
	L.Push(ud)
	return 1
}

func queryAdd(L *lua.LState) int {
	pu := checkURL(L)
	key := L.CheckString(2)
	value := L.CheckString(3)
	pu.values.Add(key, value)
	return 0
}

func queryDel(L *lua.LState) int {
	pu := checkURL(L)
	key := L.CheckString(2)
	pu.values.Del(key)
	return 0
}

func queryGet(L *lua.LState) int {
	pu := checkURL(L)
	key := L.CheckString(2)
	value := pu.values.Get(key)
	L.Push(lua.LString(value))
	return 1
}

func querySet(L *lua.LState) int {
	pu := checkURL(L)
	key := L.CheckString(2)
	value := L.CheckString(3)
	pu.values.Set(key, value)

	return 0
}

func urlEncode(L *lua.LState) int {
	pu := checkURL(L)
	pu.RawQuery = pu.values.Encode()
	L.Push(lua.LString(pu.String()))
	return 1
}
