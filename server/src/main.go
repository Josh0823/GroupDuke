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
var local = os.Getenv("PORT")
var origin = os.Getenv("ORIGIN_URL")
var redisURL = os.Getenv("REDIS_URL")

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     "Content-Type",
		AllowOrigins:     origin,
	}))

	app.Use(logRequests)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/register", registerHandler)
	app.Post("/confirm-registration", confirmRegistrationHandler)

	app.Post("/login", loginHandler)
	app.Post("/logout", logoutHandler)

	app.Get("/data", authorize(getDataHandler))
	app.Post("/add-course", authorize(addCourseHandler))
	app.Delete("/delete-course", authorize(deleteCourseHandler))

	log.Info("starting redis")
	if err := initCache(redisURL); err != nil {
		log.Fatal(err)
	}

	if dbString == "" {
		log.Fatal("DB_STRING env variable not set")
	}
	log.Info("DB_STRING: ", dbString)

	if local[0] != ':' {
		local = ":" + local
	}
	log.Fatal(app.Listen(local))
}
