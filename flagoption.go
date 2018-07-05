package flag2

type FlagOption func(f *Flag)

func JSON(node string) FlagOption {

}

func Env(name string) FlagOption {

}
