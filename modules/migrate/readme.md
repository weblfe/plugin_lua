## lua migrate module 

> ### plan <开发流程>

 1. 完成 go-migrate driver (lua://path/file.lua?params=xx)
 
   - [ ] 1.1 migrate lua driver  [数据迁移lua脚本驱动骨架]
   - [ ] 1.2 lua file driver [文件驱动]
   - [ ] 1.3 script parse and eval [ lua 脚本 解析 与执行 ]
   - [ ] 1.4 script transform to sql [ lua 脚本逻辑 转 sql ]
   - [ ] 1.5 migrate init   [migrate 相关初始化]
   - [ ] 1.6 migrate action with sql [sql 与migrate 版本 行为 关联]

 2. 完成 lua migrate 功能模块
    
   - [ ] 2.1 sql builder  [sql 构造器]
   - [ ] 2.2 data types   [数据基本类型]
   - [ ] 2.3 sql expression [sql 复合表达式]
   - [ ] 2.4 migrate functions binds [ migrate API lua 绑定]
   - [ ] 2.5 sql generate and databases driver adapter [ sql 解析生成 与 <query,cmd,syntax,protocol>数据库驱动适配器 ]

 3. 完成 go-migrate cli 和 go-lua-migrate lib deploy
    
   - [ ] 3.1 support  mysql driver (v0.1.0)
   - [ ] 3.2 继续 迭代 (postgreSQL driver)

    