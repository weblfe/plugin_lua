-- 引入 migrate
migrate = require("migrate")

-- 队列构造table 模块
create_queue_info_table = {}

create_queue_info_table.table = "{{%queue_info}}"
create_queue_info_table.connection = migrate.connDefault("default")

-- 获取migrate db
function create_queue_info_table.getDb()
    local self = create_queue_info_table
    return migrate.connection(self.connection)
end

-- 数据迁移
function create_queue_info_table.safeUp()
    local columns = {}
    local self = create_queue_info_table
    local db = self.getDb()
    -- 备注
    local tableComment = "队列信息记录表"
    -- columns 表字段定义
    columns["id"] = db.pk().comment("id")
    columns["name"] = db.string(100).comment("队列名")
    columns["status"] = db.tinyint(1).default(1).comment("队列状态")
    columns["appid"] = db.string(100).comment("应用appid")
    columns["type"] = db.string(20).comment("应用类型(mqtt,amqp,native,redis)")
    columns["consumer_max_num"] = db.integer().default(1).comment("消费协程数量限制")
    columns["properties"] = db.text().nullable().comment("队列配置属性")
    columns["comment"] = db.string(100).comment("队列备注信息")
    columns["created_at"] = db.datetime().comment("创建时间")
    columns["updated_at"] = db.datetime().comment("更新时间")
    -- db.addColumn(self.table,"deleted_at",db.string().nullable().comment("删除时间").after("created_at"))
    -- local comment = string.format("comment(\"%s\")", tableComment)
    db.createTable(self.table, columns, db.comment(tableComment))
    -- 构建索引
    db.createIndex(self.table, "idx_queue", { "name", "user" })
end

-- 回滚
function create_queue_info_table.safeDown()
    local self = create_queue_info_table
    local db = self.getDb()
    db.dropTable(self.table)
    -- db.dropIndex(create_queue_info_table.table,"idx_queue")
    -- db.dropColumn(create_queue_info_table.table,"deleted_at")
end

return create_queue_info_table