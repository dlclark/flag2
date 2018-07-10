# flag2
Package for handling flags, config files, and environment vars

The API was suggested by Peter Bourgan in his [Go for Industrial Programming](https://peter.bourgon.org/go-for-industrial-programming/#program-configuration) talk.

## Base assumptions

- Backward compatible with `flag` package.
- Existing code that uses the `flag` package that is ported to `flag2` should behave identically until you add options purposefully changing things.
    - Only look at config and EnvVar if the extra flag options are given to look there (i.e. only look in the config file if a `flag2.File(name)` option is given, only look at EnvVars if a `flag2.Env(name)` is given).  
    - You can use the `flag2.UseDefaultNamesInEnvVars()` or `flag2.UseDefaultNamesInConfigFile()` parse options to check EnvVar and Config file for ALL flags using their given names.


## Code examples

```go
var (
    foo = flag2.String("foo", "x", "foo val")
    bar = flag2.String("bar", "y", "bar val", flag2.File("bar"))
    baz = flag2.String("baz", "z", "baz val", flag2.File("baz"), flag2.Env("BAZ"))
    cfg = flag2.String("cfg", "", "JSON config file")
)
flag2.Parse(flag2.JSONFileVia("cfg"), flag2.AddEnvPrefix("MYAPP_"))
```

Or as a flag set using the defaults when parsing:

```go
var fs flag2.FlagSet
var (
    foo = fs.String("foo", "x", "foo val")
    bar = fs.String("bar", "y", "bar val")
    baz = fs.String("baz", "z", "baz val")
    cfg = fs.String("cfg", "", "JSON config file")
)
fs.Parse(os.Args, 
        UseDefaultNamesInEnvVars(),
        UseDefaultNamesInConfigFile(),
        flag2.JSONFileVia("cfg"), 
        flag2.AddEnvPrefix("MYAPP_"))
```