package redirects_test

import (
	"encoding/json"
	"os"

	redirects "github.com/ipfs-shipyard/go-ipfs-redirects"
)

func Example() {
	h := redirects.Must(redirects.ParseString(`
		# Implicit 301 redirects
		/home              /
		/blog/my-post.php  /blog/my-post
		/news              /blog
		/google            https://www.google.com

		# Redirect with a 301
		/home         /              301

		# Redirect with a 302
		/my-redirect  /              302

		# Rewrite a path
		/pass-through /index.html    200

		# Show a custom 404 for this path
		/ecommerce    /store-closed  404

		# Single page app rewrite
		/*    /index.html   200

		# Proxying
		/api/*  https://api.example.com/:splat  200

		# Forcing
		/app/*  /app/index.html  200!
  `))

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(h)
	// Output:
	// 	[
	//   {
	//     "From": "/home",
	//     "To": "/",
	//     "Status": 301,
	//     "Force": false
	//   },
	//   {
	//     "From": "/blog/my-post.php",
	//     "To": "/blog/my-post",
	//     "Status": 301,
	//     "Force": false
	//   },
	//   {
	//     "From": "/news",
	//     "To": "/blog",
	//     "Status": 301,
	//     "Force": false
	//   },
	//   {
	//     "From": "/google",
	//     "To": "https://www.google.com",
	//     "Status": 301,
	//     "Force": false
	//   },
	//   {
	//     "From": "/home",
	//     "To": "/",
	//     "Status": 301,
	//     "Force": false
	//   },
	//   {
	//     "From": "/my-redirect",
	//     "To": "/",
	//     "Status": 302,
	//     "Force": false
	//   },
	//   {
	//     "From": "/pass-through",
	//     "To": "/index.html",
	//     "Status": 200,
	//     "Force": false
	//   },
	//   {
	//     "From": "/ecommerce",
	//     "To": "/store-closed",
	//     "Status": 404,
	//     "Force": false
	//   },
	//   {
	//     "From": "/*",
	//     "To": "/index.html",
	//     "Status": 200,
	//     "Force": false
	//   },
	//   {
	//     "From": "/api/*",
	//     "To": "https://api.example.com/:splat",
	//     "Status": 200,
	//     "Force": false
	//   },
	//   {
	//     "From": "/app/*",
	//     "To": "/app/index.html",
	//     "Status": 200,
	//     "Force": true
	//   }
	// ]
}
