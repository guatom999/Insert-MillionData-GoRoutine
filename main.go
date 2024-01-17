package main

import (
	"context"
	"onemildata/config"
	"onemildata/database"
	"onemildata/server"
)

func main() {

	ctx := context.Background()

	cfg := config.GetConfig()

	db := database.NewPostgresDb(&cfg).GetDb()

	server.NewEchoServer(db).Start(ctx)

}
