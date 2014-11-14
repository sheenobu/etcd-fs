package etcdfs

import(
  . "github.com/franela/goblin"
  "testing"
  etcdm "github.com/coreos/go-etcd/etcd"
//  "fmt"
  "os"
  "io/ioutil"
)

func TestNodeFs(t *testing.T) {
  g := Goblin(t)

  roots := []string{"/x","","/z/s","/_d/x/d"}

  for i := range roots {

    root := roots[i]

    g.Describe("File_" + root, func() {
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

      g.Describe("Open", func() {
        g.It("Should be supported", func() {
          if _, e := etcd.Set(root + "/test/foo", "bar", 0); e != nil {
            g.Fail(e)
          }

          file, err := os.Open(fs.Path() + "/test/foo")

          if err != nil {
            g.Fail(err)
          }

          file.Close()
        })
      })
      g.Describe("Create", func() {
        g.It("Should be supported", func() {
          file, err := os.Create(fs.Path() + "/test/bar")

          if err != nil {
            g.Fail(err)
          }
          file.Close()

          if _, er := etcd.Get(root + "/test/bar", false, false); er != nil {
            g.Fail(er)
          }
        })
      })
      g.Describe("Delete", func() {
        g.It("Should be supported", func() {
          etcd.Set(root + "/test/barfoo", "lala", 0)

          err := os.Remove(fs.Path() + "/test/barfoo")

          if err != nil {
            g.Fail(err)
          }

          if _, er := etcd.Get(root + "/test/barfoo", false, false); er == nil {
            g.Fail("The key [" + root + "/test/barfoo] should not exist")
          }
        })
      })
      g.Describe("Read", func() {
        g.It("Should be supported", func() {
          etcd.Set(root + "/test/bar", "foo", 0)

          data, err := ioutil.ReadFile(fs.Path() + "/test/bar")

          if err != nil {
            g.Fail(err)
          }

          g.Assert(string(data)).Equal("foo")
        })
      })
      g.Describe("Write", func() {
        g.It("Should be supported", func() {
          if err := ioutil.WriteFile(fs.Path() + "/test/foobar", []byte("hello world"), 0666); err != nil {
            g.Fail(err)
          }

          res, err := etcd.Get(root + "/test/foobar", false, false)

          if err != nil {
            g.Fail(err)
          }

          g.Assert(res.Node.Value).Equal("hello world")
        })
      })
    })
  }
}
