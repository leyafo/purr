crypto = require("crypto")
a = "hello world"
encode = crypto.base64encode(a)
println(encode)
decode = crypto.base64decode(encode)
check(decode, a)