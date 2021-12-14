package config

import (
	"errors"
	"os"
)

//GetEnvironmentVariables changes bases on the Secret Manager
func GetEnvironmentVariables() (port string, host string, err error) {
	port = os.Getenv("PORT")
	if port == "" {
		return port, host, errors.New("missing Port Environment Variable")
	}

	host = os.Getenv("HOST")
	if host == "" {
		return port, host, errors.New("missing ProjectId Environment Variable")
	}

	return port, host, nil
}
