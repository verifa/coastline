package cuehack

import (
	"fmt"
	"os"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/gocode"
	"github.com/stretchr/testify/require"
)

func TestApi(t *testing.T) {

	var (
		module = "github.com/verifa/coastline/server/cuehack"
		dir    = "."
	)
	buildInstances := load.Instances([]string{module}, &load.Config{
		Dir:        dir,
		ModuleRoot: dir,
		Module:     module,
	})
	inst := cue.Build(buildInstances)[0]
	require.NoError(t, inst.Err)

	// b, err := openapi.Gen(inst, &openapi.Config{})
	// require.NoError(t, err)

	// var out bytes.Buffer
	// err = json.Indent(&out, b, "", "   ")
	// require.NoError(t, err)

	// fmt.Println(out.String())
	it, err := inst.Value().Fields(cue.Definitions(true))
	require.NoError(t, err)

	for it.Next() {
		fmt.Println("value: ", it.Value())
		fmt.Println("label: ", it.Label())
	}

	{
		b, err := gocode.Generate(dir, inst, &gocode.Config{})
		require.NoError(t, err)

		{
			err := os.WriteFile("cue_gen.go", b, 0644)
			require.NoError(t, err)
		}

		// err = ioutil.WriteFile("cue_gen.go", b, 0644)
	}
}
