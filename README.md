This is currently a work in progress.

Package insight provides idiomatic Go APIs for accessing [deps.dev API v3](https://docs.deps.dev/api/v3/).

First create a client.

```go
client := insight.NewClient()
```

Then use that client to interact with the API.

```go
ctx := context.Background()
deps, err := client.GetDependencies(ctx, "npm", "react", "18.2.0")
if err != nil {
    log.Fatal(err)
}
```