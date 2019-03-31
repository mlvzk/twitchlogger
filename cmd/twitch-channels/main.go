package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mlvzk/twitch_logger"
)

func main() {
	language := flag.String("language", "", "show only streams of specified locale ID string, empty for all languages")
	clientID := flag.String("client_id", "q6batx0epp608isickayubi39itsckt", "Client ID for the Twitch API")
	minViewers := flag.Int("min", 10, "skip channels with less viewers than min")
	loop := flag.Bool("loop", false, "")
	flag.Parse()

	twitchStreamsPager := twitch_logger.NewTwitchStreamsPager(*language, *clientID)

	streamSet := make(map[int]bool)
	csvWriter := csv.NewWriter(os.Stdout)
	for {
		streams, end, err := twitchStreamsPager.Next()
		if err != nil {
			log.Fatalln("twitchStreamsPager.Next() error:", err)
		}

		for _, stream := range streams {
			if stream.Viewers < *minViewers {
				end = true
				break
			}

			if streamSet[stream.Channel.ID] {
				continue
			}
			csvWriter.Write([]string{strconv.Itoa(stream.Channel.ID), stream.Channel.Name, strconv.Itoa(stream.Viewers)})
			streamSet[stream.Channel.ID] = true
		}
		csvWriter.Flush()

		if end {
			if *loop {
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
