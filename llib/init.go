package llib

import (
	lua "github.com/yuin/gopher-lua"
)

//Loader all of modules and global functions.
func Loader(L *lua.LState) int {
	L.PreloadModule("http", httpLoader)
	L.PreloadModule(luaCryptoModuleName, loadCryptoModule)
	L.SetGlobal("to_json", L.NewFunction(toJSON))
	L.SetGlobal("from_json", L.NewFunction(fromJSON))
	L.SetGlobal("put", L.NewFunction(printLuaVariable))
	L.SetGlobal("check", L.NewFunction(testChecking))

	return 0
}

//SetHTTPHost ...
func SetHTTPHost(host string) {
	HTTPHost = host
}
