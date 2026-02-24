package template

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"resodns/pkg/fileoperation"
)

// Segment represents either a literal string or a placeholder to be expanded.
type Segment struct {
	Literal string   // non-empty for literal segment
	Values  []string // non-empty for placeholder: list of values to expand
}

// Parse parses a template string into a list of segments. Placeholders are
// [start-end], [start-end:step], or [file:path].
func Parse(template string) ([]Segment, error) {
	var segments []Segment
	s := template

	for {
		idx := strings.Index(s, "[")
		if idx < 0 {
			if len(s) > 0 {
				segments = append(segments, Segment{Literal: s})
			}
			break
		}
		if idx > 0 {
			segments = append(segments, Segment{Literal: s[:idx]})
		}
		end := strings.Index(s[idx:], "]")
		if end < 0 {
			return nil, fmt.Errorf("unclosed '[' in template")
		}
		end += idx
		content := s[idx+1 : end]
		seg, err := parsePlaceholder(content)
		if err != nil {
			return nil, err
		}
		segments = append(segments, *seg)
		s = s[end+1:]
	}

	return segments, nil
}

// parsePlaceholder parses the content inside [...] and returns a segment with Values set.
func parsePlaceholder(content string) (*Segment, error) {
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return nil, fmt.Errorf("empty placeholder []")
	}

	// [file:path]
	if strings.HasPrefix(content, "file:") {
		path := strings.TrimSpace(content[5:])
		if path == "" {
			return nil, fmt.Errorf("file placeholder has empty path")
		}
		lines, err := fileoperation.ReadLines(path)
		if err != nil {
			return nil, fmt.Errorf("reading [file:%s]: %w", path, err)
		}
		// trim and drop empty lines
		var values []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				values = append(values, line)
			}
		}
		return &Segment{Values: values}, nil
	}

	// [start-end] or [start-end:step]
	rangeRegex := regexp.MustCompile(`^(\d+)-(\d+)(?::(\d+))?$`)
	m := rangeRegex.FindStringSubmatch(content)
	if m != nil {
		start, _ := strconv.Atoi(m[1])
		end, _ := strconv.Atoi(m[2])
		step := 1
		if m[3] != "" {
			step, _ = strconv.Atoi(m[3])
		}
		if step <= 0 {
			return nil, fmt.Errorf("invalid range step %q (must be positive)", content)
		}
		if start > end {
			return nil, fmt.Errorf("invalid range %q (start > end)", content)
		}
		var values []string
		for i := start; i <= end; i += step {
			values = append(values, strconv.Itoa(i))
		}
		return &Segment{Values: values}, nil
	}

	return nil, fmt.Errorf("unknown placeholder %q (use [start-end], [start-end:step], or [file:path])", content)
}

// Expand parses the template and returns the Cartesian product of all placeholders
// with literals concatenated. If a placeholder has no values (e.g. empty file),
// the result is empty.
func Expand(template string) ([]string, error) {
	segments, err := Parse(template)
	if err != nil {
		return nil, err
	}
	return expandSegments(segments), nil
}

func expandSegments(segments []Segment) []string {
	// Collect each segment's values; literals become a single value.
	var options [][]string
	for _, seg := range segments {
		if seg.Literal != "" {
			options = append(options, []string{seg.Literal})
		} else {
			if len(seg.Values) == 0 {
				return nil // empty placeholder => no results
			}
			options = append(options, seg.Values)
		}
	}
	if len(options) == 0 {
		return nil
	}
	// Cartesian product
	return cartesian(options)
}

func cartesian(options [][]string) []string {
	if len(options) == 0 {
		return nil
	}
	n := 1
	for _, o := range options {
		n *= len(o)
	}
	result := make([]string, 0, n)
	indices := make([]int, len(options))
	for {
		var b strings.Builder
		for i, opt := range options {
			b.WriteString(opt[indices[i]])
		}
		result = append(result, b.String())

		// next combination
		for j := len(indices) - 1; j >= 0; j-- {
			indices[j]++
			if indices[j] < len(options[j]) {
				break
			}
			indices[j] = 0
			if j == 0 {
				return result
			}
		}
	}
}

// ExpandWithMax is like Expand but returns at most max domains (max <= 0 means unlimited).
func ExpandWithMax(template string, max int) ([]string, error) {
	list, err := Expand(template)
	if err != nil {
		return nil, err
	}
	if max > 0 && len(list) > max {
		list = list[:max]
	}
	return list, nil
}

// ExpandToFile expands the template and writes the result to path, one domain per line.
// If max > 0, at most max domains are written. Returns the number of lines written.
func ExpandToFile(template string, path string, max int) (int, error) {
	list, err := ExpandWithMax(template, max)
	if err != nil {
		return 0, err
	}
	if err := fileoperation.WriteLines(list, path); err != nil {
		return 0, err
	}
	return len(list), nil
}
