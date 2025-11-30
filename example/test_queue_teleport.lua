-- Test script untuk queue_on_teleport dengan HttpGet
local Players = game:GetService("Players")

-- Ini seharusnya TIDAK di-bundle (dalam function call)
queue_on_teleport("loadstring(game:HttpGet('https://example.com/another-loader.lua'))()")

-- Ini juga seharusnya TIDAK di-bundle (dalam function call dengan syn)
syn.queue_on_teleport("loadstring(game:HttpGet('https://example.com/another-loader.lua'))()")

-- Local require seharusnya tetap di-bundle
local utils = require('./myscript/utils/fancy_print')

utils.print_fancy("Script loaded!")

print("Queue on teleport configured")
