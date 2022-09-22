// Package redirects provides Netlify style _redirects file format parsing.
package redirects

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/ucarion/urlpath"
)

// 64 KiB
const maxFileSizeInBytes = 65536

// A Rule represents a single redirection or rewrite rule.
type Rule struct {
	// From is the path which is matched to perform the rule.
	From string

	// To is the destination which may be relative, or absolute
	// in order to proxy the request to another URL.
	To string

	// Status is one of the following:
	//
	// - 3xx a redirect
	// - 200 a rewrite
	// - defaults to 301 redirect
	//
	// When proxying this field is ignored.
	//
	Status int
}

// IsRewrite returns true if the rule represents a rewrite (status 200).
func (r *Rule) IsRewrite() bool {
	return r.Status == 200
}

// IsProxy returns true if it's a proxy rule (aka contains a hostname).
func (r *Rule) IsProxy() bool {
	u, err := url.Parse(r.To)
	if err != nil {
		return false
	}

	return u.Host != ""
}

// MatchAndExpandPlaceholders expands placeholders in `r.To` and returns true if the provided path matches.
// Otherwise it returns false.
func (r *Rule) MatchAndExpandPlaceholders(urlPath string) bool {
	// get rule.From, trim trailing slash, ...
	fromPath := urlpath.New(strings.TrimSuffix(r.From, "/"))
	match, ok := fromPath.Match(urlPath)

	if !ok {
		return false
	}

	// We have a match!  Perform substitution and return the updated rule
	toPath := r.To
	toPath = replacePlaceholders(toPath, match)
	toPath = replaceSplat(toPath, match)

	r.To = toPath

	return true
}

func replacePlaceholders(to string, match urlpath.Match) string {
	if len(match.Params) > 0 {
		for key, value := range match.Params {
			to = strings.ReplaceAll(to, ":"+key, value)
		}
	}

	return to
}

func replaceSplat(to string, match urlpath.Match) string {
	return strings.ReplaceAll(to, ":splat", match.Trailing)
}

// Must parse utility.
func Must(v []Rule, err error) []Rule {
	if err != nil {
		panic(err)
	}

	return v
}

// Parse the given reader.
func Parse(r io.Reader) (rules []Rule, err error) {
	// not too large
	b, err := bufio.NewReaderSize(r, maxFileSizeInBytes+1).Peek(maxFileSizeInBytes + 1)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if len(b) > maxFileSizeInBytes {
		return nil, fmt.Errorf("redirects file size cannot exceed %d bytes", maxFileSizeInBytes)
	}

	s := bufio.NewScanner(bytes.NewReader(b))
	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		// empty
		if line == "" {
			continue
		}

		// comment
		if strings.HasPrefix(line, "#") {
			continue
		}

		// fields
		fields := strings.Fields(line)

		// missing dst
		if len(fields) <= 1 {
			return nil, fmt.Errorf("missing 'to' path: %q", line)
		}

		if len(fields) > 3 {
			return nil, fmt.Errorf("must match format 'from to [status]'")
		}

		// src and dst
		rule := Rule{
			From:   fields[0],
			To:     fields[1],
			Status: 301,
		}

		// from
		if !strings.HasPrefix(rule.From, "/") {
			return nil, fmt.Errorf("'from' path must begin with '/'")
		}

		if strings.Contains(rule.From, "*") && !strings.HasSuffix(rule.From, "*") {
			return nil, fmt.Errorf("'from' path can only end with splat")
		}

		// to
		if !strings.HasPrefix(rule.To, "/") {
			_, err := url.Parse(rule.To)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid 'to' path")
			}
		}

		// status
		if len(fields) > 2 {
			code, err := parseStatus(fields[2])
			if err != nil {
				return nil, errors.Wrapf(err, "parsing status %q", fields[2])
			}

			rule.Status = code
		}

		rules = append(rules, rule)
	}

	err = s.Err()
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// ParseString parses the given string.
func ParseString(s string) ([]Rule, error) {
	return Parse(strings.NewReader(s))
}

// parseStatus returns the status code.
func parseStatus(s string) (code int, err error) {
	if strings.HasSuffix(s, "!") {
		// See https://docs.netlify.com/routing/redirects/rewrites-proxies/#shadowing
		return 0, fmt.Errorf("forced redirects (or \"shadowing\") are not supported by IPFS gateways")
	}

	code, err = strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	if !isValidStatusCode(code) {
		return 0, fmt.Errorf("status code %d is not supported", code)
	}

	return code, nil
}

func isValidStatusCode(status int) bool {
	switch status {
	case 200, 301, 302, 303, 307, 308, 404, 410, 451:
		return true
	}
	return false
}
