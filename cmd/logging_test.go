package cmd_test

import (
	"io/ioutil"
	"launchpad.net/gnuflag"
	. "launchpad.net/gocheck"
	"launchpad.net/juju/go/cmd"
	"launchpad.net/juju/go/log"
	"path/filepath"
)

type LogSuite struct{}

var _ = Suite(&LogSuite{})

func (s *LogSuite) TestInitFlagSet(c *C) {
	l := &cmd.Log{}
	f := gnuflag.NewFlagSet("", gnuflag.ContinueOnError)
	l.InitFlagSet(f)

	err := f.Parse(false, []string{})
	c.Assert(err, IsNil)
	c.Assert(l.Path, Equals, "")
	c.Assert(l.Verbose, Equals, false)
	c.Assert(l.Debug, Equals, false)

	err = f.Parse(false, []string{"--log-file", "foo", "--verbose", "--debug"})
	c.Assert(err, IsNil)
	c.Assert(l.Path, Equals, "foo")
	c.Assert(l.Verbose, Equals, true)
	c.Assert(l.Debug, Equals, true)
}

func (s *LogSuite) TestStart(c *C) {
	defer saveLog()()
	for _, t := range []struct {
		path    string
		verbose bool
		debug   bool
		target  Checker
	}{
		{"", true, true, NotNil},
		{"", true, false, NotNil},
		{"", false, true, NotNil},
		{"", false, false, IsNil},
		{"foo", true, true, NotNil},
		{"foo", true, false, NotNil},
		{"foo", false, true, NotNil},
		{"foo", false, false, NotNil},
	} {
		l := &cmd.Log{t.path, t.verbose, t.debug}
		ctx := dummyContext(c)
		err := l.Start(ctx)
		c.Assert(err, IsNil)
		c.Assert(log.Target, t.target)
		c.Assert(log.Debug, Equals, t.debug)
	}
}

func (s *LogSuite) TestStderrLog(c *C) {
	defer saveLog()()
	l := &cmd.Log{Verbose: true}
	ctx := dummyContext(c)
	err := l.Start(ctx)
	c.Assert(err, IsNil)
	log.Printf("hello")
	c.Assert(str(ctx.Stderr), Matches, `.* JUJU hello\n`)
}

func (s *LogSuite) TestRelPathLog(c *C) {
	defer saveLog()()
	l := &cmd.Log{Path: "foo.log"}
	ctx := dummyContext(c)
	err := l.Start(ctx)
	c.Assert(err, IsNil)
	log.Printf("hello")
	c.Assert(str(ctx.Stderr), Equals, "")
	content, err := ioutil.ReadFile(filepath.Join(ctx.Dir, "foo.log"))
	c.Assert(err, IsNil)
	c.Assert(string(content), Matches, `.* JUJU hello\n`)
}

func (s *LogSuite) TestAbsPathLog(c *C) {
	defer saveLog()()
	path := filepath.Join(c.MkDir(), "foo.log")
	l := &cmd.Log{Path: path}
	ctx := dummyContext(c)
	err := l.Start(ctx)
	c.Assert(err, IsNil)
	log.Printf("hello")
	c.Assert(str(ctx.Stderr), Equals, "")
	content, err := ioutil.ReadFile(path)
	c.Assert(err, IsNil)
	c.Assert(string(content), Matches, `.* JUJU hello\n`)
}
