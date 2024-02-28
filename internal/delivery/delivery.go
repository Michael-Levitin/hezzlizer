package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Michael-Levitin/hezzlizer/internal/dto"
	"github.com/Michael-Levitin/hezzlizer/internal/logic"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type HezzlServer struct {
	logic logic.HezzlLogicI
}

func NewHezzlServer(logic logic.HezzlLogicI) *HezzlServer {
	return &HezzlServer{logic: logic}
}

func (h HezzlServer) GoodsList(w http.ResponseWriter, r *http.Request) {
	meta, err := getMeta(r)
	if err != nil {
		log.Warn().Err(err).Msg("error reading parameters")
		fmt.Fprintln(w, "error reading parameters", err)
		return
	}

	info, err := h.logic.GoodsList(context.TODO(), meta)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Goods List:\n")
	json.NewEncoder(w).Encode(info)

}

func getMeta(r *http.Request) (*dto.Meta, error) {
	var meta dto.Meta
	var err error
	queryParams := r.URL.Query()
	offset := queryParams.Get("offset")
	if offset != "" {
		meta.Offset, err = strconv.Atoi(offset)
		if err != nil {
			log.Info().Err(err).Msg("couldn't get offset")
		}
	}
	if meta.Offset == 0 {
		meta.Offset = 1
		log.Info().Err(err).Msg("couldn't get offset, setting offset = 1")
	}
	limit := queryParams.Get("limit")
	if limit != "" {
		meta.Offset, err = strconv.Atoi(limit)
		if err != nil {
			log.Info().Err(err).Msg("couldn't get limit")
		}
	}
	if meta.Limit == 0 {
		meta.Limit = 10
		log.Info().Err(err).Msg("couldn't get limit, setting limit = 10")
	}
	return &meta, nil
}