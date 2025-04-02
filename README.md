# hi-low-mini
Repo for hi-low mini (personal project)


## Setup
#### For all:
Generate protos
```
buf generate
<!-- to install refer to https://connectrpc.com/docs/go/getting-started -->
```


#### Frontend:
Install JS dependencies


#### Backend:
```
go mod tidy
go mod vendor
```

#### Running locally

Start up the server by running
```
go run ./backend/cmd/server/main.go
```

You can send curl requests to test. 
E.g.
```
curl \
    --header "Content-Type: application/json" \
    --data '{}' \
    http://localhost:8080/game_engine.v1.GameEngineService/GetGame
{"game":{"id":"000", "timestamp":"2025-03-19T18:47:11.624953Z", "pool":{"cards":[{"value":"1", "type":"CARD_TYPE_NUMBER"}]}}}
```