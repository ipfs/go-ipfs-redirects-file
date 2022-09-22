package redirects

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/tj/assert"
)

func TestRuleIsProxy(t *testing.T) {
	t.Run("without host", func(t *testing.T) {
		r := Rule{
			From: "/blog",
			To:   "/blog/engineering",
		}

		assert.False(t, r.IsProxy())
	})

	t.Run("with host", func(t *testing.T) {
		r := Rule{
			From: "/blog",
			To:   "https://blog.apex.sh",
		}

		assert.True(t, r.IsProxy())
	})
}

func TestRuleIsRewrite(t *testing.T) {
	t.Run("with 3xx", func(t *testing.T) {
		r := Rule{
			From:   "/blog",
			To:     "/blog/engineering",
			Status: 302,
		}

		assert.False(t, r.IsRewrite())
	})

	t.Run("with 200", func(t *testing.T) {
		r := Rule{
			From:   "/blog",
			To:     "/blog/engineering",
			Status: 200,
		}

		assert.True(t, r.IsRewrite())
	})

	t.Run("with 0", func(t *testing.T) {
		r := Rule{
			From: "/blog",
			To:   "/blog/engineering",
		}

		assert.False(t, r.IsRewrite())
	})
}

func TestParse(t *testing.T) {
	t.Run("with illegal force", func(t *testing.T) {
		_, err := Parse(strings.NewReader(`
		/home / 301!
		`))

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "forced redirects")
	})

	t.Run("with illegal code", func(t *testing.T) {
		_, err := Parse(strings.NewReader(`
		/home / 42
		`))

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status code 42 is not supported")
	})

	t.Run("with too large file", func(t *testing.T) {
		// create a file larger than 64 KiB, using valid rules so the only possible error is the size
		line := "/from /to 301"
		bytesPerLine := len(line)
		totalBytes := 0

		var b bytes.Buffer
		for totalBytes <= maxFileSizeInBytes {
			totalBytes += bytesPerLine
			b.WriteString(line + "\n")
		}
		text := b.String()

		_, err := ParseString(text)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redirects file size cannot exceed")
	})
}

func FuzzParse(f *testing.F) {
	testcases := []string{"/a /b 999\n",
		"/redirect-one /one.html\n/301-redirect-one /one.html 301\n/302-redirect-two /two.html 302\n/200-index /index.html 200\n/posts/:year/:month/:day/:title /articles/:year/:month/:day/:title 301\n/splat/* /redirected-splat/:splat 301\n/not-found/* /404.html 404\n/* /index.html 200\n",
		"/a /b 301!\n",
		"/a200 /b200 200\n/a301 /b301 301\n/a302 /b302 302\n/a303 /b303 303\n/a307 /b307 307\n/a308 /b308 308\n/a404 /b404 404\n/a410 /b410 410\n/a451 /b451 451\n",
		"hello\n", "/redirect-one /one.html\r\n/200-index /index.html 200\r\n", "a b 2\nc   d 42", "/a/*/b blah", "/from https://example.com 200\n/a/:blah/yeah /b/:blah/yeah"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rules, err := ParseString(orig)
		if err != nil {
			if len(rules) > 0 {
				t.Errorf("should not return rules on error")
			}
		}

		s := bufio.NewScanner(strings.NewReader(orig))

		for s.Scan() {
			line := strings.TrimSpace(s.Text())
			fields := strings.Fields(line)

			// Skip comments so we don't have to special case
			if strings.HasPrefix(line, "#") {
				continue
			}

			if err == nil && len(fields) < 2 && line != "" {
				t.Errorf("should error with less than 2 fields.  orig='%v'", orig)
				continue
			}

			if err == nil && len(fields) > 3 {
				t.Errorf("should error with more than 3 fields.  orig='%v'", orig)
				continue
			}

			if err == nil && len(fields) > 0 && !strings.HasPrefix(fields[0], "/") {
				t.Errorf("should error for from path not starting with '/'.  orig=%v", orig)
				continue
			}

			// we already handled these cases
			if len(fields) < 3 {
				continue
			}

			if err == nil && strings.Contains(fields[0], "*") && !strings.HasSuffix(fields[0], "*") {
				t.Errorf("asterisk in from not at end should error.  orig=%v", orig)
				continue
			}

			if err == nil && strings.HasSuffix(fields[2], "!") {
				t.Errorf("should error for forced redirects.  orig=%v, err=%v", orig, err)
				continue
			}

			if err == nil {
				for _, r := range rules {
					if !isValidStatusCode(r.Status) {
						t.Errorf("should error for invalid status code.  orig=%v", orig)
					}
				}
			}
		}
	})
}
