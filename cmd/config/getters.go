package config

import "github.com/spf13/pflag"

type FlagReader struct {
	flagSet *pflag.FlagSet
}

func NewFlagReader(flagSet *pflag.FlagSet) *FlagReader {
	return &FlagReader{flagSet: flagSet}
}

func (fr *FlagReader) SafeGetBool(key string, current bool) bool {
	if fr.flagSet.Changed(key) {
		value, err := fr.flagSet.GetBool(key)
		if err == nil {
			return value
		}
	}

	return current
}

func (fr *FlagReader) SafeGetString(key string, current bool) bool {
	if fr.flagSet.Changed(key) {
		value, err := fr.flagSet.GetBool(key)
		if err == nil {
			return value
		}
	}

	return current
}
