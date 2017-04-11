# etcd

### test

Run etcd in a docker container.
```
docker run --rm -p 0.0.0.0:2380:2380 -p 0.0.0.0:2379:2379 --name etcd quay.io/coreos/etcd:v3.0.16 etcd --initial-cluster "default=http://0.0.0.0:2380" --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379 --initial-advertise-peer-urls http://0.0.0.0:2380 --listen-peer-urls http://0.0.0.0:2380
```

Run the integration tests.
```
GOOS=darwin; GOARCH=amd64 go test -tags integration ./storage/etcd
```

Check the keyspace within etcd.
```
docker run --rm -e ETCDCTL_API=3 --net host --name etcdctl quay.io/coreos/etcd:v3.0.16 etcdctl get --prefix --keys-only ""
```

Cleanup the keyspace within etcd.
```
docker run --rm -e ETCDCTL_API=3 --net host --name etcdctl quay.io/coreos/etcd:v3.0.16 etcdctl del --prefix ""
```
