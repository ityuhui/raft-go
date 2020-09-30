# raft-go
My implementation of raft by golang


## generate gprc code

```
$ protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative raft_rpc/raft.proto
```