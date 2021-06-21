box.cfg {
    listen = 3301,
    wal_dir = '/var/lib/tarantool',
    memtx_dir = '/var/lib/tarantool',
}

box.once('bootstrap', function()
    box.schema.user.create('tarantool', { password = 'tarantool' })
    box.schema.user.grant('tarantool', 'read,write,execute', 'universe')

    local space_name = box.schema.space.create(
            'space_test',
            {
                engine = 'memtx',
                format = {
                    { name = 'id', type = 'uuid' },
                },
            }
    )
    space_name:create_index(
            'pk',
            {
                type = 'tree',
                unique = true,
                parts = { 'id' },
            }
    )
end)

api = {}
api.test = function(id)
    return box.space['space_test']:insert({ id })
end

return api