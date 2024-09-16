local proxy_type = "/"
local uri = ngx.re.gsub(ngx.var.backend_uri, [[^/proxy:([^/]*)]], function(m)
    if m[1] ~= "" and m[1] ~= "pass" then
        proxy_type = "/"..m[1].."/"
    end
    return "/s"
end, "jo")

uri = ngx.re.gsub(uri, [[^/s/(https?)://(\w+\.)([^/]+)/(https?)://(\w+\.)([^/]+)]], function(m)
    local uri = "/s/"..m[4].."%3A%2F%2F"..m[5]..m[6].."/"..m[1].."%3A%2F%2F"..m[2]..m[3]
    if m[2] == "www" and m[5] == "www" then
        uri = uri.."/s/"..m[6].."/"..m[3]
    elseif m[2] == m[5] and ngx.re.find(m[3], [[\.]], "jo") and ngx.re.find(m[6], [[\.]], "jo") then
        uri = uri.."/s/"..m[6].."/"..m[3]
    else
        uri = uri.."/s/"..m[5]..m[6].."/"..m[2]..m[3]
    end
    uri = uri..proxy_type..m[4].."://"..m[5]..m[6]
    return uri
end, "jo")
return uri
