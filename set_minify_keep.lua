local uri = ngx.var.backend_uri:gsub("^/[^/]+", "/minify:keep")
return uri
