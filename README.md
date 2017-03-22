## JGo Web

#### Features

- Simple routing
- Automatic template rendering based on request path
- Response object with support for common functionality, e.g.:
  - Getting form and header data
  - Redirects
- Sessions with validatable tokens
- Enforceable CSRF tokens
- Websockets (coming soon)
- Example app with:
  - User accounts
  - Shared HTML templates (e.g. header)
  - SQLite database
  - Chatroom

#### Options

- **Port** int: Port to bind to
- **Routes** []*web.Route: Custom endpoints
- **Sessions** bool: Set JGoSession cookie
- **SessionKey** string: Key used to generate and validate session tokens
- **StaticDir** string: Static assets
- **TemplateDir** string: Go HTML templates
- **InitResponse** func(r *web.Response): Function called before processing every response

#### Example usage

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
