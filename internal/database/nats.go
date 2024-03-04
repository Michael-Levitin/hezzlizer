package database

import (
	"bytes"
	"strconv"
	"time"

	"github.com/Michael-Levitin/hezzlizer/internal/dto"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type Nuts struct {
	nc *nats.Conn
}

func NewNatsSender(nc *nats.Conn) *Nuts {
	return &Nuts{nc: nc}
}

const (
	_queryStart = `
INSERT INTO goods
VALUES (Id, ProjectID, Name, Description, ProjectID, Removed, EventTime)
`
	tickerSec = 10   // ticker in seconds
	batchCap  = 1000 // batch capacity
)

var batch = make([]*dto.Item, 0, batchCap)

func (n Nuts) Send() {
	var jsonStr bytes.Buffer
	ticker := time.NewTicker(time.Second * tickerSec)

	for range ticker.C {
		n.nc.Publish("goods", []byte("Hello"))

		if len(batch) == 0 {
			continue
		}

		jsonStr.WriteString(_queryStart)
		for i, item := range batch {
			jsonStr.WriteString("(" + strconv.Itoa(item.Id) + ", ")
			jsonStr.WriteString(strconv.Itoa(item.ProjectID) + ", ")
			jsonStr.WriteString(item.Name + ", ")
			jsonStr.WriteString(item.Description + ", ")
			jsonStr.WriteString(strconv.Itoa(item.Priority) + ", ")
			jsonStr.WriteString(strconv.FormatBool(item.Removed) + ", ")
			jsonStr.WriteString(item.CreatedAt.String() + ")")
			if i < len(batch)-1 {
				jsonStr.WriteString(",\n")
			}
		}

		err := n.nc.Publish("goods", jsonStr.Bytes())
		if err != nil {
			log.Error().Err(err).Msg("nats: batch send failed")
		} else {
			log.Info().Msg("nats: batch sent" + jsonStr.String())
			log.Info().Msg("nats: batch sent")
			batch = make([]*dto.Item, 0, batchCap)
		}

		jsonStr.Reset()
	}
	log.Info().Msg("nats go finished")
}
