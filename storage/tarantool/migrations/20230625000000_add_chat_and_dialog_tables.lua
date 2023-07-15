local function up()
    local chat = box.schema.create_space('chat', {
        format = {
            { name = 'uuid',          type = 'uuid',     is_nullable = false },
            { name = 'participant_1', type = 'unsigned', is_nullable = false },
            { name = 'participant_2', type = 'unsigned', is_nullable = false }
        },
        if_not_exists = true,
    })
    chat:create_index('primary', {
        type = 'tree',
        parts = { 'uuid' },
        if_not_exists = true,
    })
    chat:create_index('chat_participants_index', {
        type = 'tree',
        parts = { 'participant_1', 'participant_2' },
        if_not_exists = true,
    })

    box.schema.func.create('check_text', {
        language = 'LUA',
        is_deterministic = true,
        if_not_exists = true,
        body = 'function(text) return (string.len(text) > 0) end'
    })

    local dialog = box.schema.create_space('dialog', {
        format = {
            { name = 'uuid',       type = 'uuid',     is_nullable = false },
            { name = 'chat_id',    type = 'uuid',     is_nullable = false },
            { name = 'from',       type = 'uuid',     is_nullable = false },
            { name = 'to',         type = 'uuid',     is_nullable = false },
            { name = 'text',       type = 'string',   is_nullable = false, constraint = 'check_text' },
            { name = 'created_at', type = 'datetime', is_nullable = false },
            { name = 'updated_at', type = 'datetime', is_nullable = false }
        },
        if_not_exists = true,
    })
    dialog:create_index('primary', {
        type = 'tree',
        parts = { 'uuid' },
        if_not_exists = true,
    })
    dialog:create_index('chat_id_index', {
        type = 'tree',
        parts = { 'chat_id' },
        if_not_exists = true,
    })

    box.schema.func.create('add_chat', {
        language = 'LUA',
        if_not_exists = true,
        body = [[
            function(participant_1, participant_2)
                local uuid = require('uuid')
                return box.space.chat:insert({ uuid.new(), participant_1, participant_2 })
            end
        ]]
    })

    box.schema.func.create('get_chat', {
        language = 'LUA',
        if_not_exists = true,
        body = [[
            function(participant_1, participant_2)
                return box.space.chat.index.chat_participants_index:get({ participant_1, participant_2 })
            end
        ]]
    })

    box.schema.func.create('send_dialog', {
        language = 'LUA',
        if_not_exists = true,
        body = [[
            function(chat_id, from, to, text, created_at, updated_at)
                local uuid = require('uuid')
                local datetime = require('datetime')
                return box.space.dialog:insert({ uuid.new(), chat_id, from, to, text, created_at, updated_at })
            end
        ]]
    })

    box.schema.func.create('list_dialog', {
        language = 'LUA',
        if_not_exists = true,
        body = [[
            function(chat_id, offset, limit)
                return box.space.dialog.index.chat_id_index:select({ chat_id }, { offset = offset, limit = limit })
            end
        ]]
    })
end

box.once('20230625000000_add_chat_and_dialog_tables', up)
