package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func (c *Config) Load(cmd *cobra.Command) error {
	var errs []error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			if val, ok := os.LookupEnv(EnvName(f.Name)); ok {
				if err := f.Value.Set(val); err != nil {
					errs = append(errs, err)
				}
			}
		}
	})
	return errors.Join(errs...)
}

const EnvPrefix = "RELAX_"

func EnvName(name string) string {
	name = strings.ToUpper(name)
	name = strings.ReplaceAll(name, "-", "_")
	return EnvPrefix + name
}
