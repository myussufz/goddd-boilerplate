package config

import "os"

// Configuration files
var (
	SystemPath      = getEnv("SYSTEM_PATH")
	MongoDBHost     = getEnv("MONGODB_HOST")
	MongoDBName     = getEnv("MONGODB_NAME")
	MongoDBUsername = getEnv("MONGODB_USERNAME")
	MongoDBPassword = getEnv("MONGODB_PASSWORD")
	SecretKey       = getEnv("SECRET_KEY")
	JWTKey          = []byte(getEnv("JWT_KEY"))
	Env             = getEnv("ENV")
)

// IsProduction :
func IsProduction() bool {
	return Env == "production"
}

// IsSandbox :
func IsSandbox() bool {
	return Env == "sandbox"
}

// IsDevelopment :
func IsDevelopment() bool {
	return Env == "development"
}

func getEnv(name string) string {
	env := os.Getenv(name)
	if env == "" {
		panic("failed to get env for " + name)
	}
	return env
}
