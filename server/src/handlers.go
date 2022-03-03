package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

///////////////////////////////////////////////////////////////////////////////
// Middleware
///////////////////////////////////////////////////////////////////////////////

func authorize(fn func(c *fiber.Ctx) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sessionToken := c.Cookies("session_token")
		if sessionToken == "" {
			log.Error("'session_token' cookie not found")
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		res, err := cache.Do("GET", sessionToken)
		if err != nil {
			log.WithError(err).Error("Error checking redis")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if res == nil {
			log.WithError(err).Error("'session_token' cookie not found in redis")
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return fn(c)
	}
}

func logRequests(c *fiber.Ctx) error {
	if c.Method() != "OPTION" {
		user := c.Cookies("net_id")
		if user != "" {
			log.WithField("user", user).Info(
				c.Method(), " ", c.OriginalURL())
		} else {
			log.Info(c.Method(), " ", c.OriginalURL())
		}
	}

	return c.Next()
}

///////////////////////////////////////////////////////////////////////////////
// Handlers
///////////////////////////////////////////////////////////////////////////////

func registerHandler(c *fiber.Ctx) error {
	// Check if the netID is a student's
	// Create a 4 digit pin to validate with
	// Send an email to the netID with the pin
	data := new(map[string]interface{})
	if err := c.BodyParser(data); err != nil {
		log.WithError(err).Error("Error parsing body")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	username := fmt.Sprint((*data)["username"])
	password := fmt.Sprint((*data)["password"])
	if username == "<nil>" || password == "<nil>" {
		err := errors.New("Posted nil `username` or `password`")
		log.WithError(err).Error("Registration error")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if val, err := dbHasUsername(username); err != nil {
		log.WithError(err).Error("Error checking if db has username")
		return c.SendStatus(fiber.StatusInternalServerError)
	} else if val {
		log.Error("username already registered")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if err := checkNetID(username); err != nil {
		log.WithError(err).Error("Error validating netID")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	pin := fmt.Sprintf("%08d", randInt(0, 99999999))
	if err := addRegistrationPin(username, pin); err != nil {
		log.WithError(err).Error("Error adding registration pin to redis")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := cachePassword(username, password); err != nil {
		log.WithError(err).Error("Error caching login credentials")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := sendRegisterEmail(username, pin); err != nil {
		log.WithError(err).Error("Error sending registration email")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func confirmRegistrationHandler(c *fiber.Ctx) error {
	data := new(map[string]interface{})
	if err := c.BodyParser(data); err != nil {
		log.WithError(err).Error("Error parsing body")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	fmt.Println(data)

	username := fmt.Sprint((*data)["username"])
	pin := fmt.Sprint((*data)["pin"])
	if username == "<nil>" || pin == "<nil>" {
		err := errors.New("`username` or `pin` not sent with post")
		log.WithError(err).Error("Error parsing data")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if val, err := getRegistrationPin(username); err != nil {
		log.WithError(err).Error("Error checking redis for pin")
		return c.SendStatus(fiber.StatusInternalServerError)
	} else if val != pin {
		fmt.Println(val)
		fmt.Println(pin)
		err = errors.New("Pin in redis != posted value")
		log.WithError(err).Error("Confirm registration failed")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	password, err := getCachedPassword(username)
	if err != nil {
		log.WithError(err).Error("Error getting cached password from redis")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := addLogin(username, password); err != nil {
		log.WithError(err).Error("Error adding login to database")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := removeRegistrationPin(username); err != nil {
		msg := fmt.Sprintf("Error removing registration pin for %v", username)
		log.WithError(err).Error(msg)
	}
	log.Info(fmt.Sprintf("Login info added for %v", username))

	if err := removeCachedPassword(username); err != nil {
		msg := fmt.Sprintf("Error removing cached password for %v", password)
		log.WithError(err).Error(msg)
	}

	return c.SendStatus(fiber.StatusOK)
}

func resetPasswordHandler(c *fiber.Ctx) error {
	// need to generate a pin to associate with the reset
	// store the __reset_pin__NETID in redis with pin
	// need to send email with link to reset password

	return c.SendStatus(fiber.StatusOK)
}

func confirmResetPasswordHandler(c *fiber.Ctx) error {
	// need to check that __reset_pin__NETID pin matches incoming pin
	// replace old password with new password

	return c.SendStatus(fiber.StatusOK)
}

func loginHandler(c *fiber.Ctx) error {
	type Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	creds := new(Credentials)
	if err := c.BodyParser(creds); err != nil {
		log.Error(fmt.Sprintf("Error parsing login: %+v", err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	expectedPassword, err := getPassword(creds.Username)
	if err != nil {
		log.WithError(err).Error("Error fetching password from database")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if expectedPassword != creds.Password {
		log.Error(fmt.Sprintf("Passwords doesn't match for %v", creds.Username))
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	expireTime := 60 * 60 // 1 hour
	sessionToken, err := addSessionToken(creds.Username, expireTime)
	if err != nil {
		log.WithError(err).Error("Failed to add new session_token to redis")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(time.Hour),
	})

	c.Cookie(&fiber.Cookie{
		Name:    "net_id",
		Value:   creds.Username,
		Expires: time.Now().Add(time.Hour),
	})

	log.Info(fmt.Sprintf("User %v logged in", creds.Username))
	return c.SendStatus(fiber.StatusOK)
}

func logoutHandler(c *fiber.Ctx) error {
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		return c.SendStatus(fiber.StatusOK)
	}

	if _, err := cache.Do("DEL", sessionToken); err != nil {
		log.WithError(err).Error("Failed to delete token in Redis")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func addCourseHandler(c *fiber.Ctx) error {
	course := new(Course)
	if err := c.BodyParser(course); err != nil {
		log.WithError(err).Error("Failed to parse course")
		return err
	}

	if err := addCourse(*course); err != nil {
		log.WithError(err).Error("Failed to add course")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Info(fmt.Sprintf("Added %v", course.CourseNumber))
	return c.SendStatus(fiber.StatusOK)
}

func deleteCourseHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func getDataHandler(c *fiber.Ctx) error {
	term := c.Params("term")

	courses, err := getCourses(term)
	if err != nil {
		log.WithError(err).Error("Failed to get courses from database")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(courses)
}

///////////////////////////////////////////////////////////////////////////////
// Utils
///////////////////////////////////////////////////////////////////////////////

// Need to add sanitizer for strings
// func checkData(course Course) error {
// 	fields := []string{"term", "course_number", "professor", "time", "link", "user"}
// 	for _, field := range fields {
// 		if val, ok := course.field; !ok || val == nil {
// 			return errors.New(fmt.Sprintf("JSON request has no %v field", field))
// 		}
// 	}

// 	return nil
// }

func checkNetID(netID string) error {
	apiKey := os.Getenv("DUKE_API_KEY")
	if apiKey == "" {
		log.Error("DUKE_API_KEY not found")
		return errors.New("DUKE_API_KEY not found")
	}
	url := fmt.Sprintf(
		"https://streamer.oit.duke.edu/ldap/people/netid/%v?access_token=%v", netID, apiKey)

	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(url)

	if err := a.Parse(); err != nil {
		return err
	}

	_, body, _ := a.Bytes()

	type APIResponse struct {
		PrimaryAffiliation string `json:"primary_affiliation"`
	}

	data := make([]APIResponse, 0)
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if len(data) < 1 {
		return errors.New("Failed to fetch data")
	}

	role := fmt.Sprint(data[0].PrimaryAffiliation)
	if role != "Student" {
		return errors.New(fmt.Sprintf("Primary affiliation for %v is %v, not Student", netID, role))
	}

	return nil
}

func sendRegisterEmail(username string, pin string) error {
	from := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")

	if from == "" || password == "" {
		return errors.New("EMAIL_USERNAME or EMAIL_PASSWORD env variable not set")
	}

	to := []string{
		fmt.Sprintf("%v@duke.edu", username),
	}

	subject := "Subject: Register for GroupDuke\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	link := fmt.Sprintf("%v/confirm/%v/%v", origin, username, pin)
	body := fmt.Sprintf("To confirm your registration, click this link: <a href=\"%v\">%v</a>", link, link)
	message := []byte(subject + mime + body)

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Sent email to %v@duke.edu", username))
	return nil
}

func randInt(low, high int) int {
	return low + rand.Intn(high-low)
}
