package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Michael-Levitin/hezzlizer/internal/dto"
	"github.com/rs/zerolog/log"
)

func main() {
	var err error
	var b []byte
	var resp *http.Response
	var list dto.GetResponse

	resp, err = http.Get("http://127.0.0.1:8080//goods/list?limit=10&offset=1")
	if err != nil {
		log.Warn().Err(err).Msg("error sending request")
	} else {
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Warn().Err(err).Msg("error reading response")
		} else {
			err = json.Unmarshal(b, &list)
			if err != nil {
				log.Warn().Err(err).Msg("error unmarshalling response")
			} else {
				log.Info().Msg(fmt.Sprintf("%+v, returned: %d items", list.Meta, len(list.Goods)))
			}
		}
	}
}
