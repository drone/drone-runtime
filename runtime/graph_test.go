package runtime

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
)

func TestIsSerial(t *testing.T) {
	spec := &engine.Spec{
		Steps: []*engine.Step{
			{Metadata: engine.Metadata{Name: "build"}},
			{Metadata: engine.Metadata{Name: "test"}},
		},
	}
	if isSerial(spec) == false {
		t.Errorf("Expect is serial true")
	}

	spec.Steps[1].DependsOn = []string{"build"}
	if isSerial(spec) == true {
		t.Errorf("Expect is serial false")
	}
}
