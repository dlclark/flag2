package flag2

type FlagOption func(*Flag) error

func JSON(name string) FlagOption {
	return func(f *Flag) error {
		f.NameInConfigFile = name
		return nil
	}
}

func Env(name string) FlagOption {
	return func(f *Flag) error {
		f.NameInEnv = name
		return nil
	}
}
