## JGo Web

- **Port**: Port to bind to
- **EnableSessions**: Set JGoSession cookie
- **TemplateDirectory**: Go HTML templates
- **StaticDirectory**: Static assets
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
        EnableSessions: true,
        TemplateDirectory: "templates",
        StaticDirectory: "public",
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
