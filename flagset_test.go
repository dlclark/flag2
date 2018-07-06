package flag2

import (
	"io/ioutil"
	"os"
	"testing"
)

func init() {
	// set an env for all the tests

	os.Setenv("MYAPP_BAZ", "bazEnvValue")
	os.Setenv("MYAPP_INT", "100")
}

func TestBasicSmoke(t *testing.T) {
	//ARRANGE
	var fs FlagSet
	var (
		def = fs.Int("def", 100, "default")
		foo = fs.String("foo", "x", "foo val")
		bar = fs.String("bar", "y", "bar val", File("bar"))
		baz = fs.String("baz", "z", "baz val", File("baz"), Env("BAZ"))
		cfg = fs.String("cfg", "", "JSON config file")
	)
	// set env, write config
	file, err := ioutil.TempFile("", "testcfg_")
	if err != nil {
		t.Fatal(err)
	}
	var cfgFile = file.Name()
	ioutil.WriteFile(cfgFile, []byte("{ \"bar\" : \"barCfgValue\" }"), 0666)

	// ACT
	err = fs.Parse([]string{"-foo", "fooValue", "-cfg", cfgFile}, JSONFileVia("cfg"), AddEnvPrefix("MYAPP_"))

	//ASSERT
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
	if want, got := cfgFile, *cfg; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
}

func TestConfigOverrides(t *testing.T) {
	//ARRANGE
	ResetForTesting(nil)

	flg1 := Int("int", 10, "", File("int"), Env("INT"))
	flg2 := String("baz", "default", "", File("baz"), Env("BAZ"))
	String("none", "default", "")

	//ACT
	CommandLine.Parse([]string{"-none", "ignore"}, AddEnvPrefix("MYAPP_"), InMemoryConfig(map[string]string{"int": "1000"}))

	//ASSERT
	//config overrides env
	if want, got := 1000, *flg1; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
	// missing config doesn't override env
	if want, got := "bazEnvValue", *flg2; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
}

func TestArgOverrides(t *testing.T) {
	//ARRANGE
	ResetForTesting(nil)

	flg1 := Int("int", 10, "", File("int"), Env("MYAPP_INT"))
	flg2 := String("baz", "default", "", Env("baz"))
	flg3 := String("oth", "othdefault", "", File("oth"))

	//ACT
	CommandLine.Parse([]string{"-int", "10000", "-baz", "argVal", "-oth", "otherVal"},
		InMemoryConfig(map[string]string{"int": "1000", "oth": "val"}))

	//ASSERT
	// arg overrides config and env
	if want, got := 10000, *flg1; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
	// arg overrides just config
	if want, got := "argVal", *flg2; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
	// arg overrides just env
	if want, got := "otherVal", *flg3; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}

}

func TestArgDefault(t *testing.T) {
	//ARRANGE
	ResetForTesting(nil)

	flg1 := Int64("int64", 10, "", File("int64"), Env("MYAPP_INT64"))
	flg2 := String("none", "default", "")
	flg3 := Int("int", 1, "")

	//ACT
	CommandLine.Parse([]string{}, InMemoryConfig(map[string]string{"none": "1000"}))

	//ASSERT
	// value not present uses default
	if want, got := int64(10), *flg1; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
	// value present in config, but not setup to pull from it
	if want, got := "default", *flg2; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
	// value present in env, but not setup to pull from it
	if want, got := 1, *flg3; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
}

func TestDefaultNames(t *testing.T) {
	//ARRANGE
	ResetForTesting(nil)

	flg1 := Int64("int64", 10, "")
	flg2 := String("none", "default", "")
	flg3 := Int("int", 1, "")

	//ACT
	CommandLine.Parse([]string{},
		InMemoryConfig(map[string]string{"none": "1000"}),
		UseDefaultNamesInEnvVars(),
		UseDefaultNamesInConfigFile(),
		AddEnvPrefix("MYAPP_"))

	//ASSERT
	// value not present uses default
	if want, got := int64(10), *flg1; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
	// value present in config
	if want, got := "1000", *flg2; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
	// value present in env
	if want, got := 100, *flg3; want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
}
