package global

import "strings"

type env struct {
}

func Env() *env {
	return &env{}
}

func (e *env) IsProd() bool {
	env := strings.ToLower(Config.App.Env)
	return env == "prod" || env == "production"
}

func (e *env) IsLocal() bool {
	env := strings.ToLower(Config.App.Env)
	return env == "local" || env == "localhost"
}

func (e *env) IsTest() bool {
	env := strings.ToLower(Config.App.Env)
	return env == "test" || env == "qa"
}
