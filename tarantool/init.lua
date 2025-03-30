box.cfg{
    listen = 3301
}

box.once('init', function()
    local polls = box.schema.space.create('polls')
    polls:format({
        { name = 'id', type = 'string' },  
        { name = 'data', type = 'string' }
    })
    polls:create_index('primary', {
        type = 'hash',
        parts = { 'id' }
    })
end)
