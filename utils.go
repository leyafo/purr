package purr

import (
	"encoding/hex"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func printLuaVariable(L *lua.LState) int {
	length := L.GetTop()
	for i := 1; i <= length; i++ {
		printLV(0, L.Get(i))
		fmt.Println()
	}
	return 0
}

func hexStr(L *lua.LState) int {
	src := L.CheckString(1)
	retStr := hex.EncodeToString([]byte(src))
	L.Push(lua.LString(retStr))
	return 1
}

func testChecking(L *lua.LState) int {
	v1 := L.CheckAny(1)
	v2 := L.CheckAny(2)
	dbg, ok := L.GetStack(1)
	if ok {
		_, err := L.GetInfo("Sl", dbg, lua.LNil)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	var v3, info string
	if L.GetTop() == 3 {
		v3 = L.CheckString(3)
		info = fmt.Sprintf("%s:%d [%s]", dbg.Source, dbg.CurrentLine, v3)
	} else {
		info = fmt.Sprintf("%s:%d ", dbg.Source, dbg.CurrentLine)
	}
	if v1.Type() != v2.Type() || v1 != v2 {
		//print with red color
		fmt.Printf("\x1b[0;31m%s %#v != %#v\tFailed\x1b[0m\n", info, v1, v2)
		L.Push(lua.LBool(false))
	} else {
		//print with green color
		fmt.Printf("\x1b[0;32m%s %#v == %#v\tPASS\x1b[0m\n", info, v1, v2)
		L.Push(lua.LBool(true))
	}
	return 1
}

func printLV(level int, v lua.LValue) {
	switch v.Type() {
	case lua.LTTable:
		tb, _ := v.(*lua.LTable)
		tb.ForEach(func(key, value lua.LValue) {
			for i := 0; i < level; i++ {
				fmt.Print("  ")
			}
			printLV(level+1, key)
			fmt.Print("=> ")
			if value.Type() == lua.LTTable {
				fmt.Println()
				printLV(level+1, value)
			} else {
				printLV(level+1, value)
				fmt.Println()
			}
		})
	case lua.LTUserData:
		fmt.Printf("%v ", v.(*lua.LUserData).Value)
	default:
		fmt.Printf("%v ", v)
	}
}
