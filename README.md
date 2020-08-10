# Etpmls-Admin-Server

## Configuration 配置

1.Copy .env.example to .env

1.将.env.example复制到.env

2.Copy storage/config/app.yaml.example to storage/config/app.yaml

2.将storage / config / app.yaml.example复制到storage / config / app.yaml

3.Copy storage/config/app_debug.yaml.example to storage/config/app_debug.yaml

3.将storage / config / app_debug.yaml.example复制到storage / config / app_debug.yaml

And configure them

并且配置它们

4
```shell script
go mod vendor
```

## Run 运行

PostgreSQL
```shell script
go run -tags=postgresql main.go
```

MySQL/MariaDB
```shell script
go run -tags=mysql main.go
```

