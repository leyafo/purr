package llib

import (
	"encoding/base64"

	lua "github.com/yuin/gopher-lua"
)

const luaCryptoModuleName = "crypto"

func loadCryptoModule(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), cryptoExports)
	L.Push(mod)
	return 1
}

var cryptoExports = map[string]lua.LGFunction{
	"base64Encode": base64Encode,
	"base64Decode": base64Decode,
}

func base64Encode(L *lua.LState) int {
	str := L.CheckString(1)
	encodeStr := base64.StdEncoding.EncodeToString([]byte(str))
	L.Push(lua.LString(encodeStr))
	return 1
}

func base64Decode(L *lua.LState) int {
	str := L.CheckString(1)
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		L.RaiseError("Call base64Decode error, msg is %s", err.Error())
		return 0
	}
	L.Push(lua.LString(string(bytes)))
	return 1
}
