package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/nu7hatch/gouuid"
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

func addSessionToken(value string, expireTime int) (string, error) {
	uu, err := uuid.NewV4()
	sessionToken := uu.String()

	_, err = cache.Do("SETEX", sessionToken, fmt.Sprint(expireTime), value)
	return sessionToken, err
}

func addRegistrationPin(netID string, pin string) error {
	_, err := cache.Do("SET", fmt.Sprintf("_pin_%v", netID), pin)
	return err
}

func getRegistrationPin(netID string) (string, error) {
	val, err := redis.String(cache.Do("GET", fmt.Sprintf("_pin_%v", netID)))
	return val, err
}

func removeRegistrationPin(netID string) error {
	_, err := cache.Do("DEL", fmt.Sprintf("_pin_%v", netID))
	return err
}

func cachePassword(netID string, password string) error {
	_, err := cache.Do("SET", fmt.Sprintf("_creds_%v", netID), password)
	return err
}

func getCachedPassword(netID string) (string, error) {
	val, err := redis.String(cache.Do("GET", fmt.Sprintf("_creds_%v", netID)))
	return val, err
}

func removeCachedPassword(netID string) error {
	_, err := cache.Do("DEL", fmt.Sprintf("_creds_%v", netID))
	return err
}

func dbHasNetID(netID string) (bool, error) {
	db, err := sql.Open("sqlite3", dbString)
	if err != nil {
		defer db.Close()
		return false, err
	}

	query := fmt.Sprintf("SELECT * FROM logins WHERE logins.username = '%v'", netID)
	rows, err := db.Query(query)
	defer db.Close()
	defer rows.Close()
	if err != nil {
		return false, err
	}

	if err := rows.Err(); err != nil {
		return false, err
	}

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

func addCourse(newCourse Course) error {
	db, err := sql.Open("sqlite3", dbString)
	if err != nil {
		defer db.Close()
		return err
	}

	stmt, err := db.Prepare("INSERT INTO courses (id, term, course_number, professor, time, link, user) VALUES (?, ?, ?, ?, ?, ?, ?)")
	defer db.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		nil,
		newCourse.Term,
		newCourse.CourseNumber,
		newCourse.Professor,
		newCourse.Time,
		newCourse.Link,
		newCourse.User,
	)
	if err != nil {
		return err
	}

	defer stmt.Close()
	return nil
}

func getCourses(term string) ([]Course, error) {
	db, err := sql.Open("sqlite3", dbString)
	if err != nil {
		defer db.Close()
		return nil, err
	}

	query := fmt.Sprintf("SELECT * FROM courses WHERE courses.term LIKE '%v'", term)
	rows, err := db.Query(query)
	defer db.Close()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	courses := make([]Course, 0)
	for rows.Next() {
		newCourse := Course{}
		err = rows.Scan(
			&newCourse.ID,
			&newCourse.Term,
			&newCourse.CourseNumber,
			&newCourse.Professor,
			&newCourse.Time,
			&newCourse.Link,
			&newCourse.User)
		if err != nil {
			return nil, err
		}

		courses = append(courses, newCourse)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func deleteCourse(idToDelete int) (int64, error) {
	db, err := sql.Open("sqlite3", dbString)
	if err != nil {
		defer db.Close()
		return 0, err
	}

	stmt, err := db.Prepare("DELETE FROM courses WHERE id = ?")
	if err != nil {
		defer stmt.Close()
		return 0, err
	}

	defer db.Close()
	defer stmt.Close()

	res, err := stmt.Exec(idToDelete)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	fmt.Printf("Deleted id=%v", idToDelete)
	return affected, nil
}

// Ensure that the netID doesn't already exist
func addLogin(username string, password string) error {
	db, err := sql.Open("sqlite3", dbString)
	if err != nil {
		defer db.Close()
		return err
	}

	query := fmt.Sprintf("SELECT * FROM logins WHERE logins.username LIKE '%v'", username)
	rows, err := db.Query(query)
	defer db.Close()
	defer rows.Close()

	if rows.Next() {
		return errors.New("Username already registered")
	}

	stmt, err := db.Prepare("INSERT INTO logins (id, username, password) VALUES (?, ?, ?)")
	defer db.Close()
	if err != nil {
		defer stmt.Close()
		return err
	}

	_, err = stmt.Exec(
		nil,
		username,
		password,
	)
	if err != nil {
		defer stmt.Close()
		return err
	}

	defer stmt.Close()
	return nil
}

func getPassword(username string) (string, error) {
	db, err := sql.Open("sqlite3", dbString)
	if err != nil {
		defer db.Close()
		return "", err
	}

	rows, err := db.Query(fmt.Sprintf("SELECT password FROM logins WHERE logins.username='%v'", username))
	if err != nil {
		defer db.Close()
		return "", err
	}
	defer db.Close()

	var password string
	if rows.Next() {
		err = rows.Scan(&password)
	}
	defer rows.Close()

	return password, nil
}
