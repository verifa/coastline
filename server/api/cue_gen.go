// Code generated by gocode.Generate; DO NOT EDIT.

package api

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/gocode/gocodec"
	_ "cuelang.org/go/pkg"
)

var cuegenvalProject = cuegenMake("Project", &Project{})

// Validate validates x.
func (x *Project) Validate() error {
	return cuegenCodec.Validate(cuegenvalProject, x)
}

var cuegenCodec, cuegenInstance = func() (*gocodec.Codec, *cue.Instance) {
	var r *cue.Runtime
	r = &cue.Runtime{}
	instances, err := r.Unmarshal(cuegenInstanceData)
	if err != nil {
		panic(err)
	}
	if len(instances) != 1 {
		panic("expected encoding of exactly one instance")
	}
	return gocodec.New(r, nil), instances[0]
}()

// cuegenMake is called in the init phase to initialize CUE values for
// validation functions.
func cuegenMake(name string, x interface{}) cue.Value {
	f, err := cuegenInstance.Value().FieldByName(name, true)
	if err != nil {
		panic(fmt.Errorf("could not find type %q in instance", name))
	}
	v := f.Value
	if x != nil {
		w, err := cuegenCodec.ExtractType(x)
		if err != nil {
			panic(err)
		}
		v = v.Unify(w)
	}
	return v
}

// Data size: 248 bytes.
var cuegenInstanceData = []byte("\x01\x1f\x8b\b\x00\x00\x00\x00\x00\x00\xffD\xceMK\xc3@\x10\xc6\xf1y\xd2\bu\xa9~\x03!T\x10\x11\x92x\x13\v\x1e\x04\u045b\x14\x8f\x96\x1e\xc6e\u06ae&\u0650\xdd\xf4\xa0(\xa8m\xf5[\xafD|\xb9\xfe\x98\xff\xf0\uc10f\bQ\xf8$\x847\xa2\x93\xf0\xda\x03\x06\xa6r\x9e+-\x17\xec\xb9s\xf4\x10\xdfX\xeb\x11\x11\xe21\xfb\x05\x06\x84\xadKS\x88C\xd8\x10\xd1^XG\xc0\xeed\xaa[\xc9f\xa6\xf8)7\x84\xb0\":\f\xef=\xa0\xff\xef+B\x84\xf8\x9aK\xe9\x1e\xc5\u07e8\x88(\xac\xbb!\x00\x0e\xe6\xc6/\u06bbL\xdb2_Jcf\x9ck\xcb\xce\x17\xa6\x92\xdcI\xb3\x94&\xe7\xda\x00\xe8\xbbZt\xa6[\xc1U\xcd\xfa\x81\xe7\x92pm\x94\xda\x1f7\xf6^\xb4\x1f%Oj\xbb\xe2RF\xc9\xd9\xcbpr\x9e\xder\xfax\x9c\x9e\xa6\xd3\xech\xa8\x9e\xd5\xdf\xd9o\xa0\x88\xbe\x02\x00\x00\xff\xff\xec\xeb\x12\xbe\x14\x01\x00\x00")
