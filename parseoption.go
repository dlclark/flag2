package flag2

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// ErrNoValue is the error returned if the config file doesn't contain a value
// for the requested flag
var ErrNoValue = errors.New("flag: help requested")

type ParseOption func(*FlagSet) error

func JSONVia(flagName string) ParseOption {
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

func EnvPrefix(prefix string) ParseOption {
	return func(f *FlagSet) error {
		f.VisitAll(func(flag *Flag) {
			if flag.NameInEnv != "" {
				flag.NameInEnv = prefix + flag.NameInEnv
			}
		})
		return nil
	}
}
