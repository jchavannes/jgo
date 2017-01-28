## JGo Web

- **Port**: Port to bind to
- **Sessions**: Set JGoSession cookie
- **TemplateDir**: Go HTML templates
- **StaticDir**: Static assets
- **Routes**: Custom endpoints

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
- Persistent sessions
- Both templates and static assets in root
