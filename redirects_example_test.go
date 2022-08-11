package redirects

import (
	"encoding/json"
	"os"
)

func Example() {
	h := Must(ParseString(`
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
  `))

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(h)
	// Output:
	// 	[
	//   {
	//     "From": "/home",
	//     "To": "/",
	//     "Status": 301
	//   },
	//   {
	//     "From": "/blog/my-post.php",
	//     "To": "/blog/my-post",
	//     "Status": 301
	//   },
	//   {
	//     "From": "/news",
	//     "To": "/blog",
	//     "Status": 301
	//   },
	//   {
	//     "From": "/google",
	//     "To": "https://www.google.com",
	//     "Status": 301
	//   },
	//   {
	//     "From": "/home",
	//     "To": "/",
	//     "Status": 301
	//   },
	//   {
	//     "From": "/my-redirect",
	//     "To": "/",
	//     "Status": 302
	//   },
	//   {
	//     "From": "/pass-through",
	//     "To": "/index.html",
	//     "Status": 200
	//   },
	//   {
	//     "From": "/ecommerce",
	//     "To": "/store-closed",
	//     "Status": 404
	//   },
	//   {
	//     "From": "/*",
	//     "To": "/index.html",
	//     "Status": 200
	//   },
	//   {
	//     "From": "/api/*",
	//     "To": "https://api.example.com/:splat",
	//     "Status": 200
	//   }
	// ]
}
