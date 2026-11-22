local M = {}

function M.get_settings_fields(ctx)
    return {
        title = {
            type = "str",
            required = true,
            label = "Название привычки",
            max_length = 100,
        },
        target = {
            type = "int",
            required = true,
            label = "Цель в день",
            default = 1,
            min = 1,
            max = 100,
        },
    }
end

function M.on_quest_create(ctx, params)
    local now = ctx.now
    local day = 24 * 60 * 60

    return {
        {
            type = "create_quest_state",
            starts_at = now,
            ends_at   = now + day,
            internal  = { streak = 0, best_streak = 0, note = "" },
        },
    }
end

return M