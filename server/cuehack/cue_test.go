package cuehack

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/gocode/gocodec"
	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/server/oapi"
)

func TestAPI(t *testing.T) {
	var (
		module = "github.com/verifa/coastline/server/cuehack"
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
		module = "github.com/verifa/coastline/server/cuehack"
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

	reqs := []oapi.NewRequest{
		{
			Type:        "ArtifactoryRepoRequest",
			RequestedBy: "someone",
			Spec: map[string]interface{}{
				"repo": "hello",
			},
		},
		{
			Type:        "JenkinsServerRequest",
			RequestedBy: "someone",
			Spec: map[string]interface{}{
				"name": "hello",
			},
		},
	}

	codec := gocodec.New(r, nil)
	for _, tt := range reqs {
		t.Run(tt.Type, func(t *testing.T) {
			i := map[string]interface{}{
				"request": tt,
			}
			{
				err := codec.Validate(inst.Value(), i)
				require.NoError(t, err)
			}
		})
	}

}
