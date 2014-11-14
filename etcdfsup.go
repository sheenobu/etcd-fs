package main

import (
  "flag"
  "log"
  "os"
  "os/exec"
  "sync"
  "time"
  . "etcdfs"
  "github.com/hanwen/go-fuse/fuse/pathfs"
  "github.com/hanwen/go-fuse/fuse/nodefs"
  "github.com/landjur/go-uuid"
)

/*
 * Uploads a given directory into etcd.
 */
func main() {
  flag.Parse()
  if len(flag.Args()) < 2 {
    log.Fatal("Usage:\n  etcd-fs ETCDENDPOINT SRCROOT <DESTROOT> ")
  }

  srcRoot := flag.Arg(1)

  rootPath := "/"
  if flag.Arg(2) != "" {
    rootPath = flag.Arg(2)
  }

  tmpFile := os.TempDir()
  uuid, err := uuid.NewV4()
  dir := tmpFile + "/" + uuid.String()

  os.MkdirAll(dir,os.ModeDir)

  etcdFs := EtcdFs{FileSystem: pathfs.NewDefaultFileSystem(), EtcdEndpoint: flag.Arg(0), EtcdRootPath: rootPath }
  etcdFs.NewEtcdClient().SetDir(rootPath,0)
  nfs := pathfs.NewPathNodeFs(&etcdFs, nil)
  server, _, err := nodefs.MountRoot(dir, nfs.Root(), nil)
  if err != nil {
    log.Fatalf("Mount fail: %v\n", err)
  }

  var wg sync.WaitGroup

  wg.Add(1)
  go func(){
    wg.Done()
    server.Serve()
  }()

  time.Sleep(50000)

  wg.Wait()

  exec.Command("/bin/cp","-Rv",srcRoot,dir).CombinedOutput()

  server.Unmount()
}
