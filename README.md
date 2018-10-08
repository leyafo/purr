# purr
purr 是用来测试 HTTP 服务的 mock 工具。使用 purr 你可以很轻易模拟一个 HTTP 请求，并检查 HTTP Response 是否符合预期返回。它还可以很方便集成到你以后的 golang 项目里测试已有的 http services。

## Getting Started
purr 使用 lua 做为 DSL (Domain Special Language) 来实现 http 请求。这里 lua 使用的实现是 [gopher-lua](https://github.com/yuin/gopher-lua)。

### installation
`go get github.com/leyafo/purr`  
GopherLua supports >= Go1.9.

### usage
使用 purr 做为测试工具有两种使用方法，一种是集成到已有的 golang 项目里，另一种是单独启动一个进程测试任何语言写的 http 项目，一下是这两种方式使用的测试代码。
```go
//main.go
func demoHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, client")
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "standalone" {
		http.HandleFunc("/", demoHandle)
		go func() {
			http.ListenAndServe(":9527", nil)
        }()
        //如果是其他的语言，你可以只传入 host 和测试文件路径就能测试 http 服务。
		purr.RunTest("http://127.0.0.1:9527/", "./")  
	} else {
        //集成到已有的 golang 服务中
		purr.RunTestWithServer(http.HandlerFunc(demoHandle), "./")
	}
}
```

以下是用 lua 写的测试用例：  
```lua
--hello_test.lua
status, header, body = GET("/")
expect(status, 200)
```
更多示例详见 [example](example) 目录。

## Documentation
go 相关的接口文档你可以通过[在线godoc](https://godoc.org/github.com/leyafo/purr)直接查看，lua 相关的文档在 [doc/lua_doc_cn.md](doc/lua_doc_cn.md)。

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.

## Author
李亚夫 - [leyafo](http://www.leyafo.com)