local proxy_type = "/"
local uri = ngx.re.gsub(ngx.var.backend_uri, [[^/proxy:([^/]*)]], function(m)
    if m[1] ~= "" and m[1] ~= "pass" then
        proxy_type = "/"..m[1].."/"
    end
    return "/s"
end, "jo")

uri = ngx.re.gsub(uri, [[^/s/(https?)://((?:www\.)?([^/]+))/(https?)://((?:www\.)?([^/]+))]], "/s/$1%3A%2F%2F$2/$4%3A%2F%2F$5/s/$3/$6"..proxy_type.."$1://$2", "jo")
return uri
