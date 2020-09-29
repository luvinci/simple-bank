### 如何在Golang中编写和运行数据库迁移？

安装好MySQL或其他数据库，创建一个数据库`simple_bank`（后面需要用到）:

1. 安装`migrate`（https://github.com/golang-migrate/migrate/releases）

2. 在项目根路径下创建`db/migration`文件夹，用来存储所有的迁移文件

3. 创建migrations

    ```shell
    migrate create -ext sql -dir db/migration -seq init_schema
    # -ext 文件扩展名为sql
    # -seq 标志生成的迁移文件的顺序版本号(init_schema任意写，比如：create_accounts)
    # init_schema为迁移的名称
    ```

    可以看到`db/migration`下生成了两个文件：

    ```
    000001_init_schema.up.sql
    000001_init_schema.down.sql
    ```

    打开要转移的包含sql语句的文件`*.sql`，将所有sql语句复制到`000001_init_schema.up.sql`里面，比如我这里的sql语句就是创建了三个表（为了方便展示删除了多余内容）：

    ```
    CREATE TABLE accounts (
    	...
    ) ENGINE = InnoDB;
    CREATE TABLE entries (
    	...
    )ENGINE = InnoDB;
    CREATE TABLE transfers (
    	...
    )ENGINE = InnoDB;
    ```

    对于`000001_init_schema.down.sql`文件，我们应还原其对应`000001_init_schema.up.sql`中sql语句所做的更改，因为up中是创建三个表，所以对应的down就是删除这三个表（如果有外键关联，注意执行顺序）：

    ```
    DROP TABLE IF EXISTS accounts;
    DROP TABLE IF EXISTS entries;
    DROP TABLE IF EXISTS transfers;
    ```

4. 运行migrations

    ```shell
    migrate -path db/migration -database "postgresql://postgres:123456@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up
    # -path 指定迁移文件的文件夹
    # -database 指定DSN
    # -verbose 用来在迁移时打印详细日志记录
    # 默认情况下，postgres容器未启用SSL，所以加上sslmode=disable
    ```

### Makefile

```makefile
migrateup:
	migrate -path db/migration -database "postgresql://postgres:123456@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:123456@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down

.PHONY: migrateup migratedown
```

`make migrateup`向上迁移，`make migratedown`向下迁移。



