if ngx.ctx.is_html == nil then
    ngx.ctx.is_html = ngx.re.find(ngx.arg[1], [[<!DOCTYPE]], "joi") and true or false
end

if ngx.ctx.html_rewriter == nil and ngx.ctx.is_html then
    local cjson = require "cjson"
    local lolhtml = require "lolhtml"
    ngx.ctx.html_rewriter = coroutine.wrap(function(chunk)
        local buffered, is_buffered = "", false
        local builder = lolhtml.new_rewriter_builder()
        builder:add_element_content_handlers({
            selector = lolhtml.new_selector("script#mf-state"),
            text_handler = function(chunk)
                buffered = buffered..chunk:get_text()
                is_buffered = chunk:is_last_in_text_node()
            end,
        })
        local rewriter = lolhtml.new_rewriter({
            builder = builder,
            sink = function() end,
        })
        while chunk do
            rewriter:write(chunk)
            chunk = ""
            if is_buffered then
                for k, v in pairs(cjson.decode(buffered)) do
                    chunk = cjson.encode(v.body)
                    break
                end
                buffered, is_buffered = "", false
            end
            chunk = coroutine.yield(chunk)
        end
        rewriter:close()
        return ""
    end)
end

if ngx.ctx.html_rewriter then
    ngx.arg[1] = ngx.ctx.html_rewriter(ngx.arg[1])
end
