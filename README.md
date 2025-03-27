This is currently a work in progress.

Package insight provides idiomatic Go APIs for accessing [deps.dev](https://deps.dev)
public API.

First create a client.

```go
client := insight.NewClient()
```

Then use that client to interact with the API.

```go
deps, err := client.GetDependencies("npm", "react", "18.2.0")
if err != nil {
    log.Fatal(err)
}
```