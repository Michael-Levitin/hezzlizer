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

func (h HezzlServer) GoodCreate(w http.ResponseWriter, r *http.Request) {
	item, err := getGoodCreate(r)
	if err != nil {
		log.Warn().Err(err).Msg("error reading parameters")
		fmt.Fprintln(w, "error reading parameters")
		return
	}

	info, err := h.logic.GoodCreate(context.TODO(), item)
	if err != nil {
		log.Warn().Err(err).Msg("error executing h.logic.GoodCreate")
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
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

func getGoodCreate(r *http.Request) (*dto.Item, error) {
	var err error
	if err = r.ParseForm(); err != nil {
		return nil, fmt.Errorf("ParseForm() err: %v", err)
	}

	queryParams := r.URL.Query()
	projectId := queryParams.Get("projectId")

	var projectIdNum int
	if projectId != "" {
		projectIdNum, err = strconv.Atoi(projectId)
		if err != nil {
			log.Info().Err(err).Msg("couldn't get projectId")
			return nil, fmt.Errorf("couldn't parse URL parameters")
		}
	}
	item := dto.Item{
		ProjectID:   projectIdNum,
		Name:        r.FormValue("name"),
	}

	return &item, nil
}
