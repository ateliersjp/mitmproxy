if ngx.var.request_method ~= "GET" then
    return 1
else
    return 0
end
