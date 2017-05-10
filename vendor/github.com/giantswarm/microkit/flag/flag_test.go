package flag

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Test_Init(t *testing.T) {
	e := ""
	f := testFlag{}
	s := ""

	// Make sure the uninitialized flag structure does not have any values set.
	{
		e = ""
		s = f.Config.Dirs
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = ""
		s = f.Config.Files
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = ""
		s = f.Server.Listen.Address
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = ""
		s = f.Server.TLS.CaFile
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = ""
		s = f.Server.TLS.CrtFile
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = ""
		s = f.Server.TLS.KeyFile
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = ""
		s = f.Foo
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
	}

	Init(&f)

	// Make sure the initialized flag structure does have the proper values set.
	{
		e = "config.dirs"
		s = f.Config.Dirs
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = "config.files"
		s = f.Config.Files
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = "server.listen.address"
		s = f.Server.Listen.Address
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = "server.tls.cafile"
		s = f.Server.TLS.CaFile
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = "server.tls.crtfile"
		s = f.Server.TLS.CrtFile
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = "server.tls.keyfile"
		s = f.Server.TLS.KeyFile
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = "foo"
		s = f.Foo
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
	}
}

func Test_Parse(t *testing.T) {
	f := testFlag{}
	Init(&f)

	v := viper.New()

	fs := pflag.NewFlagSet("test-flag-set", pflag.ContinueOnError)
	expectedAddress := "http://127.0.0.1:8000"
	fs.String(f.Server.Listen.Address, expectedAddress, "Test help usage.")
	expectedFoo := 74
	fs.Int(f.Foo, expectedFoo, "Test help usage.")

	Parse(v, fs)

	address := v.GetString(f.Server.Listen.Address)
	if address != expectedAddress {
		t.Fatal("expected", expectedAddress, "got", address)
	}

	foo := v.GetInt(f.Foo)
	if foo != expectedFoo {
		t.Fatal("expected", expectedFoo, "got", foo)
	}
}

type testFlag struct {
	Config testConfig
	Server testServer
	Foo    string
}

type testConfig struct {
	Dirs  string
	Files string
}

type testServer struct {
	Listen testListen
	TLS    testTLS
}

type testListen struct {
	Address string
}

type testTLS struct {
	CaFile  string
	CrtFile string
	KeyFile string
}
