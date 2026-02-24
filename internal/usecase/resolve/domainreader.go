package resolve

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"resodns/pkg/procreader"
)

type DomainReader struct {
	source          io.ReadCloser
	sourceScanner   *bufio.Scanner
	subdomainReader *procreader.ProcReader
	domains         []string
	sanitizer       DomainSanitizer
}

var _ io.Reader = (*DomainReader)(nil)

type DomainSanitizer func(domain string) string

func NewDomainReader(source io.ReadCloser, domains []string, sanitizer DomainSanitizer) *DomainReader {
	r := &DomainReader{
		source:        source,
		sourceScanner: bufio.NewScanner(source),
		domains:       domains,
		sanitizer:     sanitizer,
	}
	r.subdomainReader = procreader.New(r.nextSubdomains)
	return r
}

func (r *DomainReader) Read(p []byte) (int, error) {
	return r.subdomainReader.Read(p)
}

func (r *DomainReader) nextSubdomains(size int) ([]byte, error) {
	if !r.sourceScanner.Scan() {
		r.source.Close()
		if err := r.sourceScanner.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}
	var output bytes.Buffer
	word := r.sourceScanner.Text()
	if len(r.domains) == 0 {
		domain := r.processDomain(word)
		output.WriteString(domain)
	} else {
		for _, domain := range r.domains {
			if strings.ContainsRune(domain, '*') {
				domain = strings.ReplaceAll(domain, "*", word)
			} else {
				domain = fmt.Sprintf("%s.%s", word, domain)
			}
			output.WriteString(r.processDomain(domain))
		}
	}
	return output.Bytes(), nil
}

func (r *DomainReader) processDomain(domain string) string {
	if r.sanitizer != nil {
		domain = r.sanitizer(domain)
	}
	return domain + "\n"
}
