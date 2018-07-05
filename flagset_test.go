package flag2

import (
	"os"
	"testing"
)

func BasicSmoke(t *testing.T) {
	var fs FlagSet
	var (
		foo = fs.String("foo", "x", "foo val")
		bar = fs.String("bar", "y", "bar val", JSON("bar"))
		baz = fs.String("baz", "z", "baz val", JSON("baz"), Env("BAZ"))
		cfg = fs.String("cfg", "", "JSON config file")
	)
	fs.Parse(os.Args, JSONVia("cfg"), EnvPrefix("MYAPP_"))

}
