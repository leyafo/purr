package llib

import (
	"net/url"

	lua "github.com/yuin/gopher-lua"
)

/*Export
type URL
    func Parse(rawurl string) (*URL, error)
    func ParseRequestURI(rawurl string) (*URL, error)
    func (u *URL) EscapedPath() string
    func (u *URL) Hostname() string
    func (u *URL) IsAbs() bool
    func (u *URL) MarshalBinary() (text []byte, err error)
    func (u *URL) Parse(ref string) (*URL, error)
    func (u *URL) Port() string
    func (u *URL) Query() Values
    func (u *URL) RequestURI() string
    func (u *URL) ResolveReference(ref *URL) *URL
    func (u *URL) String() string
    func (u *URL) UnmarshalBinary(text []byte) error
*/
var urlMethods = map[string]lua.LGFunction{
	"escaped_path":      urlEscapedPath,
	"hostname":          urlHostname,
	"is_abs":            urlIsAbs,
	"marshal_binary":    urlMarshalBinary,
	"parse":             urlParse,
	"port":              urlPort,
	"query":             urlQuery,
	"requestURI":        urlRequestURI,
	"resolve_reference": urlResolveReference,
	"string":            urlString,
	"unmarshal_binary":  urlUnmarshalBinary,
	"set_query":         urlSetQuery,
	"set_path":          urlSetPath,
}

const luaURLTypeName = "url"

func constructURLUD(L *lua.LState, url *url.URL) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = url
	L.SetMetatable(ud, L.GetTypeMetatable(luaURLTypeName))
	return ud
}

//static methods
func newURL(L *lua.LState) int {
	reqURL := &url.URL{}
	host := internalHTTPHost
	if L.GetTop() >= 1 {
		scheme := L.CheckString(1)
		reqURL.Scheme = scheme
		if L.GetTop() >= 2 {
			reqURL.Host = L.CheckString(2)
		} else {
			reqURL.Host = host
		}
	}

	ud := constructURLUD(L, reqURL)
	L.Push(ud)
	return 1
}

func parseRawURL(L *lua.LState) int {
	rawurl := L.CheckString(1)
	urlObj, err := url.Parse(rawurl)
	if err != nil {
		L.ArgError(1, err.Error())
		return 0
	}
	ud := constructURLUD(L, urlObj)
	L.Push(ud)
	return 1
}

func parseRequestURI(L *lua.LState) int {
	rawurl := L.CheckString(1)
	urlObj, err := url.ParseRequestURI(rawurl)
	if err != nil {
		L.ArgError(1, err.Error())
		return 0
	}
	ud := constructURLUD(L, urlObj)
	L.Push(ud)
	return 1
}

func checkURL(L *lua.LState) *url.URL {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*url.URL); ok {
		return v
	}
	L.ArgError(1, luaURLTypeName+" expected")
	return nil
}

func urlEscapedPath(L *lua.LState) int {
	urlObj := checkURL(L)
	retStr := urlObj.EscapedPath()
	L.Push(lua.LString(retStr))
	return 1
}

func urlHostname(L *lua.LState) int {
	urlObj := checkURL(L)
	retStr := urlObj.Hostname()
	L.Push(lua.LString(retStr))
	return 1
}

func urlIsAbs(L *lua.LState) int {
	urlObj := checkURL(L)
	retbool := urlObj.IsAbs()
	L.Push(lua.LBool(retbool))
	return 1
}

func urlMarshalBinary(L *lua.LState) int {
	urlObj := checkURL(L)
	bytes, err := urlObj.MarshalBinary()
	if err != nil {
		L.RaiseError("Call MarshalBinary error, msg is %s", err.Error())
		return 0
	}
	L.Push(lua.LString(string(bytes)))
	return 1
}

func urlParse(L *lua.LState) int {
	urlObj := checkURL(L)
	ref := L.CheckString(2)
	urlRet, err := urlObj.Parse(ref)
	if err != nil {
		L.RaiseError("Call Parse error, msg is %s", err.Error())
		return 0
	}

	ud := constructURLUD(L, urlRet)
	L.Push(ud)
	return 1
}

func urlPort(L *lua.LState) int {
	urlObj := checkURL(L)
	port := urlObj.Port()
	L.Push(lua.LString(port))
	return 1
}

func urlQuery(L *lua.LState) int {
	urlObj := checkURL(L)
	urlValues := urlObj.Query()
	ud := constructURLValuesUD(L, &urlValues)
	L.Push(ud)
	return 1
}

func urlRequestURI(L *lua.LState) int {
	urlObj := checkURL(L)
	uriRet := urlObj.RequestURI()
	L.Push(lua.LString(uriRet))
	return 1
}

func urlResolveReference(L *lua.LState) int {
	urlObj := checkURL(L)
	ud := L.CheckUserData(2)
	urlRef, ok := ud.Value.(*url.URL)
	if !ok {
		L.ArgError(1, luaURLTypeName+" expected")
		return 0
	}
	urlRet := urlObj.ResolveReference(urlRef)
	udRet := constructURLUD(L, urlRet)
	L.Push(udRet)
	return 1
}

func urlString(L *lua.LState) int {
	urlObj := checkURL(L)
	strRet := urlObj.String()
	L.Push(lua.LString(strRet))
	return 1
}

func urlUnmarshalBinary(L *lua.LState) int {
	urlObj := checkURL(L)
	text := L.CheckString(2)
	err := urlObj.UnmarshalBinary([]byte(text))
	if err != nil {
		L.RaiseError("Call UnmarshalBinary error, msg is %s", err.Error())
	}
	return 0
}

func urlSetQuery(L *lua.LState) int {
	urlObj := checkURL(L)
	urlObj.RawQuery = L.CheckString(2)
	return 0
}

func urlSetPath(L *lua.LState) int {
	urlObj := checkURL(L)
	urlObj.Path = L.CheckString(2)
	return 0
}
