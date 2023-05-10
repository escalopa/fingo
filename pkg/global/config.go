package global

import "github.com/spf13/viper"

func LoadConfig(x interface{}, name string, path string, form string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(form)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(x)
	if err != nil {
		return err
	}

	return nil
}
