package redirects

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tj/assert"
)

func TestRule_IsProxy(t *testing.T) {
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

func TestRule_IsRewrite(t *testing.T) {
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

func Test_Parse(t *testing.T) {
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
