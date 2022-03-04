package main

import (
	"os"

	"firebase.google.com/go/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

var cache redis.Conn
var client *db.Client
var dbCreds = os.Getenv("FIREBASE_CREDS_FILE")
var dbURL = os.Getenv("FIREBASE_DB_URL")
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

	app.Post("/reset-password", resetPasswordHandler)
	app.Post("/confirm-reset-password", confirmResetPasswordHandler)

	app.Post("/login", loginHandler)
	app.Post("/logout", logoutHandler)

	app.Post("/contact", contactHandler)

	app.Get("/data", authorize(getDataHandler))
	app.Post("/add-course", authorize(addCourseHandler))
	app.Delete("/delete-course", authorize(deleteCourseHandler))

	log.Info("starting redis")
	if err := initCache(redisURL); err != nil {
		log.Fatal(err)
	}

	log.Info("connecting to firebase")
	if err := initFirebase(); err != nil {
		log.WithError(err).Fatal("Error initializing firebase")
	}

	if local[0] != ':' {
		local = ":" + local
	}
	log.Fatal(app.Listen(local))
}
