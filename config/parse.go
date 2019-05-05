package config

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"os"
	"os/user"
)

func ParseConfig() (Config, error){
	conf := Config{}
	errConfig := godotenv.Load(".env")
	u,_ := user.Current()
	if errConfig != nil{
		if _, err := os.Stat(u.HomeDir + "/.opf/.env"); os.IsNotExist(err){
			return Config{}, errors.New("Please create .env file on root path.")
		}
		err := godotenv.Load(u.HomeDir + "/.opf/.env")
		if err != nil{
			return Config{}, errors.New("Please create .env file on root path.")
		}
	}


	conf.Instagram.Login = os.Getenv("INSTAGRAM_LOGIN")
	conf.Instagram.Password = os.Getenv("INSTAGRAM_PASSWORD")

	conf.Twitter.Api.Key = os.Getenv("TWITTER_CONSUMER_SECRET")
	conf.Twitter.Login = os.Getenv("TWITTER_CONSUMER")
	conf.Twitter.Password = os.Getenv("TWITTER_ACCESS_TOKEN")
	conf.Twitter.Api.SKey = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	conf.Api.Host = os.Getenv("API_HOST")
	conf.Api.Port = os.Getenv("API_PORT")
	conf.Api.Verbose = os.Getenv("API_VERBOSE")

	conf.Common.HistoryFile = os.Getenv("OPERATIVE_HISTORY")

	conf.Database.Driver = os.Getenv("DB_DRIVER")
	conf.Database.Name = os.Getenv("DB_NAME")
	conf.Database.Host = os.Getenv("DB_HOST")
	conf.Database.User = os.Getenv("DB_USER")
	conf.Database.Pass = os.Getenv("DB_PASS")

	return conf, nil
}
