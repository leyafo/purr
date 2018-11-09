package purr

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

//Injection set a lua fuction into lua state
var Injection lua.LGFunction

//RunTestWithServer Running a http backend server and then load the test case script.
func RunTestWithServer(handle http.Handler, scriptPath string) {
	L := lua.NewState()
	defer L.Close()

	SetTestCasesPath(L, scriptPath)
	LoadModules(L)

	ts := httptest.NewServer(handle)
	SetInternalHost(ts.URL)
	defer ts.Close()

	if Injection != nil {
		Injection(L)
	}

	doTest(L, scriptPath)
}

//RunTest without server
func RunTest(httpHost string, scriptPath string) {
	L := lua.NewState()
	defer L.Close()

	SetTestCasesPath(L, scriptPath)
	LoadModules(L)
	SetInternalHost(httpHost)

	if Injection != nil {
		Injection(L)
	}

	doTest(L, scriptPath)
}

func doTest(L *lua.LState, scriptPath string) {
	if strings.Index(scriptPath, ".lua") != -1 {
		err := L.DoFile(scriptPath)
		if err != nil {
			panic(err.Error())
		}
	} else {
		dir, err := os.Open(scriptPath)
		if err != nil {
			panic(err.Error())
		}
		names, err := dir.Readdirnames(0)
		if err != nil {
			panic(err.Error())
		}
		for _, name := range names {
			if strings.LastIndex(name, "_test.lua") != -1 {
				err = L.DoFile(path.Join(scriptPath, name))
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}
