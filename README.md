## JGo Web

#### Features

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

#### Options

- **Port** _int_: Port to bind to
- **Routes** _[]*web.Route_: Custom endpoints
- **Sessions** _bool_: Set JGoSession cookie
- **SessionKey** _string_: Key used to generate and validate session tokens
- **StaticDir** _string_: Static assets
- **TemplateDir** _string_: Go HTML templates
- **InitResponse** _func(r *web.Response)_: Function called before processing every response

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
