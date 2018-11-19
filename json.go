package purr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type jsonV struct {
	lua.LValue
}

func (v jsonV) MarshalJSON() ([]byte, error) {
	return valueToJSON(v.LValue)
}

func toJSON(L *lua.LState) int {
	v := L.CheckAny(1)
	bytes, err := valueToJSON(v)
	if err != nil {
		L.ArgError(1, err.Error())
		return 0
	}
	L.Push(lua.LString(string(bytes)))
	return 1
}

func fromJSON(L *lua.LState) int {
	value := L.CheckAny(1)
	var v interface{}
	var err error
	if value.Type() == lua.LTUserData {
		ud, ok := value.(*lua.LUserData)
		if !ok {
			L.ArgError(3, "*bytes.Buffer object expected. ")
		}
		err = json.NewDecoder(ud.Value.(*bytes.Buffer)).Decode(&v)
	} else if value.Type() == lua.LTString {
		jstr := value.String()
		err = json.NewDecoder(strings.NewReader(jstr)).Decode(&v)
	} else {
		L.ArgError(1, "*bytes.Buffer or string expect.")
		return 0
	}
	if err != nil {
		fmt.Println(err.Error())
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(decode(L, v))
	return 1
}

//from https://github.com/layeh/gopher-json/blob/master/json.go#L150
func decode(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(decode(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), decode(L, item))
		}
		return tbl
	}
	return lua.LNil
}

func valueToJSON(v lua.LValue) (data []byte, err error) {
	switch v.Type() {
	case lua.LTNil:
		data, err = json.Marshal(nil)
	case lua.LTBool:
		b, _ := v.(lua.LBool)
		data, err = json.Marshal(bool(b))
	case lua.LTNumber:
		n, _ := v.(lua.LNumber)
		data, err = json.Marshal(float64(n))
	case lua.LTString:
		s, _ := v.(lua.LString)
		data, err = json.Marshal(string(s))
	case lua.LTTable:
		var array []jsonV
		m := make(map[string]jsonV)
		tb, _ := v.(*lua.LTable)
		tb.ForEach(func(k lua.LValue, val lua.LValue) {
			value := jsonV{LValue: val}
			if k.Type() != lua.LTNumber && k.Type() != lua.LTString {
				err = fmt.Errorf("Unspport table key: %s", k)
				return
			}
			if k.Type() == lua.LTNumber {
				array = append(array, value)
			}
			key := fmt.Sprintf("%v", k)
			m[key] = value
		})
		if len(m) == len(array) {
			data, err = json.Marshal(array) //table is general array
		} else {
			data, err = json.Marshal(m) //table is key-value array
		}
	default:
		err = fmt.Errorf("Cannot convert %s", v.String())
	}
	return
}
