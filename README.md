# Best Go Micro Service template

## Architecture
### 1. Gin Restful API
### 2. Grpc server
### 3. Database mysql(migrate, downgrate)
### 4. Loging logrus
### 5. Configuration
### ....

## Install Go
Refer to: https://golang.org/doc/install

## Install MariaDB
Create/configure an admin user in the mariadb database

Reference:

[Linux]
https://www.digitalocean.com/community/tutorials/how-to-install-mariadb-on-ubuntu-18-04

[Windows]
https://mariadb.com/kb/en/installing-mariadb-msi-packages-on-windows/

## Configure
Change `[mariadb]` to your LOCAL DATABASE CONFIG for development in `etc/order_config.toml`:
```
# Example:

[mariadb]
address  = "127.0.0.1"
username = "myusername"
password = "123456"
port     = 3306
database = "casorder"
```
Then, go to file `utils/config/config.go` and change config path:
```
[Linux]
viper.AddConfigPath("Absolute/Path/To/cas-order/etc") (Example: "/home/mypc/some_dir/cas-order/etc")

[Windows]
viper.AddConfigPath("Absolute\\Path\\To\\cas-order\\etc") (Example: "E:\\SomeWorkSpace\\cas-order\\etc")
```


## Update Modules
```sh
go mod tidy
```

## Build executables

### Linux
```
Build:
export GOBIN=./bin/
go build -o ${GOBIN} ./cmd/api/
go build -o ${GOBIN} ./cmd/manage/
go build -o ${GOBIN} ./cmd/taskmanager/

Use:
./bin/api
./bin/manage --help
./bin/taskmanager
```

### Windows
```
Build:
go build -o .\bin\ .\cmd\api\
go build -o .\bin\ .\cmd\manage\
go build -o .\bin\ .\cmd\taskmanager\

Use:
.\bin\api.exe
.\bin\manage.exe --help
.\bin\taskmanager.exe
```

## Migrate DB
```md
[LINUX]
./bin/manage db-migrate

[WINDOWS]
.\bin\manage.exe db-migrate
```

## Install Protoc
```sh
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
go get -u -d github.com/golang-migrate/migrate/cmd/migrate
```

## Build Protoc 

```sh
cd taskmanager/grpc
./clear.sh
./build.sh
```

## Test HTTP Request
```sh
curl -X GET http://127.0.0.1:8090/casorder/api/v1/health/check
```

## Test
```sh
./bin/manage test-db
./bin/manage test-grpc
```
