package purr

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"math/rand"
	"strings"
	"time"
)

//HEXSET for RandString use
const HEXSET = "0123456789abcdef"

//CHARSET for RandString use
const CHARSET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

//RandString ...
func RandString(length uint, randomSet string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = randomSet[seededRand.Intn(len(randomSet))]
	}
	return string(b)
}

func luaUUID(L *lua.LState) int {
	//uuid format 8-4-4-4-12
	uuid := fmt.Sprintf("%s-%s-%s-%s-%s", RandString(8, HEXSET), RandString(4, HEXSET), RandString(4, HEXSET), RandString(4, HEXSET), RandString(12, HEXSET))
	L.Push(lua.LString(strings.ToUpper(uuid)))
	return 1
}

func luaRandString(L *lua.LState) int {
	length := L.CheckInt(1)
	retStr := RandString(uint(length), CHARSET)
	L.Push(lua.LString(retStr))
	return 1
}

func luaTimestamp(L *lua.LState) int {
	now := time.Now().UnixNano()
	t := int64(now / int64(time.Millisecond))
	L.Push(lua.LNumber(t))
	return 1
}
