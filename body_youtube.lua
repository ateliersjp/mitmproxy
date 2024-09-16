if ngx.ctx.is_html == nil then
    ngx.ctx.is_html = ngx.re.find(ngx.arg[1], [[<!DOCTYPE]], "joi") and true or false
end

if ngx.ctx.html_rewriter == nil and ngx.ctx.is_html then
    local cjson = require "cjson"
    local lolhtml = require "lolhtml"
    ngx.ctx.html_rewriter = coroutine.wrap(function(chunk)
        local buffered, is_buffered, is_last = "", true, false
        local builder = lolhtml.new_rewriter_builder()
        builder:add_element_content_handlers({
            selector = lolhtml.new_selector("body script"),
            text_handler = function(chunk)
                if is_buffered and not is_last then
                    buffered = buffered..chunk:get_text()
                    is_last = chunk:is_last_in_text_node()
                end
            end,
        })
        local rewriter = lolhtml.new_rewriter({
            builder = builder,
            sink = function() end,
        })
        while chunk do
            rewriter:write(chunk)
            chunk = ""
            if is_buffered and is_last then
                buffered = string.match(buffered, "var ytInitialPlayerResponse = (%{.*%});")
                if buffered then
                    local streaming_data = cjson.decode(buffered).streamingData
                    if ngx.var.media_type ~= "" then
                        local adaptive_formats = {}
                        setmetatable(adaptive_formats, cjson.array_mt)
                        for i, v in ipairs(streaming_data.adaptiveFormats) do
                            if v.mimeType:find(ngx.var.media_type) then
                                table.insert(adaptive_formats, v)
                            end
                        end
                        streaming_data.adaptiveFormats = adaptive_formats
                    end
                    chunk = cjson.encode(streaming_data)
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
