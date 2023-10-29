package env

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"

	"github.com/labstack/gommon/log"
)

//go:embed dotenv/*
var dotenvFS embed.FS

type (
	Env struct {
		App RestEnv
		Log LogEnv
	}

	RestEnv struct {
		Timeout      time.Duration `env:"TIMEOUT,default=2s"`
		Production   bool          `env:"PRODUCTION,default=false"`
		AppName      string        `env:"APP_NAME,required"`
		ReadTimeout  time.Duration `env:"READ_TIMEOUT,default=10s"`
		WriteTimeout time.Duration `env:"WRITE_TIMEOUT,default=10s"`
		ListenPort   string        `env:"LISTEN_PORT,default=:8080"`
		Http2        bool          `env:"DISABLED_HTTP2,default=false"`
	}

	LogEnv struct {
		LogLocation      string `env:"LOG_LOCATION,default=STD"`
		LogLevel         string `env:"LOG_LEVEL, default=DEBUG"`
		TimeFieldName    string `env:"TIME_FIELD_NAME,default=time"`
		MessageFieldName string `env:"MSG_FIELD_NAME,default=message"`
		ErrorFieldName   string `env:"ERR_FIELD_NAME,default=error"`
		TimeFieldFormat  string `env:"TIME_FIELD_FORMAT,default=UNIX"`
		// Throttle
		DebugBurstLvl    uint32        `env:"DEBUG_BURST_LEVEL,default=5"`
		DebugBurstPeriod time.Duration `env:"DEBUG_BURST_PERIOD,default=1s"`
		DebugN           uint32        `env:"N_DEBUG,default=100"`
		InfoBurstLvl     uint32        `env:"INFO_BURST_LEVEL,default=50"`
		InfoBurstPertiod time.Duration `env:"INFO_BURST_PERIOD,default=1s"`
		InfoN            uint32        `env:"N_INFO,default=100"`

		WarnBurstLvl     uint32        `env:"WARN_BURST_LEVEL,default=50"`
		WarnBurstPertiod time.Duration `env:"WARN_BURST_PERIOD,default=1s"`
		WarnN            uint32        `env:"N_WARN,default=100"`
	}
)

func NewEnv(ctx context.Context) *Env {
	env := loadEnvFile(ctx)
	if !env.App.Production {
		fmt.Println("Running App in Development Env 🔥")
	}

	return env

}

// Load and process *.env into struct
func loadEnvFile(ctx context.Context) *Env {

	env_files := readEnvFiles("dotenv")

	err := godotenv.Load(env_files...)
	if err != nil {
		log.Fatalf("[Error] - Load env file, %v", err.Error())
	}

	env := &Env{}

	err = envconfig.Process(ctx, env)
	if err != nil {
		log.Fatalf("[Error] - serialize env file, %v", err.Error())
	}

	return env
}

// read nested embedded dotenv file  , create tempFile for goDotEnv to load
func readEnvFiles(path string) []string {
	env_files := []string{}
	var walkDir func(string)

	tmp := os.TempDir()
	walkDir = func(path string) {
		files, err := dotenvFS.ReadDir(path)
		if err != nil {
			log.Fatalf("[Error] - Read nested dotenv Dir : %v", err)
		}

		for _, file := range files {
			filePath := filepath.Join(path, file.Name())
			if file.IsDir() {
				walkDir(filePath)
			} else {
				fileName := filepath.Join(tmp, file.Name())
				data, err := dotenvFS.ReadFile(filePath)
				if err != nil {
					log.Fatalf("[Error] - Read embedded Env file , %v", err.Error())
				}
				err = os.WriteFile(fileName, data, 0600)

				if err != nil {
					log.Fatalf("[Error] - Write Temp Env file , %v", err.Error())
				}

				env_files = append(env_files, fileName)

			}

		}
	}

	walkDir(path)
	return env_files

}
