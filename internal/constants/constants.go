package constants

import (
	"strings"

	"kira-url/internal/env"
)

const (
	BASE_TTL int = 300
)

var BaseDomain = ensureTrailingSlash(env.GetEnvString("SERVER_DOMAIN", "http://localhost:3536/"))

// ensureTrailingSlash keeps short URL concatenation safe: BaseDomain is
// prepended directly to the stored slug, so the separator must be present.
func ensureTrailingSlash(domain string) string {
	if domain == "" || strings.HasSuffix(domain, "/") {
		return domain
	}
	return domain + "/"
}
