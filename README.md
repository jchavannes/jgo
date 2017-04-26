# JGo Web

## Features

- Simple routing
- Automatic template rendering based on request path
- Response object with support for common functionality, e.g.:
  - Getting form and header data
  - Redirects
- Sessions with validatable tokens
- Enforceable CSRF tokens
- WebSockets
- Example app with:
  - User accounts
  - Shared HTML templates (e.g. header)
  - SQLite database
  - Chatroom

## Docs

GoDocs: [https://godoc.org/github.com/jchavannes/jgo](https://godoc.org/github.com/jchavannes/jgo)

## Usage

```go
package main

import (
    "github.com/jchavannes/jgo/web"
)

func main() {
    server := web.Server{
        Port: 80,
        Routes: []web.Route{{
            Pattern: "/hello",
            Handler: func(r *web.Response) {
                r.Write("world")
            },
        }},
    }
    server.Run()
}
```

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
