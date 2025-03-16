package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ViperReader struct {
	instance *viper.Viper
}

func NewViperReader(viper *viper.Viper) *ViperReader {
	return &ViperReader{
		instance: viper,
	}
}

func (vr *ViperReader) SafeGetBool(key string, current bool) bool {
	if vr.instance.InConfig(key) {
		return vr.instance.GetBool(key)
	}

	return current
}

func (vr *ViperReader) SafeGetString(key string, current string) string {
	if vr.instance.InConfig(key) {
		return vr.instance.GetString(key)
	}

	return current
}

func (vr *ViperReader) SafeGetInt(key string, current int) int {
	if vr.instance.InConfig(key) {
		return vr.instance.GetInt(key)
	}

	return current
}

type DualParamReader struct {
	viper *viper.Viper
	flags *pflag.FlagSet
}

func NewDualReader(viper *viper.Viper, flags *pflag.FlagSet) *DualParamReader {
	return &DualParamReader{
		viper: viper,
		flags: flags,
	}
}

func (vr *DualParamReader) SafeGetBool(key string, current bool) bool {
	if vr.flags.Changed(key) {
		res, _ := vr.flags.GetBool(key)
		return res
	}
	if vr.viper.InConfig(key) {
		return vr.viper.GetBool(key)
	}

	return current
}

func (vr *DualParamReader) SafeGetString(key string, current string) string {
	if vr.flags.Changed(key) {
		res, _ := vr.flags.GetString(key)
		return res
	}
	if vr.viper.InConfig(key) {
		return vr.viper.GetString(key)
	}

	return current
}

func (vr *DualParamReader) SafeGetInt(key string, current int) int {
	if vr.flags.Changed(key) {
		res, _ := vr.flags.GetInt(key)
		return res
	}
	if vr.viper.InConfig(key) {
		return vr.viper.GetInt(key)
	}

	return current
}
