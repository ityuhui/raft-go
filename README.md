# raft-go
My implementation of raft by golang


## generate gprc code

```
$ protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative raft_rpc/raft.proto
```

### Test

```
$ raft-go.exe --Iam=localhost:5050 --Peers=localhost:5051,localhost:5052

$ raft-go.exe --Iam=localhost:5051 --Peers=localhost:5052,localhost:5050

$ raft-go.exe --Iam=localhost:5052 --Peers=localhost:5050,localhost:5051
```