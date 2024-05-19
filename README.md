# Go Service Template

Basic bootstrap for writing http service in Golang, without usage of any http framework or router. Initial idea was to use only packages from standard library. The only exception right now is `golang.org/x/sync` package, which is used for [errgroup.Group](https://pkg.go.dev/golang.org/x/sync/errgroup).