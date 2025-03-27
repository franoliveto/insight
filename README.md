This is currently a work in progress.

Package insight provides idiomatic Go APIs for accessing [deps.dev](https://deps.dev)
public API.

### How to use it

First create a client.

```go
client := insight.NewClient()
```

Then use that client to interact with the API.

```go
key := insight.VersionKey{System: "npm", Name: "react", Version: "18.2.0"}
deps, err := client.GetDependencies(key)
if err != nil {
    log.Fatal(err)
}
```