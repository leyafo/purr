local _M = {}
local http = require("http")

local url = http.url

_M.url = url.new("http","127.0.0.1:8000")

return _M
