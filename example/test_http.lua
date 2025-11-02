-- Test script with HTTP dependency
local httpModule = loadstring(game:HttpGet("https://raw.githubusercontent.com/Roblox/testez/master/LICENSE"))()

print("Loaded HTTP module")
print(httpModule)
