package resolve

import (
	"bufio"
	"io"
	"strings"

	"resodns/pkg/wildcarder"
	"github.com/d3mondev/resolvermt"
)

type CacheReader struct {
	reader  io.ReadCloser
	scanner *bufio.Scanner
}

func NewCacheReader(r io.ReadCloser) *CacheReader {
	return &CacheReader{reader: r, scanner: bufio.NewScanner(r)}
}

func (c *CacheReader) Read(w io.Writer, cache *wildcarder.DNSCache, maxCount int) (count int, err error) {
	const (
		stateNewAnswerSection = iota
		stateSaveAnswer
		stateSkip
	)
	var curDomain string
	var curState int
	var domainSaved bool
	var found int
	for c.scanner.Scan() {
		line := c.scanner.Text()
		if line == "" {
			curState = stateNewAnswerSection
			if maxCount > 0 && found == maxCount {
				break
			}
			continue
		}
		switch curState {
		case stateNewAnswerSection:
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				curState = stateSkip
				continue
			}
			domain := strings.TrimSuffix(parts[0], ".")
			if domain == "" {
				curState = stateSkip
				continue
			}
			curDomain = domain
			domainSaved = false
			curState = stateSaveAnswer
			fallthrough
		case stateSaveAnswer:
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				curState = stateSkip
				continue
			}
			domain := curDomain
			rrtypeStr := parts[1]
			answer := parts[2]
			var rrtype resolvermt.RRtype
			switch rrtypeStr {
			case "A":
				rrtype = resolvermt.TypeA
			case "AAAA":
				rrtype = resolvermt.TypeAAAA
			case "CNAME":
				answer = strings.TrimSuffix(answer, ".")
				rrtype = resolvermt.TypeCNAME
			default:
				continue
			}
			if !domainSaved {
				found++
				domainSaved = true
				if w != nil {
					w.Write([]byte(domain + "\n"))
				}
			}
			if cache != nil {
				cache.Add(domain, []wildcarder.DNSAnswer{{Type: rrtype, Answer: answer}})
			}
			if cache == nil && w == nil {
				curState = stateSkip
			}
		case stateSkip:
		}
	}
	return found, c.scanner.Err()
}

func (c *CacheReader) Close() error {
	return c.reader.Close()
}
