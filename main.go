package main

import (
	"project-app-bioskop/cmd"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/wire"
	"project-app-bioskop/pkg/database"
	"project-app-bioskop/pkg/utils"
)

func main() {
	config, err := utils.ReadConfiguration()
	if err != nil {
		// print error
	}

	db, err := database.InitDB(config.DB)
	if err != nil {
		// print err
	}

	repo := repository.NewRepository(db)
	route := wire.Wiring(*repo)
	cmd.APiserver(route)
}
