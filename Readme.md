# Kazokku Users API

## clone this repository
```bash
git clone https://github.com/your-username/kazokku-users.git
```

## install dependencies
```bash
go mod tidy
```


## Configure the following environment variables for the server:

```bash
//server setting
SERVER_URI=localhost
SERVER_PORT=3000

//database setting
DB_HOST=localhost
DB_PORT=5432
DB_NAME=kazokku_users
DB_USERNAME=postgres
DB_PASSWORD=dimasslalu123
DB_POOL_MIN=10
DB_POOL_MAX=100
DB_TIMEOUT=10
DB_MAX_IDLE_TIME_SECOND=60


//Redis setting
REDIS_PASSWORD:dimasslalu123
```


## Run the applications

```bash
go run main.go
```

## Postman Documentation

```bash
https://documenter.getpostman.com/view/23663611/2sAYBSktdW
```