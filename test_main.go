package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	lua "github.com/yuin/gopher-lua"
	llib "gitlab.ucloudadmin.com/yafo.li/purr/llib"
)

var luaPath string

func setPath(L *lua.LState) {
	tb := L.GetGlobal("package")
	p, _ := tb.(*lua.LTable)
	path := p.RawGetString("path")
	p.RawSetString("path", lua.LString(fmt.Sprintf("%s;%s", path, luaPath+"/?.lua")))
	L.SetGlobal("package", p)
}

func main() {
	args := os.Args
	luaPath = args[1]
	var luaFile string
	if len(args) > 2 {
		luaFile = args[2]
	}
	L := lua.NewState()
	defer L.Close()

	setPath(L)
	llib.Loader(L)

	config := loadConfiguration(L)
	httpConfig := config.RawGetString("http").(*lua.LTable)
	var hostName string
	hostName = httpConfig.RawGetString("host").String()

	llib.SetHTTPHost(hostName)

	if len(args) > 2 {
		err := L.DoFile(path.Join(luaPath, luaFile))
		if err != nil {
			panic(err.Error())
		}
	} else {
		dir, err := os.Open(luaPath)
		if err != nil {
			panic(err.Error())
		}
		names, err := dir.Readdirnames(0)
		if err != nil {
			panic(err.Error())
		}
		for _, name := range names {
			if strings.LastIndex(name, "_test.lua") != -1 {
				err = L.DoFile(path.Join(luaPath, name))
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}

//LoadSecretKeys ...
func loadConfiguration(L *lua.LState) *lua.LTable {
	err := L.DoFile(path.Join(luaPath, "config.lua"))
	if err != nil {
		panic(err.Error())
	}
	config := L.Get(-1).(*lua.LTable)
	L.Pop(1)
	return config
}
