if ngx.ctx.is_html == nil then
    ngx.ctx.is_html = ngx.re.find(ngx.arg[1], [[<!DOCTYPE]], "joi") and true or false
end

if ngx.ctx.html_rewriter == nil and ngx.ctx.is_html then
    local lolhtml = require "lolhtml"
    ngx.ctx.html_rewriter = coroutine.wrap(function(chunk)
        local buffered = ""
        local builder = lolhtml.new_rewriter_builder()
        builder:add_element_content_handlers({
            selector = lolhtml.new_selector("script, link[as=script]"),
            element_handler = function(el)
                el:remove()
            end
        })
        builder:add_element_content_handlers({
            selector = lolhtml.new_selector("noscript"),
            element_handler = function(el)
                el:remove_and_keep_content()
            end,
        })
        local rewriter = lolhtml.new_rewriter({
            builder = builder,
            sink = function(output)
                buffered = buffered..output
            end,
        })
        while chunk do
            rewriter:write(chunk)
            chunk, buffered = buffered, ""
            chunk = coroutine.yield(chunk)
        end
        rewriter:close()
        return buffered
    end)
end

if ngx.ctx.html_rewriter then
    ngx.arg[1] = ngx.ctx.html_rewriter(ngx.arg[1])
end
