-- main.lua
-- План: ежедневная привычка.
-- Хук вызывается хостом, когда пользователь запускает новый квест.

local M = {}

--- @param ctx table  -- контекст хоста: ctx.now (unix timestamp), ctx.config, ctx.log
--- @param params table -- параметры запуска (пока пусто, задел на будущее)
--- @return table -- QuestCreateAction: { type = "create_quest_state", ... }
function M.on_quest_create(ctx, params)
    local now = ctx.now
    local day = 24 * 60 * 60

    return {
        type = "create_quest_state",
        starts_at = now,
        ends_at   = now + day,
        internal = {
            streak      = 0,
            best_streak = 0,
            note        = "",
        },
    }
end

return M