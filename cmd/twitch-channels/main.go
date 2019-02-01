package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"time"
	"strconv"

	"github.com/mlvzk/twitch_logger"
)

func main() {
	var language string
	var clientId string
	var minViewers int
	var loop bool
	flag.StringVar(&language, "language", "", "show only streams of specified locale ID string, empty for all languages")
	flag.StringVar(&clientId, "client_id", "q6batx0epp608isickayubi39itsckt", "Client ID for the Twitch API")
	flag.IntVar(&minViewers, "min", 10, "skip channels with less viewers than min")
	flag.BoolVar(&loop, "loop", false, "")
	flag.Parse()

	twitchStreamsPager := twitch_log.NewTwitchStreamsPager(language, clientId)

	streamSet := make(map[int]bool)
	csvWriter := csv.NewWriter(os.Stdout)
	for {
		streams, end, err := twitchStreamsPager.Next()
		if err != nil {
			log.Fatalln("twitchStreamsPager.Next() error:", err)
		}

		for _, stream := range streams {
			if stream.Viewers < minViewers {
				end = true
				break
			}

			if streamSet[stream.Channel.Id] {
				continue
			}
			csvWriter.Write([]string{strconv.Itoa(stream.Channel.Id), stream.Channel.Name, strconv.Itoa(stream.Viewers)})
			streamSet[stream.Channel.Id] = true
		}
		csvWriter.Flush()

		if end {
			if loop {
				twitchStreamsPager.Reset()
				time.Sleep(120 * time.Second)
			} else {
				break
			}
		}
	}

	if err := csvWriter.Error(); err != nil {
		log.Fatalln("CSV error:", err)
	}
}
