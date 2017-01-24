## JGo Web

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
        TemplateDirectory: "templates",
        StaticDirectory: "public",
        EnableSessions: true,
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
