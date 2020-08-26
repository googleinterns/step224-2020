package main

import (
	"fmt"
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober/metrics"
	"github.com/google/cloudprober/probes"

	"github.com/google/cloudprober/probes/options"
	probes_configpb "github.com/google/cloudprober/probes/proto"
	testdatapb "github.com/google/cloudprober/probes/testdata"
	targetspb "github.com/google/cloudprober/targets/proto"

	cmd "github.com/googleinterns/step224-2020/cmd"
)


// testProbe implements the probes.Probe interface, while providing
// facilities to examine the probe status for the purpose of testing.
type testProbe struct {
	intialized      bool
	runningStatusCh chan bool
}

func (p *testProbe) Init(name string, opts *options.Options) error {
	p.intialized = true
	p.runningStatusCh = make(chan bool)
	return nil
}

func (p *testProbe) Start(ctx context.Context, dataChan chan *metrics.EventMetrics) {
	p.runningStatusCh <- true

	// If context is done (used to stop a running probe before removing it),
	// change probe state to not-running.
	<-ctx.Done()
	p.runningStatusCh <- false
	close(p.runningStatusCh)
}

func generateExtensionProbeDef(name string) *probes_configpb.ProbeDef {
	probeDef := &probes_configpb.ProbeDef{
		Name: proto.String(name),
		Type: probes_configpb.ProbeDef_EXTENSION.Enum(),
		Targets: &targetspb.TargetsDef{
			Type: &targetspb.TargetsDef_DummyTargets{},
		},
	}

	proto.SetExtension(probeDef, testdatapb.E_FancyProbe, &testdatapb.FancyProbe{Name: proto.String("fancy")})
	return probeDef
}

func main() {
	// Create a probe and then register and add it to the prober
	probeName1 := "testExtension1"

	probes.RegisterProbeType(200, func() probes.Probe {
		return &testProbe{}
	})
	probe := generateExtensionProbeDef(probeName1)

	extensions := proto.RegisteredExtensions(probe)
	fmt.Println("Extensions: ", extensions)
	fmt.Println("Extension: ",  *extensions[200])
	fmt.Println("Probe config: ", probe)

	cmd.RegisterAndAddProbe(int(testdatapb.E_FancyProbe.Field), probe, &testProbe{})
	fmt.Println("Probe added!")
}
