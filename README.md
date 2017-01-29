## JGo Web

- **Port**: Port to bind to
- **Routes**: Custom endpoints
- **Sessions**: Set JGoSession cookie
- **StaticDir**: Static assets
- **TemplateDir**: Go HTML templates

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
            Handler: func(r *web.Request) {
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
- Persistent sessions
- Both templates and static assets in root
