ngx.header.Content_Length = nil
ngx.header.Cache_Control = "max-age=300"
ngx.ctx.is_html = ngx.re.find(ngx.header.Content_Type or "", [[^text/html]], "jo") and true or false
