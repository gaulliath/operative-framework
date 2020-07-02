package engine

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/user"
)

type Files struct {
	Config string `json:"config"`
	Job    string `json:"job"`
	Base   string `json:"base"`
	Export string `json:"export"`
}

func Preload() (Files, error) {
	// Load Configuration File
	configFile := ".env"
	err := godotenv.Load(".env")

	if err != nil {

		// Generate Default .env File
		u, errU := user.Current()
		if errU != nil {
			fmt.Println("Please create a '.env' file on root path.")
			return Files{}, err
		}
		if _, err := os.Stat(u.HomeDir + "/.opf/.env"); os.IsNotExist(err) {
			if _, err := os.Stat(u.HomeDir + "/.opf/"); os.IsNotExist(err) {
				_ = os.Mkdir(u.HomeDir+"/.opf/", os.ModePerm)
			}

			if _, err := os.Stat(u.HomeDir + "/.opf/webhooks"); os.IsNotExist(err) {
				_ = os.Mkdir(u.HomeDir+"/.opf/webhooks", os.ModePerm)
			}

			log.Println("Generating default .env on '" + u.HomeDir + "' directory...")
			path, errGeneration := GenerateEnv(u.HomeDir + "/.opf/.env")
			if errGeneration != nil {
				return Files{}, err
			}
			err := godotenv.Load(path)
			if err != nil {
				log.Println(err.Error())
				return Files{}, err
			}
		}
		configFile = u.HomeDir + "/.opf/.env"
		configJob := u.HomeDir + "/.opf/cron/"
		opfBaseDirectory := u.HomeDir + "/.opf/"
		opfExport := opfBaseDirectory + "export/"

		return Files{
			Config: configFile,
			Job:    configJob,
			Base:   opfBaseDirectory,
			Export: opfExport,
		}, nil
	}

	return Files{}, nil
}
