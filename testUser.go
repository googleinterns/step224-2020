package main

import (
	"fmt"
	cmd "github.com/googleinterns/step224-2020/cmd"
	// "github.com/google/cloudprober/probes"
	"github.com/googleinterns/step224-2020/cloudprober/myprobe"
)

func main() {
//	client := cmd.GetClient()
//	probes.RegisterProbeType(int(myprobe.E_RedisProbe.Field),
//		func() probes.Probe { return &myprobe.Probe{} })

	probe := "cmd/myproberRedis.cfg"
	probe2 := "cmd/myproberRedis2.cfg"
	fmt.Println("Probe path")

//	cmd.AddProbe(probe, client)
	cmd.AddProbe(probe)

	cmd.RegisterAndAddProbe(probe2, &myprobe.Probe{})

	// probes.RegisterProbeType(probeType,
              // func() probes.Probe { return &myprobe.Probe{} })
	// cmd.AddProbe(probe2)
	fmt.Println("Probe added!")
}
