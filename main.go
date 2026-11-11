package main

import (
	"log"

	"firlydev/backend-test/app/configs"
	"firlydev/backend-test/app/database"
	"firlydev/backend-test/app/middlewares"
	"firlydev/backend-test/app/routers"
	"firlydev/backend-test/features/applicant/data"
	"firlydev/backend-test/features/applicant/service"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg := configs.Load()
	db, err := database.Connect(cfg.DBPath)
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}
	if cfg.SeedData {
		if err := database.SeedDemoData(db); err != nil {
			log.Fatal("database seeding failed: ", err)
		}
	}
	repo := data.NewRepository(db)
	appService := service.New(repo)
	e := echo.New()
	middlewares.Register(e)
	routers.Register(e, appService)
	log.Printf("server listening on :%s", cfg.Port)
	log.Fatal(e.Start(":" + cfg.Port))
}
