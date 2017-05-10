# etcdv2

### test

Run etcd in a docker container.
```
docker run --rm -p 0.0.0.0:2379:2379 --name etcdv2 quay.io/coreos/etcd:v2.3.7 -advertise-client-urls http://0.0.0.0:2379 -listen-client-urls http://0.0.0.0:2379
```

Run the integration tests.
```
GOOS=darwin; GOARCH=amd64 go test -tags integration ./storage/etcdv2
```

Check the keyspace within etcd.
```
docker run --rm --net host --name etcdctl quay.io/coreos/etcd:v3.0.16 etcdctl ls --recursive
```

Cleanup the keyspace within etcd.
```
docker run --rm --net host --name etcdctl quay.io/coreos/etcd:v3.0.16 etcdctl rm --recursive /foo
```
