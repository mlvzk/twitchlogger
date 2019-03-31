package twitch_logger

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync/atomic"
)

var (
	apiBaseURL = "https://api.twitch.tv/kraken"
)

type TwitchStreamsResponse struct {
	Total   int            `json:"_total"`
	Streams []TwitchStream `json:"streams"`
}

type TwitchStream struct {
	Viewers int           `json:"viewers"`
	Channel TwitchChannel `json:"channel"`
}

type TwitchChannel struct {
	ID   int    `json:"_id"`
	Name string `json:"name"`
}

type TwitchStreamsPager struct {
	offset   int32
	limit    int32
	language string
	clientID string
}

// NewTwitchStreamsPager returns the initial TwitchStreamsPager
func NewTwitchStreamsPager(language, clientID string) *TwitchStreamsPager {
	return &TwitchStreamsPager{
		offset:   0,
		limit:    100,
		language: language,
		clientID: clientID,
	}
}

// Next returns the TwitchStreams for the current page and progresses.
// bool is true if it hits the end.
func (pager *TwitchStreamsPager) Next() ([]TwitchStream, bool, error) {
	offset := atomic.SwapInt32(&pager.offset, pager.offset+100)
	request, err := http.NewRequest("GET", apiBaseURL+"/streams/?language="+pager.language+"&limit="+strconv.Itoa(int(pager.limit))+"&offset="+strconv.Itoa(int(offset)), nil)
	if err != nil {
		return []TwitchStream{}, false, err
	}

	request.Header.Add("Client-ID", pager.clientID)
	request.Header.Add("Accept", "application/vnd.twitchtv.v5+json")
	response, err := (&http.Client{}).Do(request)
	if err != nil {
		return []TwitchStream{}, false, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []TwitchStream{}, false, err
	}
	twitchStreamsResponse := TwitchStreamsResponse{}
	json.Unmarshal(body, &twitchStreamsResponse)

	pager.offset += pager.limit

	if twitchStreamsResponse.Total <= int(pager.offset) {
		return twitchStreamsResponse.Streams, true, nil
	}

	return twitchStreamsResponse.Streams, false, nil
}

// Reset resets the pager to it's initial position
func (pager *TwitchStreamsPager) Reset() {
	pager.offset = 0
}
