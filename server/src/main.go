package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

var cache redis.Conn
var dbString = os.Getenv("DB_STRING")
var local = ":3000"
var origin = os.Getenv("ORIGIN_URL")

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     "Content-Type",
		AllowOrigins:     origin,
	}))

	app.Use(logRequests)

	app.Post("/register", registerHandler)
	app.Post("/login", loginHandler)
	app.Post("/logout", logoutHandler)

	// app.Get("/data", authorize(getDataHandler))
	app.Get("/data", getDataHandler)
	app.Post("/add-course", authorize(addCourseHandler))
	app.Delete("/delete-course", authorize(deleteCourseHandler))

	log.Info("starting redis")
	if err := initCache(os.Getenv("REDIS_URL")); err != nil {
		log.Fatal(err)
	}

	log.Fatal(app.Listen(local))
}
