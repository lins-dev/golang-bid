package main

import "fmt"

func main()  {
	fmt.Println("Hello World!!!!!")
	fmt.Println("API")
}

// run Air
// air --build.cmd "go build -o ./bin/api ./cmd/api" --build.bin "./bin/api"

// create migrations
// in folder: internal/store/pgstore/migrations
// tern new create_users_table

// generate SQLc files
// sqlc generate -f ./internal/store/pgstore/sqlc.yml