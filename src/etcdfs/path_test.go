package etcdfs

import(
  . "github.com/franela/goblin"
  "testing"
  etcdm "github.com/coreos/go-etcd/etcd"
  "fmt"
  "os"
)

func TestPathFs(t *testing.T) {
  g := Goblin(t)

  roots := []string{"/x","","/z/s","/_d/x/d"}

  for i := range roots {

    root := roots[i]

    g.Describe("Path_" + root, func() {
      var etcd *etcdm.Client
      var fs testEtcdFsMount

      g.Before(func() {
        etcd = etcdm.NewClient([]string{testEtcdEndpoint})
      })

      g.BeforeEach(func() {
        etcd.RawDelete(root + "/test", true, true)
        etcd.SetDir(root + "/test", 0)
        fs = NewTestEtcdFsMount(root)
      })

      g.AfterEach(func() {
        fs.Unmount()
      })

      g.Describe("ls", func() {
        g.It("Should be supported", func() {
          etcd.Set(root + "/test/a", "a", 0)
          etcd.SetDir(root + "/test/b", 0)

          f, err1 := os.Open(fs.Path() + "/test")

          if err1 != nil {
            g.Fail(err1)
          }
          defer f.Close()

          files, err2 := f.Readdir(0)

          if err2 != nil {
            g.Fail(err2)
          }

          g.Assert(len(files) == 2).IsTrue()

          file1 := files[0]
          file2 := files[1]

          switch file1.Name() {
            case "a":
              g.Assert(file1.IsDir()).IsFalse()
              g.Assert(uint64(file1.Size())).Equal(uint64(1))
              g.Assert(file1.Mode().String()).Equal("-rw-rw-rw-")
            case "b":
              g.Assert(file1.IsDir()).IsTrue()
              g.Assert(file1.Mode().String()).Equal("drw-rw-rw-")
            default:
              g.Fail(fmt.Sprintf("Didn't expect file [%s]", file1.Name()))
          }
          switch file2.Name() {
            case "a":
              g.Assert(file2.IsDir()).IsFalse()
              g.Assert(uint64(file2.Size())).Equal(uint64(1))
              g.Assert(file2.Mode().String()).Equal("-rw-rw-rw-")
            case "b":
              g.Assert(file2.IsDir()).IsTrue()
              g.Assert(file2.Mode().String()).Equal("drw-rw-rw-")
            default:
              g.Fail(fmt.Sprintf("Didn't expect file [%s]", file2.Name()))
          }
        })
      })
      g.Describe("mkdir", func() {
        g.It("Should be supported", func() {
          if e := os.Mkdir(fs.Path() + "/test/foo", os.ModeDir | 0666); e != nil {
            g.Fail(e)
          }
          res, err := etcd.Get(root + "/test/foo", false, false)
          if err != nil {
            g.Fail(err)
          }
          g.Assert(res.Node.Dir).IsTrue()
        });
        g.It("Should support creating with parents", func() {
          if e := os.MkdirAll(fs.Path() + "/test/a/b/c/foo", os.ModeDir | 0666); e != nil {
            g.Fail(e)
          }
          res, err := etcd.Get(root + "/test/a/b", false, false)
          if err != nil {
            g.Fail(err)
          }
          g.Assert(res.Node.Dir).IsTrue()
        });
      })
      g.Describe("rmdir", func() {
        g.It("Should be supported", func() {
          etcd.CreateDir(root + "/test/foo", 0)
          if e := os.Remove(fs.Path() + "/test/foo"); e != nil {
            g.Fail(e)
          }
          _, err := etcd.Get(root + "/test/foo", false, false)
          if err == nil {
            g.Fail(root + "/test/foo should not exist in etcd.")
          }
        });
        g.It("Should support removing with children", func() {
          etcd.CreateDir(root + "/test/foo/bar", 0)
          if e := os.RemoveAll(fs.Path() + "/test/foo"); e != nil {
            g.Fail(e)
          }
          _, err := etcd.Get(root + "/test/foo/bar", false, false)
          if err == nil {
            g.Fail(root + "/test/foo should not exist in etcd.")
          }
        });
      })
    })
  }
}
