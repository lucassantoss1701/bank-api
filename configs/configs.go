package configs

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

var configuration *Config

type Config struct {
	AppName  string `mapstructure:"APP_NAME" default:"bank-api"`
	Server   server
	Database database
	Security security
}

type database struct {
	Type string `mapstructure:"DB_TYPE" default:"mysql"`
	User string `mapstructure:"DB_USER" default:"test"`
	Pass string `mapstructure:"DB_PASS" default:"test"`
	Host string `mapstructure:"DB_HOST" default:"bank"`
	Port string `mapstructure:"DB_PORT" default:"3306"`
	Name string `mapstructure:"DB_NAME" default:"bank"`
}

type server struct {
	Host string `mapstructure:"SERVER_HOST" default:"8000"`
}

type security struct {
	Secret string `mapstructure:"SECRET" default:""`
}

func getMappedEnvs(configStruct reflect.Type) []string {
	result := make([]string, 0)

	for i := 0; i < configStruct.NumField(); i++ {
		field := configStruct.Field(i)
		if configName := field.Tag.Get("mapstructure"); configName != "" {
			result = append(result, configName)
		}
		if field.Type.Kind() == reflect.Struct {
			result = append(result, getMappedEnvs(field.Type)...)
		}
	}
	return result
}

func setDefaultValues(configStruct reflect.Type) {
	for i := 0; i < configStruct.NumField(); i++ {
		field := configStruct.Field(i)
		configName := field.Tag.Get("mapstructure")
		defaultValue := field.Tag.Get("default")

		if configName != "" && defaultValue != "" {
			viper.SetDefault(configName, defaultValue)
		}

		if field.Type.Kind() == reflect.Struct {
			setDefaultValues(field.Type)
		}
	}
}

func Load() error {
	configuration = &Config{}

	envFile := ".env-development"
	envPath := "."

	viper.AddConfigPath(envPath)
	viper.SetConfigName(envFile)
	viper.SetConfigType("env")

	setDefaultValues(reflect.TypeOf(Config{}))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("[Method: Config.Load()]", envFile, "not found, load by environment variables")

			viper.AutomaticEnv()
			mapped := getMappedEnvs(reflect.TypeOf(Config{}))
			for _, env := range mapped {
				viper.BindEnv(env)
			}

		} else {
			return err
		}
	} else {
		fmt.Println("[Method: Config.Load()] Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		return err
	}

	if err := viper.Unmarshal(&configuration.Server); err != nil {
		return err
	}

	if err := viper.Unmarshal(&configuration.Database); err != nil {
		return err
	}

	if err := viper.Unmarshal(&configuration.Security); err != nil {
		return err
	}

	return nil

}

// Get returns a Config Structure
func Get() *Config {
	return configuration
}
