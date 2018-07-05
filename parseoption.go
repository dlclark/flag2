package flag2

type ParseOption func(f *FlagSet)

func JSONVia(flagName string) ParseOption {

}

func EnvPrefix(prefix string) ParseOption {

}
