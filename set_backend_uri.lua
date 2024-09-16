local uri = ngx.re.sub(ngx.var.backend_uri, "^/[^/]+", "", "jo")
return uri
