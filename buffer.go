package purr

import (
	"bytes"
	"fmt"
	"io"

	"crypto/md5"

	lua "github.com/yuin/gopher-lua"
)

const luaBufferTypeName = "buffer"

func luaBufferLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), bufferExports)
	L.Push(mod)
	return 1
}

var bufferExports = map[string]lua.LGFunction{
	"size":    bufferSize,
	"get_md5": getBufMD5,
}

func bufferSize(L *lua.LState) int {
	buf := getBufferUserData(L, 1)
	if buf == nil {
		return 0
	}
	L.Push(lua.LNumber(buf.Len()))
	return 1
}

func getBufMD5(L *lua.LState) int {
	buf := getBufferUserData(L, 1)
	if buf == nil {
		return 0
	}
	h := md5.New()
	if _, err := io.Copy(h, buf); err != nil {
		L.ArgError(1, err.Error())
		return 0
	}
	md5Str := fmt.Sprintf("%x", h.Sum(nil))
	L.Push(lua.LString(md5Str))
	return 1
}

func getBufferUserData(L *lua.LState, index int) *bytes.Buffer {
	ud := L.CheckUserData(index)
	buf, ok := ud.Value.(*bytes.Buffer)
	if !ok {
		L.ArgError(index, "bytes.buffer expected!!")
		return nil
	}
	return buf
}
