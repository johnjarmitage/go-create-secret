# go-create-secret
creation of docker container that creates a secret within k8s

```
go get -d -v ./...
GOOS=linux go build -o secret-maker .
docker build -t secret_maker:1.0.0 .
```
