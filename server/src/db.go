package main

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/gomodule/redigo/redis"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/nu7hatch/gouuid"
	"google.golang.org/api/option"
)

type Course struct {
	ID           int    `json:"id"`
	Term         string `json:"term"`
	CourseNumber string `json:"course_number"`
	Professor    string `json:"professor"`
	Time         string `json:"time"`
	Link         string `json:"link"`
	User         string `json:"user"`
}

func initCache(redisURL string) error {
	conn, err := redis.DialURL(redisURL)
	if err != nil {
		return err
	}

	cache = conn
	return nil
}

func initFirebase() error {
	opt := option.WithCredentialsFile(dbCreds)
	config := &firebase.Config{
		DatabaseURL: dbURL,
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return err
	}

	c, err := app.Database(context.Background())
	if err != nil {
		return err
	}

	client = c
	return nil
}

func addSessionToken(value string, expireTime int) (string, error) {
	uu, err := uuid.NewV4()
	sessionToken := uu.String()

	_, err = cache.Do("SETEX", sessionToken, fmt.Sprint(expireTime), value)
	return sessionToken, err
}

func addRegistrationPin(username string, pin string) error {
	_, err := cache.Do("SET", fmt.Sprintf("_pin_%v", username), pin)
	return err
}

func getRegistrationPin(username string) (string, error) {
	val, err := redis.String(cache.Do("GET", fmt.Sprintf("_pin_%v", username)))
	return val, err
}

func removeRegistrationPin(username string) error {
	_, err := cache.Do("DEL", fmt.Sprintf("_pin_%v", username))
	return err
}

func cachePassword(username string, password string) error {
	_, err := cache.Do("SET", fmt.Sprintf("_creds_%v", username), password)
	return err
}

func getCachedPassword(username string) (string, error) {
	val, err := redis.String(cache.Do("GET", fmt.Sprintf("_creds_%v", username)))
	return val, err
}

func removeCachedPassword(username string) error {
	_, err := cache.Do("DEL", fmt.Sprintf("_creds_%v", username))
	return err
}

func addCourse(course Course) error {
	term := course.Term

	_, err := client.NewRef("courses").
		Child(term).
		Push(context.Background(), course)

	return err
}

func getCourses(term string) ([]Course, error) {
	var results []Course

	terms := []string{term}
	if term == "" {
		terms = []string{"Sp22", "Fa21", "Su21"}
	}

	for _, term := range terms {
		var result map[string]interface{}
		err := client.NewRef("courses").Child(term).Get(context.Background(), &result)
		if err != nil {
			return nil, err
		}

		for _, obj := range result {
			m := obj.(map[string]interface{})
			newCourse := Course{
				CourseNumber: fmt.Sprint(m["course_number"]),
				Link:         fmt.Sprint(m["link"]),
				Professor:    fmt.Sprint(m["professor"]),
				Term:         fmt.Sprint(m["term"]),
				Time:         fmt.Sprint(m["time"]),
				User:         fmt.Sprint(m["user"]),
			}

			results = append(results, newCourse)
		}
	}

	return results, nil
}

// REPLACE
func deleteCourse(idToDelete int) (int64, error) {
	return 0, nil
}

func dbHasUsername(username string) (bool, error) {
	var result string
	err := client.NewRef("logins").Child(username).Get(context.Background(), &result)

	return result != "", err
}

func addLogin(username string, password string) error {
	err := client.NewRef("logins").
		Child(username).
		Set(context.Background(), password)

	return err
}

func getPassword(username string) (string, error) {
	var result string
	err := client.NewRef("logins").Child(username).Get(context.Background(), &result)

	if err != nil {
		return "", err
	}
	return result, nil
}
