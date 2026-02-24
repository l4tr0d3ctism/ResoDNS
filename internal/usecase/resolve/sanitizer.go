package resolve

import "strings"

func DefaultSanitizer(domain string) string {
	domain = strings.ToLower(domain)
	domain = strings.TrimPrefix(domain, "*.")
	for i := 0; i < len(domain); i++ {
		c := domain[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' {
			continue
		}
		return ""
	}
	return domain
}
