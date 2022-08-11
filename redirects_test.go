package redirects

import (
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

func TestParse(t *testing.T) {
	t.Run("with illegal force", func(t *testing.T) {
		_, err := Parse(strings.NewReader(`
		/home / 301!
		`))

		assert.Error(t, err, "forced redirects should return an error")
	})
}
