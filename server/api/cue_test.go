package api

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/gocode/gocodec"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	var (
		module = "github.com/verifa/coastline/server/api"
		dir    = "."
	)
	c := cuecontext.New()
	buildInstances := load.Instances([]string{module}, &load.Config{
		Dir:        dir,
		ModuleRoot: dir,
		Module:     module,
	})

	value := c.BuildInstance(buildInstances[0])
	require.NoError(t, value.Err())
	it, err := value.Value().Fields(cue.Definitions(true))
	require.NoError(t, err)

	for it.Next() {
		fmt.Println("value: ", it.Value())
		fmt.Println("label: ", it.Selector())

		it2, err := it.Value().Fields(cue.Definitions(true))
		require.NoError(t, err)

		for it2.Next() {
			if it2.Selector().IsDefinition() {
				fmt.Println("value: ", it2.Value())
				fmt.Println("label: ", it2.Selector())

			}
		}
	}
}

func TestValidate(t *testing.T) {
	var (
		module = "github.com/verifa/coastline/server/api"
		dir    = "."
	)
	// c := cuecontext.New()
	r := &cue.Runtime{}
	buildInstances := load.Instances([]string{module}, &load.Config{
		Dir:        dir,
		ModuleRoot: dir,
		Module:     module,
	})

	inst, err := r.Build(buildInstances[0])
	require.NoError(t, err)
	// value := c.BuildInstance(buildInstances[0])
	require.NoError(t, inst.Value().Err())

	// var i interface{}
	i := map[string]interface{}{
		"service": "abc",
	}

	codec := gocodec.New(r, nil)
	{
		err := codec.Validate(inst.Value(), i)
		require.NoError(t, err)
	}

}
