# raft-go
An implementation of raft


## Generate gprc code

```
$ protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative raft_rpc/raft.proto
```

### Build

#### Build the daemon
```shell
cd raft-deamon
go build
```

#### Build the client
```shell
cd raft-client
go build
```

### Usage

#### Start raft-daemon
```
raft-daemon.exe --I=localhost:5050 --Peers=localhost:5051,localhost:5052

raft-daemon.exe --I=localhost:5051 --Peers=localhost:5052,localhost:5050

raft-daemon.exe --I=localhost:5052 --Peers=localhost:5050,localhost:5051
```

#### Use raft-client

```
raft-client.exe --server localhost:5050 --set a=10

raft-client.exe --server localhost:5050 --get a
```
