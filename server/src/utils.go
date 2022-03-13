package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

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

func sendEmail(to []string, subject string, body string) error {
	from := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")

	if from == "" || password == "" {
		return errors.New("EMAIL_USERNAME or EMAIL_PASSWORD env variable not set")
	}

	subject = fmt.Sprintf("Subject: %v\n", subject)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(subject + mime + body)

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}

func randInt(low, high int) int {
	return low + rand.Intn(high-low)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
