config = require("config")
http = require("http")
url = config.url

url_values = http.url_values.new()
url_values:add("key", "value")
url_values:add("method", "get")

url:set_query(url_values:encode())
url:set_path("/hello")

local header = http.header.new("content-type", "application/json")
status_code, header, body = http.get(url, header, {})

put(response, header, body)
check(status_code, 200)
