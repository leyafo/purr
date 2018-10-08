# lua 文档
purr 里面的 go 实现没有太多暴露的接口，你可以查看 [godoc]() 获取这些接口的文档。本文档主要是针对 lua 代码辅助文档。

## http 模块
http 包装了 url 和 url values, 为了方便 http 的 header 直接用 lua 的 table 来代替。

### GET POST PUT DELETE HEAD
为了方便，你可以直接调用这几个全局方法，它们分别对应的是 http 请求的 method。它们接受的参数是 url, header, body。返回的结果是 status code, header, body。  
传入的 body 接受 lua.LTable, lua.LString, FILE(通过 fopen 打开文件得到) 类型。  
返回的 body 是 lua.LString。
header 都是 lua.LTable，里面是 key/value 形式的键值对。

### http.url 模块
url 模块主要是为了操作 url query 而包装的 lua 模块。为了简洁，你可以在调用 GET POST 方法时，直接传入 `?key=val&key1=val1` 这样的 query string。当碰到需要处理特殊字符的情况时，你需要调用这个模块的相关方法来处理。  

**http.url.new(path)**  
path 为 string 类型。  
构造一个 url 里面可以传入 path。  

**query_add(key, value)**  
key,value 为 string 类型。  
新增一个 query 字段。  

**query_del(key)**  
key 为 string 类型。  
删除一个 query 字段。  

**query_set(key,value)**  
key,value 为 string 类型。  
修改一个 query 字段。  

**query_get(key)value**  
key,value 为 string 类型。  
获取一个 query 字段的 value。  

**encode()encoded_str**  
encoded_str 为 string 类型。  
对 query 进行编码输出。  

##crypt 模块  
crypt 模块里面集中了一些 authorization 需要用到的加密算法，你可以使用这里面的方法对请求进行加密。  

**base64encode(str)encodedstr**  
对 str 进行 base64 encode。  

**base64decode(encodedstr)decodedstr**    
接收一个 base64 string 返回 decode string。  

**hmac(type, [str1,str2,str3...])encryptstr**  
接收多个 string(以 array 形式传入) 并对其进行对应算法加密。其中 type 支持 `sha1` 和 `sha256`。  

**sha1sum(str)encryptedstr**  
对 str 进行 hmac sha1 加密。  

**sha256sum(str)encryptedstr**  
对 str 进行 hmac sha256 加密。  

##其他辅助函数
**to_json(luatype)str**  
把任何 lua type 转化为 string。

**from_json(json_str)luatype**  
把 json_str 转化为 lua type。

**println(luatype)**  
打印一个 lua 变量，lua table 会格式化打印。

**check(luatype1, luatype2, describe_str)bool**  
**expect(luatype1, luatype2, describe_str)bool**  
检查两个变量是否相等，如果不相等返回 false，同时打印一条红色字符输出的控制台信息。如果相等返回 ture，同时打印一条绿色字符输出的控制台信息。  
describe_str 是打印时附加的输出信息。

**hex(str)str**  
把 str 变成 16 进制输出。

**sha256sum(str)str**  
对 str 进行 hmac sha256 加密

**sha1sum(str)str**  
对 str 进行 hmac sha1 加密

**uuid()str**  
返回一串随机生成的 uuid 字符串。

**rand_string(length)**  
返回一串随机字符串，length 是随机字符串长度。

**timestamp() number**  
返回当前时间的 timestamp，长度是 13 位。

**fopen(filepath)FILE**  
打开一个文件，返回的 FILE 可以直接作为 http body。

**fsize(FILE)number**
获取文件大小

**fclose(FILE)**
关闭打开的文件
