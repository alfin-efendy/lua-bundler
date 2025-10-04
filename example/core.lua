local m = {}

function m:greet(name)
    return "Hello, " .. name .. "!"
end

function m:farewell(name)
    return "Goodbye, " .. name .. "!"
end

return m