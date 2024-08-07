ngx.header.content_length = nil
ngx.header.expires = 300
ngx.ctx.is_html = ngx.header.content_type:find("html", 1, true) and true or false
