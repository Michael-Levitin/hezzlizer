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

type Body struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func main() {
	zerolog.SetGlobalLevel(-1) // -1 = Trace
	//for i := 0; i < 10; i++ {
	//	goodCreate()
	//}
	//for i := 0; i < 10; i++ {
	//	goodUpdate()
	//}
	for i := 0; i < 10; i++ {
		goodRemove()
	}
	//for i := 0; i < 10; i++ {
	//	getList()
	//}

}

func goodUpdate() {
	notify := "update"
	resp, err := http.Post(randUrlUpdate(),
		"application/json",
		bytes.NewBuffer(randBody(randName(), randName())))
	respHandle(err, notify, resp)
}

func goodCreate() {
	notify := "create"
	resp, err := http.Post("http://127.0.0.1:8080/good/create?projectId=1",
		"application/json",
		bytes.NewBuffer(randBody(randName())))
	respHandle(err, notify, resp)
}

func goodRemove() {
	notify := "remove"
	resp, err := http.Post(randUrlRemove(),
		"application/json",
		bytes.NewBuffer(randBody()))
	respHandle(err, notify, resp)
}

func respHandle(err error, notify string, resp *http.Response) {
	if err != nil {
		log.Warn().Err(err).Msg(notify + ": error sending request")
		return
	}
	if resp.StatusCode != 200 {
		log.Warn().Msg(notify + fmt.Sprintf(": page returned status: %s", resp.Status))
		return
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn().Err(err).Msg(notify + ": error reading response")
		return
	}

	var item dto.Item
	err = json.Unmarshal(b, &item)
	if err != nil {
		log.Warn().Err(err).Msg(notify + ": error unmarshalling response " + string(b))
		return
	}
	log.Trace().Msg(fmt.Sprintf(notify+": %+v", item))
}

func getList() {
	var list dto.GetResponse

	resp, err := http.Get(randUrlList())
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

func randUrlUpdate() string {
	return fmt.Sprintf("http://127.0.0.1:8080/good/update?id=%d&projectId=1", randInt(20))
}

func randUrlRemove() string {
	return fmt.Sprintf("http://127.0.0.1:8080/good/remove?id=%d&projectId=1", randInt(20))
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
	return string(name)
}

func randBody(arg ...string) []byte {
	b := Body{}

	for i, field := range arg {
		switch {
		case i == 0:
			b.Name = field
		case i == 1:
			b.Description = field
		}
	}

	bts, err := json.Marshal(b)
	if err != nil {
		return []byte{}
	}

	return bts
}
