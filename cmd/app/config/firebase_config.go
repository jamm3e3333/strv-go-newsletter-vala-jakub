package config

import "github.com/ilyakaznacheev/cleanenv"

type FirebaseConfig struct {
	SAKeyFileEnc  string `env:"CONFIG_FIREBASE_SERVICE_ACCOUNT_KEY_FILE_ENCODED"`
	DBUrl         string `env:"CONFIG_FIREBASE_DATABASE_URL"`
	IsTestingMode bool   `env:"CONFIG_FIREBASE_IS_TESTING_MODE"`
}

func CreateFirebaseConfig() (FirebaseConfig, error) {
	var cfg FirebaseConfig
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
