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

- **AllowedExtensions** _[]string_: Allowed file extensions for static files (`*` for all)
- **DisableAutoRender** _bool_: Disables automatically rendering templates
- **Port** _int_: Port to bind to
- **PreHandler** _func(*web.Response)_: Function called before processing responses
- **Routes** _[]web.Route_: Custom endpoints
- **SessionKey** _string_: Key used to generate and validate session tokens
- **StaticFilesDir** _string_: Static assets
- **TemplatesDir** _string_: Go HTML templates
- **UseSessions** _bool_: Set JGoSession cookie

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
