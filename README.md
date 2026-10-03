### run dlv debug

```shell
dlv debug --headless --listen=:2345 --accept-multiclient --continue --only-same-user=false --log
```

### Create .env file

```bash
touch .env
```

with these envs

```bash
GH_TOKEN=1
TRELLO_TOKEN=1
TRELLO_KEY=1
```

### Running tests

```bash
go test
```

### running api

it starts a server at `https://localhost:12769`

```bash
go run cmd/api/main.go
```
