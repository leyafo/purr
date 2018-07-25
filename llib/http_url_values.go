package llib

import (
	"net/url"

	lua "github.com/yuin/gopher-lua"
)

/*Export
type Values
    func ParseQuery(query string) (Values, error)
    func (v Values) Add(key, value string)
    func (v Values) Del(key string)
    func (v Values) Encode() string
    func (v Values) Get(key string) string
    func (v Values) Set(key, value string)
*/

const luaURLValuesTypeName = "url_values"

/*
//register to global module
func registerURLValuesType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaURLValuesTypeName)
	L.SetGlobal(luaURLValuesTypeName, mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newURLValues))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), urlValueMethods))
}
*/

func constructURLValuesUD(L *lua.LState, urlValues *url.Values) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = urlValues
	L.SetMetatable(ud, L.GetTypeMetatable(luaURLValuesTypeName))
	return ud
}

// constructor
func newURLValues(L *lua.LState) int {
	ud := constructURLValuesUD(L, &url.Values{})
	L.Push(ud)
	return 1
}

func parseQuery(L *lua.LState) int {
	query := L.CheckString(1)
	urlValues, err := url.ParseQuery(query)
	if err != nil {
		L.ArgError(1, err.Error())
		return 0
	}
	ud := constructURLValuesUD(L, &urlValues)
	L.Push(ud)
	return 1
}

var urlValueMethods = map[string]lua.LGFunction{
	"add":    urlValuesAdd,
	"del":    urlValuesDel,
	"encode": urlValuesEncode,
	"get":    urlValuesGet,
	"set":    urlValuesSet,
}

func checkURLValues(L *lua.LState) *url.Values {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*url.Values); ok {
		return v
	}
	L.ArgError(1, luaURLValuesTypeName+" expected")
	return nil
}

func urlValuesAdd(L *lua.LState) int {
	urlValues := checkURLValues(L)
	key := L.CheckString(2)
	value := L.CheckString(3)
	urlValues.Add(key, value)
	return 0
}
func urlValuesDel(L *lua.LState) int {
	urlValues := checkURLValues(L)
	key := L.CheckString(2)
	urlValues.Del(key)
	return 0
}

func urlValuesEncode(L *lua.LState) int {
	urlValues := checkURLValues(L)
	encodeStr := urlValues.Encode()
	L.Push(lua.LString(encodeStr))
	return 1
}
func urlValuesGet(L *lua.LState) int {
	urlValues := checkURLValues(L)
	key := L.CheckString(2)
	value := urlValues.Get(key)
	L.Push(lua.LString(value))
	return 1
}

func urlValuesSet(L *lua.LState) int {
	urlValues := checkURLValues(L)
	key := L.CheckString(2)
	value := L.CheckString(3)
	urlValues.Set(key, value)

	return 0
}
