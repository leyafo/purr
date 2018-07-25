crypto = require("crypto")
a = "hello world"
encode = crypto.base64Encode(a)
put(encode)
decode = crypto.base64Decode(encode)
check(decode, a)

