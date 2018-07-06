package flag2

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// ErrNoValue is the error returned if the config file doesn't contain a value
// for the requested flag
var ErrNoValue = errors.New("flag: help requested")

// ParseOption is a hook that allows mixins to update a FlagSet after command line
// args are parsed, but before envvar and config file.
type ParseOption func(*FlagSet) error

// JSONFileVia will set the config file based on the value of the given flag
func JSONFileVia(flagName string) ParseOption {
	return func(f *FlagSet) error {
		if f.configFile != nil {
			return f.failf("error handling %v flag, config file already set", flagName)
		}

		flag := f.Lookup(flagName)

		f.configFile = &jsonConfigFile{
			fs:           f,
			fileNameFlag: flag,
		}

		return nil
	}
}

type jsonConfigFile struct {
	fs           *FlagSet
	fileNameFlag *Flag
	file         *os.File
	data         map[string]string
}

func (j *jsonConfigFile) FileName() string {
	return j.fileNameFlag.Value.String()
}

func (j *jsonConfigFile) ConfigValue(name string) (string, error) {
	if j.data == nil {
		return "", ErrNoValue
	}

	value, ok := j.data[name]
	if !ok {
		return "", ErrNoValue
	}

	return value, nil
}

func (j *jsonConfigFile) Open() error {
	if j.file != nil {
		return nil
	}

	// open the file
	fileName := j.FileName()

	if fileName == "" {
		// no file, no values
		return nil
	}

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	j.file = file

	bytes, err := ioutil.ReadAll(j.file)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &j.data); err != nil {
		return err
	}

	return nil
}

func (j *jsonConfigFile) Close() {
	if j.file != nil {
		j.file.Close()
	}
}

// AddEnvPrefix will prepend all flag Env names with a given prefix
func AddEnvPrefix(prefix string) ParseOption {
	return func(f *FlagSet) error {
		f.VisitAll(func(flag *Flag) {
			if flag.NameInEnv != "" {
				flag.NameInEnv = prefix + flag.NameInEnv
			}
		})
		return nil
	}
}

// UseDefaultNamesInConfigFile will allow every flag to be found in the Config file
func UseDefaultNamesInConfigFile() ParseOption {
	return func(f *FlagSet) error {
		f.VisitAll(func(flag *Flag) {
			if flag.NameInConfigFile == "" {
				flag.NameInConfigFile = flag.Name
			}
		})
		return nil
	}
}

// UseDefaultNamesInEnvVars will allow every flag to be found in the EnvVars
func UseDefaultNamesInEnvVars() ParseOption {
	return func(f *FlagSet) error {
		f.VisitAll(func(flag *Flag) {
			if flag.NameInEnv == "" {
				flag.NameInEnv = strings.ToUpper(flag.Name)
			}
		})
		return nil
	}
}

// InMemoryConfig will set the config file data to a given map of data
func InMemoryConfig(data map[string]string) ParseOption {
	return func(f *FlagSet) error {
		f.configFile = &memoryConfigFile{data: data}
		return nil
	}
}

type memoryConfigFile struct {
	data map[string]string
}

func (m *memoryConfigFile) ConfigValue(name string) (string, error) {
	val, found := m.data[name]
	if !found {
		return "", ErrNoValue
	}
	return val, nil
}

func (m *memoryConfigFile) FileName() string {
	return "in memory"
}

func (m *memoryConfigFile) Open() error {
	return nil
}

func (m *memoryConfigFile) Close() {}
