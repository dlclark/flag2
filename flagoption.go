package flag2

// A FlagOption allows setting properties for a specific flag
type FlagOption func(*Flag) error

// File sets the flag's key name in a config file
func File(name string) FlagOption {
	return func(f *Flag) error {
		f.NameInConfigFile = name
		return nil
	}
}

// Env sets the flag's key name in the EnvVars
func Env(name string) FlagOption {
	return func(f *Flag) error {
		f.NameInEnv = name
		return nil
	}
}
