package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thilinapiy/stargazer_exporter/stargazer"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	_          prometheus.Collector = &collector{}
	sendClient                      = &http.Client{Timeout: 10 * time.Second}

	oldStats     []stargazer.MissesBlock
	blockAddress = kingpin.Flag("block-address", "Hash address of the block that needs to monitor").Required().String()
	bindPort     = kingpin.Flag("bind-port", "Port which listens for prometheus to scrape").Default(":9119").String()
	chatID       = kingpin.Flag("chat-id", "Telegram chat group id").Required().String()
	token        = kingpin.Flag("bot-token", "Telegram bot secret token").Required().String()
)

type collector struct {
	MissedBlocksTotal *prometheus.Desc
	NewMissesBlocks   *prometheus.Desc
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
		NewMissesBlocks: prometheus.NewDesc(
			"new_missed_blocks",
			"If there are new missed blocks this will return true, compaired to previous scrape",
			[]string{"FirstBlock", "LastBlock"},
			nil,
		),

		stats: stats,
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.MissedBlocksTotal,
		c.NewMissesBlocks,
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

	if oldStats != nil {
		for _, s := range getChanges(oldStats, stats) {
			FirstBlock := s.StartHeight
			LastBlock := s.EndHeight
			count, _ := strconv.Atoi(s.Count)

			ch <- prometheus.MustNewConstMetric(
				c.NewMissesBlocks,
				prometheus.CounterValue,
				float64(count),
				FirstBlock,
				LastBlock,
			)
		}
	}

	oldStats = stats
}

func getChanges(oldStats, newStats []stargazer.MissesBlock) []stargazer.MissesBlock {
	k := 0
	for ; oldStats[0].StartHeight != newStats[k].StartHeight; k++ {

	}
	if k != 0 {
		sendMessage(newStats[:k], false)
	}
	return newStats[:k]
}

func sendMessage(msg []stargazer.MissesBlock, alert bool) bool {

	url := "https://api.telegram.org/bot" + *token + "/sendMessage"

	message := ""
	critical := 0
	msgBuff := "The Chainflow Validator missed\n"
	for _, block := range msg {
		msgBuff += fmt.Sprintf("%s blocks between block numbers %s-%s\n", block.Count, block.StartHeight, block.EndHeight)
		intCount, _ := strconv.Atoi(block.Count)
		critical += intCount
	}
	msgBuff += fmt.Sprintf("at %s", time.Now().UTC())

	if critical > 24 {
		message = "⚠️" + " Critical Alert!\n" + msgBuff
	}

	for _, chat := range strings.Split(*chatID, ",") {
		body := map[string]interface{}{
			"chat_id":              chat,
			"text":                 message,
			"disable_notification": alert,
		}

		bytesRepresentation, err := json.Marshal(body)
		if err != nil {
			log.Fatalln(err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytesRepresentation))
		req.Header.Set("Content-Type", "application/json")

		resp, er := sendClient.Do(req)
		if er != nil {
			log.Fatal("Error in request send")
			return false
		}

		if err != nil {
			log.Fatal("Error in request create")
			return false
		}
		defer resp.Body.Close()
	}
	return true
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
