#  API Reporting

## ⚙️ Specifications
    Written in Go version : 1.19

## 📚 Repo Structure
```
├── db
├── docs
├── handler
├── libs
│   ├── logger
│   └── util
├── logs
├── migrations
├── seeder
├── repository
├── schema
├── server
│   └── middleware
└── service
```

- `db` contains initiator to open connection database
- `docs` contains documentation of project
- `handler` contains go package layer to handle requests from http (request layer)
- `libs` contains shared code that can be used on each packages
- `logs` contains logging file
- `migrations` contains migrations file
- `seeder` contains example of table records in database
- `repository` contains go package layer to serve a requests from service (source data layer)
- `schema` contains shared code that can be used on other packages in context entity structure
- `server` contains a go http server and middleware
- `service` contains go package layer to serve a requests from handler (business logic layer)

## 🔧 Running Locally
To run this project you need some preparation :
- `create database 'reporting'` 
- `installing migrator tools` download from [golang migrate](https://github.com/golang-migrate/migrate/releases) in release page
- `migrate -path ./migrations -database "mysql://root:secret@tcp(localhost:3306)/reporting" -verbose up` run this command to up a migration
- `go mod tidy` installing a module
- `go run .` run it

## 🔧 Migration Guide 
`migrate create -ext sql -dir ./migrations -format unix -tz UTC name_of_migration` - Create a migration
`migrate -path ./migrations -database "mysql://root:secret@tcp(localhost:3306)/reporting" -verbose up`  - migrate up migration
`migrate -path ./migrations -database "mysql://root:secret@tcp(localhost:3306)/reporting" -verbose down` - migrate down migration


## 📦 Go Library

Using [Go Chi](https://github.com/go-chi/chi) as router for building HTTP services, looking a [Docs](https://github.com/go-chi/chi).


## 📰 Go Article

[Download Golang Binnary](https://go.dev/dl/)

[How to install Go in PC / Laptop / Server](https://go.dev/doc/install)

## 📚 Go Book

[Go Tutorial - Bahasa](https://dasarpemrogramangolang.novalagung.com/)

## 💡 Go Command

[CMD List Golang](https://go.dev/cmd/go/)

## 🧷 Recommended IDE

[Visual Studio Code](https://code.visualstudio.com/)

## 🔧 Recommended Extension Visual Studio Code

[GO Extension on VSCode](https://marketplace.visualstudio.com/items?itemName=golang.go)