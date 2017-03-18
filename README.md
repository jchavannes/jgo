## JGo Web

- **Port**: Port to bind to
- **Routes**: Custom endpoints
- **Sessions**: Set JGoSession cookie
- **SessionKey**: Key used to generate and validate session tokens
- **StaticDir**: Static assets
- **TemplateDir**: Go HTML templates
- **InitResponse**: Function called before processing every response

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
