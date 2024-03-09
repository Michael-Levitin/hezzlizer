package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/Michael-Levitin/hezzlizer/internal/dto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(0) // -1 = Trace
	for i := 0; i < 10; i++ {
		getList(randUrlList())
	}

	for i := 0; i < 10; i++ {
		createGood(randName())
	}

}

func createGood(s string) {
	var item dto.Item

	resp, err := http.Post("http://127.0.0.1:8080/good/create?projectId=1", "application/json", bytes.NewBuffer([]byte(s)))
	if err != nil {
		log.Warn().Err(err).Msg("error sending request")
		return
	}
	if resp.StatusCode != 200 {
		log.Warn().Msg(fmt.Sprintf("Page returned status: %s", resp.Status))
		return
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn().Err(err).Msg("error reading response")
		return
	}
	err = json.Unmarshal(b, &item)
	if err != nil {
		log.Warn().Err(err).Msg("error unmarshalling response")
		return
	}
	log.Trace().Msg(fmt.Sprintf("Create: %+v", item))
}

func getList(url string) {
	var resp *http.Response
	var list dto.GetResponse
	var b []byte
	var err error
	resp, err = http.Get(url)
	if err != nil {
		log.Warn().Err(err).Msg("error sending request")
		return
	}
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Warn().Err(err).Msg("error reading response")
		return
	}
	err = json.Unmarshal(b, &list)
	if err != nil {
		log.Warn().Err(err).Msg("error unmarshalling response")
		return
	}
	log.Trace().Msg(fmt.Sprintf("%+v, returned: %d items", list.Meta, len(list.Goods)))
}

func randUrlList() string {
	return fmt.Sprintf("http://127.0.0.1:8080/goods/list?limit=%d&offset=%d", randInt(20), randInt(20))
}

func randInt(n int) int {
	rand.Seed(time.Now().UnixNano() + rand.Int63())
	return rand.Intn(n)
}

func randName() string {
	lenght := randInt(15) + 5
	name := make([]uint8, lenght)
	for i := 0; i < lenght; i++ {
		name[i] = uint8(randInt(26) + 97)
	}
	return `{"name":"` + string(name) + `"}`
}
