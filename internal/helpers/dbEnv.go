package helpers

import (
	"os"

	"github.com/g-ton/stori-candidate/env"
)

func GetEnvDB() env.EnvApp {
	/* env config reading from
	.env.development
	.env.testaws
	*/
	envVar := os.Getenv("GIN_MODE")
	if envVar == "release" {
		return env.GetEnv(".env.testaws")
	}
	return env.GetEnv(".env.development")
}
