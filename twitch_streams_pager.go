package twitch_logger

import (
	"net/http"
	"encoding/json"
	"strconv"
	"io/ioutil"
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
	Id   int    `json:"_id"`
	Name string `json:"name"`
}

type TwitchStreamsPager struct {
	offset   int
	limit    int
	language string
	clientId string
}

func NewTwitchStreamsPager(language, clientId string) *TwitchStreamsPager {
	return &TwitchStreamsPager{
		offset: 0,
		limit: 100,
		language: language,
		clientId: clientId,
	}
}

func (pager *TwitchStreamsPager) Next() ([]TwitchStream, bool, error) {
	request, err := http.NewRequest("GET", apiBaseURL+"/streams/?language=" + pager.language + "&limit="+strconv.Itoa(pager.limit)+"&offset="+strconv.Itoa(pager.offset), nil)
	if err != nil {
		return []TwitchStream{}, false, err
	}

	request.Header.Add("Client-ID", pager.clientId)
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

	if twitchStreamsResponse.Total <= pager.offset {
		return twitchStreamsResponse.Streams, true, nil
	}

	return twitchStreamsResponse.Streams, false, nil
}

func (pager *TwitchStreamsPager) Reset() {
	pager.offset = 0
}
