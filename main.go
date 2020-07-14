package main

import "os"

func main() {
	a := App{}
	a.Init(getEnv("APP_DB_HOST", "localhost"),
		getEnv("APP_DB_PORT", "5432"),
		getEnv("APP_DB_USERNAME", "testdb"),
		getEnv("APP_DB_PASSWORD", "testdb"),
		getEnv("APP_DB_NAME", "testdb"))

	a.Run(":" + getEnv("APP_SERVICE_PORT", "8080"))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
