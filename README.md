## JGo Web

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

---

#### Todo

- Web sockets
- Authentication example
