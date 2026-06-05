function handle_rate_limit(req, resp)
    print("Lua rate limit handler called")
    if resp then
        print("StatusCode: " .. resp:statusCode())
    else
        print("Response object is nil")
    end
    
    if resp:statusCode() == 429 then
        print("Status is 429, setting body")
        resp:headers("Content-Type", "application/json")
        resp:body('{"http_status_code": 429, "message": "Rate limit exceeded. Please try again later.", "error_code": "RATE_LIMIT_EXCEEDED"}')
    end
end

