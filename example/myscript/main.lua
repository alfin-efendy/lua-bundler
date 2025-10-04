-- Main entry point
local EzUI = loadstring(game:HttpGet('https://raw.githubusercontent.com/alfin-efendy/ez-rbx-ui/refs/heads/main/ui.lua'))()

local CoreModule = require('../core.lua')
local FancyPrint = require('./utils/fancy_print.lua')
local UI = require('./ui.lua')

local window = EzUI.CreateWindow({
    Name = "Lua Bundler Example",
    Width = 700,
    Height = 400,
    Opacity = 0.9,
    AutoAdapt = true
})

UI:Init(window, CoreModule, FancyPrint)
UI:CreateTab()