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
	L.SetGlobal("hex", L.NewFunction(hexStr))
	L.SetGlobal("sha256sum", L.NewFunction(sha256sum))

	return 0
}

//You can set Host as test server
var internalHTTPHost string

//SetHTTPHost ...
func SetHTTPHost(host string) {
	internalHTTPHost = host
}
