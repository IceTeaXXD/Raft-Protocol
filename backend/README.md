# IF3230 - Sistem Paralel dan Terdistribusi
> Tugas Besar - Raft Protocol (BE)

## How to run
1. Run the server
```
go run cmd/server/main.go
```
Arguemnts:
- `-port` (default: 8080) : Port number for the server
- `-host` (default: localhost) : Host for the server


## Running Examples
1. Leader
```
go run cmd/server/main.go -host=localhost -port=8080
```

2. Followers
```
go run ./cmd/server/main.go -leaderHost=localhost -leaderPort=8080 -host=localhost -port=8081
go run ./cmd/server/main.go -leaderHost=localhost -leaderPort=8080 -host=localhost -port=8082
go run ./cmd/server/main.go -leaderHost=localhost -leaderPort=8080 -host=localhost -port=8083
go run ./cmd/server/main.go -leaderHost=localhost -leaderPort=8080 -host=localhost -port=8084
go run ./cmd/server/main.go -leaderHost=localhost -leaderPort=8080 -host=localhost -port=8085
```

3. Client
```
go run cmd/client/main.go -host=localhost -port=8080
```

4. Unit Test (BONUS)
```
cd internal/store
go test -v
```