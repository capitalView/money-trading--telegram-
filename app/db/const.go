package db

import "main/utils"

var (
	User     = utils.GoDotEnvVariable("DB_USER")
	Password = utils.GoDotEnvVariable("DB_PASSWORD")
	Name     = utils.GoDotEnvVariable("DB_NAME")
	Host     = utils.GoDotEnvVariable("DB_HOST")
)
