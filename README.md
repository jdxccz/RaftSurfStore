# RaftSurfStore

## Description

This is a cloud-based multi-user synchronous storage system, which is divided into meta stores for storing versions and hashes and block stores for storing data. It built the Raft mechanism to ensure fault-tolerant meta store and the Chord mechanism to ensure scalable and distributed block store.

## Usage

1. Start Server
```console
make run-server
```

2. Start Client
```console
make run-client $BASEDIR
```

3. Debug
```console
protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/SurfStore.proto
go run cmd/server/main.go -f $CONFIGFILE -m $MODE -i $ID
go run cmd/client/main.go -f $CONFIGFILE -b $BASEDIR
```

## EXAMPLE

You should execute the following commands:
1. start server
```console
make run-server
```
2. sync db1 & cloud
```console
make run-client db1
```

3. sync db2 & cloud
```console
make run-client db2
```

4. sync db1 & cloud again
```console
make run-client db1
```

You can find that the files inside the folders db1 and db2 have been synchronized.