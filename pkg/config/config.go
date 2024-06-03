package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	GlobalToken string
	LimitBody   int64
	HttpLogs    bool
}

func LoadEnv() (*Env, error) {
	if isDocker := os.Getenv("DOCKER_ENV") == "true"; !isDocker {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	global_token := os.Getenv("GLOBAL_AUTH_TOKEN")
	if global_token == "" {
		return nil, errors.New("global token not defined")
	}

	limit_body, err := strconv.ParseInt(os.Getenv("BODY_SIZE"), 10, 64)
	if err != nil {
		return nil, err
	}

	if limit_body == 0 {
		limit_body = 5
	}

	env := Env{
		GlobalToken: global_token,
		LimitBody:   limit_body,
		HttpLogs:    os.Getenv("HTTP_LOGS") == "true",
	}

	return &env, nil
}
