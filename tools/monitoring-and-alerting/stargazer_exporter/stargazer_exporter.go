package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thilinapiy/stargazer_exporter/stargazer"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	_ prometheus.Collector = &collector{}

	blockAddress = kingpin.Flag("block-address", "Hash address of the block that needs to monitor").Required().String()
	bindPort     = kingpin.Flag("bind-port", "Port which listens for promethius to scrape").Default(":9119").String()
)

type collector struct {
	MissedBlocksTotal *prometheus.Desc
	stats             func() ([]stargazer.MissesBlock, error)
}

func newCollector(stats func() ([]stargazer.MissesBlock, error)) prometheus.Collector {
	return &collector{
		MissedBlocksTotal: prometheus.NewDesc(
			"stargazer_missed_blocks_total",
			"Description of Stargazer",
			[]string{"FirstBlock", "LastBlock"},
			nil,
		),

		stats: stats,
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.MissedBlocksTotal,
	}

	for _, d := range ds {
		ch <- d
	}
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	stats, err := c.stats()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.MissedBlocksTotal, err)
		return
	}

	for _, s := range stats {
		FirstBlock := s.StartHeight
		LastBlock := s.EndHeight
		count, _ := strconv.Atoi(s.Count)

		ch <- prometheus.MustNewConstMetric(
			c.MissedBlocksTotal,
			prometheus.CounterValue,
			float64(count),
			FirstBlock,
			LastBlock,
		)
	}
}

func main() {

	kingpin.Parse()

	stats := func() ([]stargazer.MissesBlock, error) {
		return stargazer.GetMissedGroups(*blockAddress)
	}

	c := newCollector(stats)
	prometheus.MustRegister(c)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	log.Printf("starting exporter on %q", *bindPort)
	if err := http.ListenAndServe(*bindPort, mux); err != nil {
		log.Fatalf("cannot start exporter: %s", err)
	}
}
