http = require("http")

url_values = http.url_values.new()
url_values:add("remoteurl", "https://127.0.0.1:28888/hello")
url_values:add("method", "get")
url_values:add("content-type", "application/json")
query = url_values:encode()
local header = http.header.new("content-type", "application/json")
put(query)
status_code, header, body = http.get("/cb?"..query, header, {})

put(response, header, body)
check(status_code, 200)
