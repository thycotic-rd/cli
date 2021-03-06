package cli

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/posener/complete"
)

// PredictorWrapper merges a flag with its predictor
type PredictorWrapper struct {
	complete.Predictor
	Flag           *FlagWrapper
	PredictNothing bool
}

// PredictorWrappers maps a flag name to its predictor wrapper
type PredictorWrappers map[string]PredictorWrapper

// FlagValue is a generalized storage for a flag's value
type FlagValue struct {
	Val          string
	FlagType     string
	Name         string
	DefaultValue string
	Shorthand    string
}

type FlagWrapper struct {
	Val          *FlagValue
	FlagType     string
	Name         string
	FriendlyName string
	Shorthand    string
	Usage        string
	Global       bool
	Hidden       bool
}

func (f FlagValue) IsBool() bool {
	return f.FlagType == "bool"
}

func (f *FlagValue) Set(v string) error {
	if f.FlagType == "" || f.FlagType == "string" {
		if v != "" && len(v) > 1 && strings.HasPrefix(v, "@") {
			f.FlagType = "file"
			fname := v[1:]
			if b, err := ioutil.ReadFile(fname); err != nil {
				return err
			} else {
				f.Val = string(b)
				return nil
			}
		}
	}
	f.Val = v
	return nil
}

func (f *FlagValue) Type() string {
	return f.FlagType
}

func (f *FlagValue) String() string {
	return f.Val
}

func (w PredictorWrappers) Merge(wrappers PredictorWrappers, errorOnConflict bool) (PredictorWrappers, error) {
	if errorOnConflict {
		for k, v := range wrappers {
			if v, ok := w[k]; ok && v.Predictor != nil {
				return nil, fmt.Errorf("duplicate predictor: '%s'", k)
			}
			w[k] = v
		}
	} else {
		for k, v := range wrappers {
			w[k] = v
		}
	}
	return w, nil
}
