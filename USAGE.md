### 第一步

将gcrawl.sql导入mysql数据库

### 第二步

修改main.go 里面的参数

```go
var (
	mysql_user = "root"
	mysql_pwd  = "root"
	mysql_db   = "gcrawl"
	conditions = map[string]interface{}{ //  https://github.com/louzhu123/gcrawl
		"position": "后端开发",
	}
)
```

### 第三步

```shell
go run main.go
```

