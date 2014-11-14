package main

import (
  "flag"
  "log"
  . "etcdfs"
  "github.com/hanwen/go-fuse/fuse/pathfs"
  "github.com/hanwen/go-fuse/fuse/nodefs"
)

func main() {
  flag.Parse()
  if len(flag.Args()) < 2 {
    log.Fatal("Usage:\n  etcd-fs MOUNTPOINT ETCDENDPOINT <ROOT>")
  }

  rootPath := "/"
  if flag.Arg(2) != nil {
    rootPath = flag.Arg(2)
  }

  etcdFs := EtcdFs{FileSystem: pathfs.NewDefaultFileSystem(), EtcdEndpoint: flag.Arg(1), EtcdRootPath: rootPath }
  etcdFs.NewEtcdClient().mkdir(rootPath)
  nfs := pathfs.NewPathNodeFs(&etcdFs, nil)
  server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil)
  if err != nil {
    log.Fatalf("Mount fail: %v\n", err)
  }
  server.Serve()
}
