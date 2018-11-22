package purr

import (
	"fmt"
	"path"

	lua "github.com/yuin/gopher-lua"
)

//LoadModules all of modules and global functions.
func LoadModules(L *lua.LState) int {
	L.PreloadModule("http", httpLoader)
	L.PreloadModule(luaCryptoModuleName, luaCryptoLoader)
	L.PreloadModule(luaBufferTypeName, luaBufferLoader)

	L.SetGlobal("to_json", L.NewFunction(toJSON))
	L.SetGlobal("from_json", L.NewFunction(fromJSON))
	L.SetGlobal("println", L.NewFunction(printLuaVariable))
	L.SetGlobal("check", L.NewFunction(testChecking))
	L.SetGlobal("expect", L.NewFunction(testChecking))
	L.SetGlobal("hex", L.NewFunction(hexStr))
	L.SetGlobal("sha256sum", L.NewFunction(sha256sum))
	L.SetGlobal("sha1sum", L.NewFunction(sha1sum))
	L.SetGlobal("uuid", L.NewFunction(luaUUID))
	L.SetGlobal("rand_string", L.NewFunction(luaRandString))
	L.SetGlobal("timestamp", L.NewFunction(luaTimestamp))
	L.SetGlobal("fopen", L.NewFunction(lOpenFile))
	L.SetGlobal("fclose", L.NewFunction(lCloseFile))
	L.SetGlobal("fsize", L.NewFunction(lFileSize))
	L.SetGlobal("fbuffer", L.NewFunction(lFileBuffer))

	//for http usage
	L.SetGlobal("PUT", L.NewFunction(httpPut))
	L.SetGlobal("GET", L.NewFunction(httpGet))
	L.SetGlobal("POST", L.NewFunction(httpPost))
	L.SetGlobal("DELETE", L.NewFunction(httpDelete))
	L.SetGlobal("HEAD", L.NewFunction(httpHead))

	return 0
}

func httpLoader(L *lua.LState) int {
	//http exports are static functions
	mod := L.SetFuncs(L.NewTable(), httpExports)

	urlModule := L.NewTypeMetatable(purrURLTypeName)
	// static methods
	L.SetField(urlModule, "new", L.NewFunction(newURL))
	// methods
	L.SetField(urlModule, "__index", L.SetFuncs(L.NewTable(), purrURLMethods))
	L.SetField(mod, purrURLTypeName, urlModule)

	L.Push(mod)
	return 1
}

//RegisterModule registe a lua module into state
func RegisterModule(L *lua.LState, typeName string, staticFuncs, methods map[string]lua.LGFunction) {
	mt := L.NewTypeMetatable(typeName)
	for name, luaFunc := range staticFuncs {
		L.SetField(mt, name, L.NewFunction(luaFunc))
	}
	if methods != nil {
		L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), methods))
	}
}

//SetTestCasesPath for test case path
func SetTestCasesPath(L *lua.LState, scriptPath string) {
	tb := L.GetGlobal("package")
	p, _ := tb.(*lua.LTable)
	luaEnv := p.RawGetString("path")
	scriptPath = path.Join(path.Dir(scriptPath), "?.lua")
	luaEnv = lua.LString(fmt.Sprintf("%s;%s", luaEnv, scriptPath))
	p.RawSetString("path", luaEnv)
	L.SetGlobal("package", p)
}
