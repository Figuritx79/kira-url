package constants

import "kira-url/internal/env"

const (
	BASE_TTL int = 300
)

var BaseDomain = env.GetEnvString("DOMAIN", "http://localhost:3536/")
