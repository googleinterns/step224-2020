import (
	"context"
	"fmt"
	"sync"
	"io"
	"os"
	"time"
	"cloud.google.com/go/storage"
	"github.com/google/cloudprober/logger"
	"github.com/google/cloudprober/metrics"
	"github.com/google/cloudprober/probes/options"
	"github.com/google/cloudprober/targets/endpoint"
	"github.com/googleinterns/step224-2020/blob/evan/hermesMain/cmd"
)

// Probe holds aggregate information about all probe runs, per-target.
type HermesProbe struct {
	name    string
	c       *ProbeConf
	targets []string
	opts    *options.Options

	res map[string]*metrics.EventMetrics  // Results by target
	l   *logger.Logger
}

// Init initializes the probe with the given params.
func (p *HermesProbe) Init(name string, opts *options.Options) error {
	p.c, ok := opts.ProbeConf.(*ProbeConf)
	if !ok {
		return fmt.Errorf("Invalid Hermes probe configuration")
	}
	p.name = name
	p.opts = opts
	p.l = opts.Logger
	p.res = make(map[string]*metrics.EventMetrics)

	return nil
}

// Start starts and runs the probe indefinitely.
func (p *HermesProbe) Start(ctx context.Context, dataChan chan *metrics.EventMetrics) {
	probeTicker := time.NewTicker(p.opts.Interval)

	for {
		select {
		case <-ctx.Done():
			probeTicker.Stop()
			return
		case <-probeTicker.C:
			// On probe tick, write data to the channel and run probe.
			for _, em := range p.res {
				dataChan <- em
			}
			p.targets = endpoint.NamesFromEndpoints(p.opts.Targets.ListEndpoints())
			p.initProbeMetrics()
			probeCtx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(p.opts.Timeout))
			p.runProbe(probeCtx)
			cancelFunc()
		}
	}
}

// initProbeMetrics initializes missing probe metrics.
func (p *HermesProbe) initProbeMetrics() {
	for _, target := range p.targets {
		if p.res[target] != nil {
			continue
		}
		var latVal metrics.Value
		if p.opts.LatencyDist != nil {
			latVal = p.opts.LatencyDist.Clone()
		} else {
			latVal = metrics.NewFloat(0)
		}
		p.res[target] = metrics.NewEventMetrics(time.Now()).
			AddMetric("total", metrics.NewInt(0)).
			AddMetric("success", metrics.NewInt(0)).
			AddMetric("latency", latVal).
			AddLabel("probe", p.name).
			AddLabel("dst", target)
	}
}

// runProbeForTarget runs probe for a single target.
func (p *HermesProbe) runProbeForTarget(ctx context.Context, target string) error {
	bucket := p.c.GetBucket()
 	object := p.c.GetObject()
	client, err := storage.NewClient(ctx)
          if err != nil {
                  return fmt.Errorf("storage.NewClient: %v", err)
          }
          defer client.Close()
 
         ctx, cancel := context.WithTimeout(ctx, p.opts.Timeout)
		 defer cancel()

		 f, err := os.Open("notes.txt")
		 if err != nil {
				 return fmt.Errorf("os.Open: %v", err)
		 }
		 defer f.Close()
 
		 // Upload an object with storage.Writer.
		 wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
		 if _, err = io.Copy(wc, f); err != nil {
				 return fmt.Errorf("io.Copy: %v", err)
		 }
		 if err := wc.Close(); err != nil {
				 return fmt.Errorf("Writer.Close: %v", err)
		 }
		 fmt.Fprintf(w, "Blob %v uploaded.\n", object)
		 return nil

}

// runProbe runs probe for all targets and update EventMetrics.
func (p *HermesProbe) runProbe(ctx context.Context) {
	p.targets = endpoint.NamesFromEndpoints(p.opts.Targets.ListEndpoints())

	var wg sync.WaitGroup
	for _, target := range p.targets {
		wg.Add(1)

		go func(target string, em *metrics.EventMetrics) {
			defer wg.Done()
			em.Metric("total").AddInt64(1)
			start := time.Now()
			err := p.runProbeForTarget(ctx, target) // run probe just for a single target
			if err != nil {
				p.l.Errorf(err.Error())
				return
			}
			em.Metric("success").AddInt64(1)
			em.Metric("latency").AddFloat64(time.Now().Sub(start).Seconds() / p.opts.LatencyUnit.Seconds())
		}(target, p.res[target])

	}

	wg.Wait()
}
