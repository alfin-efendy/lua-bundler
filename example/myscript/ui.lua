local m = {}

local Window
local Core
local FancyPrint

function m:Init(window, coreModule, fancyPrintModule)
    Window = window
    Core = coreModule
    FancyPrint = fancyPrintModule
end


function m:CreateTab()
    local tab = Window:AddTab({
        Name = "Hello World",
        Icon = "❤️",
    })

    tab:AddButton("Greet", function()
        local message = Core:greet("User")
        FancyPrint:FancyPrint(message)
    end)
end

return m