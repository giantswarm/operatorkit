package flag

import (
	"encoding/json"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	microerror "github.com/giantswarm/microkit/error"
)

func Init(f interface{}) {
	b, err := json.Marshal(f)
	if err != nil {
		panic(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}

	for k, v := range m {
		m[k] = toValue([]string{strings.ToLower(k)}, k, v)
	}
	b, err = json.Marshal(m)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, f)
	if err != nil {
		panic(err)
	}
}

func Parse(v *viper.Viper, fs *pflag.FlagSet) {
	v.BindPFlags(fs)
}

func Merge(v *viper.Viper, fs *pflag.FlagSet, dirs, files []string) error {
	// We support multiple config files. Viper cannot do that on its own. So we
	// configure a new viper for each config file that we are interested in and
	// merge the found configurations into the viper given by the client.
	for _, f := range files {
		newViper := viper.New()
		newViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		newViper.AutomaticEnv()
		for _, configDir := range dirs {
			newViper.AddConfigPath(configDir)
		}
		newViper.SetConfigName(f)

		err := newViper.ReadInConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// In case there is no config file given we simply go ahead to check
				// the other ones. If we do not find any configuration using config
				// files, we go ahead to check the process environment.
			} else {
				return microerror.MaskAny(err)
			}
		}

		fs.VisitAll(func(f *pflag.Flag) {
			if newViper.IsSet(f.Name) {
				v.Set(f.Name, newViper.Get(f.Name))
			}
		})
	}

	return nil
}

func toValue(path []string, key string, val interface{}) interface{} {
	m, ok := val.(map[string]interface{})
	if ok {
		for k, v := range m {
			m[k] = toValue(append([]string{strings.ToLower(k)}, path...), k, v)
		}

		return m
	}

	res := strings.Join(reverse(path), ".")
	return res
}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}
