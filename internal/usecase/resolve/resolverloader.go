package resolve

import (
	"bufio"
	"io"
	"os"
	"strings"

	"resodns/internal/app/ctx"
)

type DefaultResolverLoader struct{}

func NewDefaultResolverFileLoader() *DefaultResolverLoader {
	return &DefaultResolverLoader{}
}

func (l *DefaultResolverLoader) Load(c *ctx.Ctx, filename string) error {
	if filename == "" {
		return nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	resolvers, err := loadResolvers(file)
	if err != nil {
		return err
	}
	if len(resolvers) > 0 {
		c.Options.TrustedResolvers = resolvers
	}
	return nil
}

func loadResolvers(r io.Reader) ([]string, error) {
	var resolvers []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		resolvers = append(resolvers, s)
	}
	return resolvers, scanner.Err()
}
