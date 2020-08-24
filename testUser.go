package main

import (
	"fmt"
	cmd "github.com/googleinterns/step224-2020/cmd"
	"github.com/google/cloudprober/probes"
	"github.com/googleinterns/step224-2020/cloudprober/myprobe"
)

func main() {
	probe := "cmd/myprober.cfg"
	probe2 := "cmd/myproberRedis2.cfg"
	fmt.Println("Probe path")

	fmt.Println("Extension: ", int(myprobe.E_RedisProbe.Field))

	cmd.RegisterAndAddProbeFromFilepath(probe, &myprobe.Probe{})

	probes.RegisterProbeType(200, func() probes.Probe {
		return &myprobe.Probe{}
	})

	cmd.RegisterAndAddProbeFromFilepath(probe2, &myprobe.Probe{})

	fmt.Println("Probe added!")
}
