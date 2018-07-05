# flag2
Package for handling flags, config files, and environment vars

API as suggested by Peter Bourgan in his [Go for Industrial Programming](https://peter.bourgon.org/go-for-industrial-programming/#program-configuration) talk.

```go
var fs flag2.FlagSet
var (
    foo = fs.String("foo", "x", "foo val")
    bar = fs.String("bar", "y", "bar val", flag2.JSON("bar"))
    baz = fs.String("baz", "z", "baz val", flag2.JSON("baz"), flag2.Env("BAZ"))
    cfg = fs.String("cfg", "", "JSON config file")
)
fs.Parse(os.Args, flag2.JSONVia("cfg"), flag2.EnvPrefix("MYAPP_"))

```