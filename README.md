# IPFS `_redirects` File Parser

This is a parser for the IPFS Web Gateway's `_redirects` file format.

## Specification

Follow specification work at https://github.com/ipfs/specs/pull/290

## Format
Currently only supports `from`, `to` and `status`.

```
from to [status]
```

## Example

```sh
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

# Show a custom 404 for everything under this path
/ecommerce/*  /store-closed  404

# Single page app rewrite
/*    /index.html   200

# Proxying
/api/*  https://api.example.com/:splat  200
```

---

## Credit
This project was forked from [tj/go-redirects](https://github.com/tj/go-redirects).  Thank you TJ for the initial work. 🙏
