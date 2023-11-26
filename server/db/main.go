package db

import "os"

const DB_HOST_KEY = "DB_HOST"
const DB_PORT_KEY = "DB_PORT"
const DB_PASSWORD_KEY = "DB_PASSWORD"
const DB_USER_KEY = "DB_USER"
const DB_NAME_KEY = "DB_NAME"

func GetDSN() string {
	dbHost := os.Getenv(DB_HOST_KEY)
	dbPort := os.Getenv(DB_PORT_KEY)
	dbUser := os.Getenv(DB_USER_KEY)
	dbPassword := os.Getenv(DB_PASSWORD_KEY)
	dbName := os.Getenv(DB_NAME_KEY)

	return dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
}
