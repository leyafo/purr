package purr

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"hash"

	lua "github.com/yuin/gopher-lua"
)

const luaCryptoModuleName = "crypto"

func luaCryptoLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), cryptoExports)
	L.Push(mod)
	return 1
}

var cryptoExports = map[string]lua.LGFunction{
	"base64encode": base64Encode,
	"base64decode": base64Decode,
	"hmac":         hMAC,
	"sha1sum":      sha1sum,
	"sha256sum":    sha256sum,
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

func hMAC(L *lua.LState) int {
	cryptoType := L.CheckString(1)
	key := []byte(L.CheckString(2))
	var hm hash.Hash
	switch cryptoType {
	case "sha1":
		hm = hmac.New(sha1.New, key)
		break
	case "sha256":
		hm = hmac.New(sha256.New, key)
		break
	default:
		L.ArgError(1, "invalid crypto type")
		return 0
	}
	cryptoDates := L.CheckTable(3)
	tbLen := cryptoDates.Len()
	for i := 1; i <= tbLen; i++ {
		v := cryptoDates.RawGetInt(i)
		hm.Write([]byte(lua.LVAsString(v)))
	}
	retStr := string(hm.Sum(nil))
	L.Push(lua.LString(retStr))

	return 1
}

func sha256sum(L *lua.LState) int {
	src := L.CheckString(1)
	sum := sha256.New()
	sum.Write([]byte(src))
	dst := sum.Sum(nil)
	L.Push(lua.LString(string(dst)))
	return 1
}

func sha1sum(L *lua.LState) int {
	src := L.CheckString(1)
	sum := sha1.New()
	sum.Write([]byte(src))
	dst := sum.Sum(nil)
	L.Push(lua.LString(string(dst)))
	return 1
}
