### The repo

This is a repo to run some integrations between github and trello

example

- Get all pull requests from a github user and create cards on trello with name, pr status, reviews
- Get last version of changelog github repo file
- Get reviews from a pull request and count

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
go test internal/client/trello/*
go test internal/client/github/*
...
```

### running api

it starts a server at `https://localhost:12769`

```bash
go run cmd/api/main.go
```
