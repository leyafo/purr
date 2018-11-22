package purr

import (
	"io/ioutil"
	"os"

	lua "github.com/yuin/gopher-lua"
)

const luaFileTypeName = "FILE"

//lOpenFile open file as read only mode.
func lOpenFile(L *lua.LState) int {
	filePath := L.CheckString(1)
	file, err := os.Open(filePath)
	if err != nil {
		L.ArgError(1, err.Error())
		return 0
	}

	ud := L.NewUserData()
	ud.Value = file
	L.SetMetatable(ud, L.GetTypeMetatable(luaFileTypeName))
	L.Push(ud)
	return 1
}

func lCloseFile(L *lua.LState) int {
	f := checkFile(L)
	f.Close()
	return 0
}

func checkFile(L *lua.LState) *os.File {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*os.File); ok {
		return v
	}
	L.ArgError(1, luaFileTypeName+" expected")
	return nil
}

func lFileSize(L *lua.LState) int {
	f := checkFile(L)
	fi, err := f.Stat()
	if err != nil {
		L.ArgError(1, err.Error())
		return 0
	}
	L.Push(lua.LNumber(fi.Size()))
	return 1
}

func lFileBuffer(L *lua.LState) int {
	f := checkFile(L)
	b, err := ioutil.ReadAll(f)
	if err != nil {
		L.ArgError(1, err.Error())
		return 0
	}

	ud := L.NewUserData()
	ud.Value = b
	L.SetMetatable(ud, L.GetTypeMetatable(luaBufferTypeName))
	L.Push(ud)
	return 1
}
