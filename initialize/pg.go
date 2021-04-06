package initialize

import (
    "context"
    "fmt"
    "gin-web/models"
    "gin-web/pkg/global"
    pg "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/schema"
    "time"
)

// 初始化pg数据库
func Pg() {
    dsn := fmt.Sprintf(
        "postgresql://%s:%s@%s:%d/%s?&%s",
        global.Conf.Pg.Username,
        global.Conf.Pg.Password,
        global.Conf.Pg.Host,
        global.Conf.Pg.Port,
        global.Conf.Pg.Database,
        global.Conf.Pg.Query,
    )
    // 隐藏密码
    showDsn := fmt.Sprintf(
        "postgresql://%s@%s:%d/%s?&%s",
        global.Conf.Pg.Username,
        global.Conf.Pg.Host,
        global.Conf.Pg.Port,
        global.Conf.Pg.Database,
        global.Conf.Pg.Query,
    )
    global.Log.Info("数据库连接DSN: ", showDsn)
    init := false
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Conf.System.ConnectTimeout)*time.Second)
    defer cancel()
    go func() {
        for {
            select {
            case <-ctx.Done():
                if !init {
                    panic(fmt.Sprintf("初始化pg异常: 连接超时(%ds)", global.Conf.System.ConnectTimeout))
                }
                // 此处需return避免协程空跑
                return
            }
        }
    }()
    db, err := gorm.Open(pg.Open(dsn), &gorm.Config{
        // 禁用外键(指定外键时不会在pg创建真实的外键约束)
        DisableForeignKeyConstraintWhenMigrating: true,
        // 指定表前缀
        NamingStrategy: schema.NamingStrategy{
            TablePrefix: global.Conf.Pg.TablePrefix + "_",
        },
        // 查询全部字段, 某些情况下*不走索引
        QueryFields: true,
    })
    if err != nil {
        panic(fmt.Sprintf("初始化pg异常: %v", err))
    }
    init = true
    // 开启pg日志
    if global.Conf.Pg.LogMode {
        db = db.Debug()
    }
    global.Pg = db
    // 表结构
    autoMigratePg()
    global.Log.Info("初始化pg完成")
    // 初始化数据库日志监听器
    binlog()
}

// 自动迁移表结构
func autoMigratePg() {
    global.Pg.AutoMigrate(
        new(models.SysUser),
        new(models.SysRole),
        new(models.SysMenu),
        new(models.SysApi),
        new(models.SysCasbin),
        new(models.SysWorkflow),
        new(models.SysWorkflowLine),
        new(models.SysWorkflowLog),
        new(models.RelationUserWorkflowLine),
        new(models.SysLeave),
        new(models.SysOperationLog),
        new(models.SysMessage),
        new(models.SysMessageLog),
        new(models.SysMachine),
    )
}
