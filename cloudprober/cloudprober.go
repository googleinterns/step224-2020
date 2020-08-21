package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/golang/glog"
	"github.com/google/cloudprober"
	"github.com/googleinterns/step224-2020/cloudprober/myprobe"
	"github.com/google/cloudprober/probes"
	"github.com/google/cloudprober/web"
)

var (
	configFile = flag.String("config_file", "", "Config file")
)

func configFileToString(fileName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Failed to read the config file: %v", err)
	}
	return string(b)
}

func getConfig() string {
	if *configFile != "" {
		return configFileToString(*configFile)
	} else {
		fmt.Println("Config file not found")
		return ""
	}
}

func main() {
	flag.Parse()

	// Register stubby probe type
	probes.RegisterProbeType(int(myprobe.E_RedisProbe.Field), // TODO: CHANGE THIS
		func() probes.Probe { return &myprobe.Probe{} })

	err := cloudprober.InitFromConfig(getConfig())
	if err != nil {
		glog.Exitf("Error initializing cloudprober. Err: %v", err)
	}

	// web.Init sets up web UI for cloudprober.
	web.Init()

	cloudprober.Start(context.Background())

	// Wait forever
	select {}
}
