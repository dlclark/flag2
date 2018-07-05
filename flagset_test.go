package flag2

import (
	"os"
	"testing"
)

func TestBasicSmoke(t *testing.T) {
	var fs FlagSet
	var (
		def = fs.Int("def", 100, "default")
		foo = fs.String("foo", "x", "foo val")
		bar = fs.String("bar", "y", "bar val", JSON("bar"))
		baz = fs.String("baz", "z", "baz val", JSON("baz"), Env("BAZ"))
		cfg = fs.String("cfg", "", "JSON config file")
	)
	// set env, write config
	os.Setenv("MYAPP_BAZ", "bazEnvValue")

	err := fs.Parse([]string{"-foo", "fooValue", "-bar", "barCfgValue", "-cfg", "cfgValue"}, JSONVia("cfg"), EnvPrefix("MYAPP_"))
	if err != nil {
		t.Fatalf("Unexpected error during parse: %v", err)
	}

	if want, got := 100, *def; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
	if want, got := "fooValue", *foo; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
	if want, got := "barCfgValue", *bar; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
	if want, got := "bazEnvValue", *baz; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
	if want, got := "cfgValue", *cfg; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
}
