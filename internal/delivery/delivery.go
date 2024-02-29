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
	item, err := getParam(r)
	if err != nil {
		log.Warn().Err(err).Msg("error reading parameters")
		fmt.Fprintln(w, err)
		return
	}

	item, err = h.logic.GoodCreate(context.Background(), item)
	if err != nil {
		log.Warn().Err(err).Msg("error executing h.logic.GoodCreate")
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h HezzlServer) GoodUpdate(w http.ResponseWriter, r *http.Request) {
	item, err := getParam(r)
	if err != nil {
		fmt.Fprintln(w, "error reading parameters: ", err)
		return
	}

	item, err = h.logic.GoodUpdate(context.Background(), item)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h HezzlServer) GoodRemove(w http.ResponseWriter, r *http.Request) {
	item, err := getParam(r)
	if err != nil {
		fmt.Fprintln(w, "error reading parameters: ", err)
		return
	}

	item, err = h.logic.GoodRemove(context.Background(), item)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	itemShort := dto.ItemShort{
		Id:        item.Id,
		ProjectID: item.ProjectID,
		Removed:   item.Removed,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(itemShort)
}

func (h HezzlServer) GoodsList(w http.ResponseWriter, r *http.Request) {
	meta, err := getMeta(r)
	if err != nil {
		log.Warn().Err(err).Msg("error reading parameters")
		fmt.Fprintln(w, "error reading parameters", err)
		return
	}

	info, err := h.logic.GoodsList(context.Background(), meta)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Goods List:\n")
	json.NewEncoder(w).Encode(info)

}

func (h HezzlServer) GoodReprioritize(w http.ResponseWriter, r *http.Request) {
	item, err := getParam(r)
	if err != nil {
		log.Warn().Err(err).Msg("error reading parameters")
		fmt.Fprintln(w, "error reading parameters", err)
		return
	}

	priorities, err := h.logic.GoodReprioritize(context.Background(), item)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "update item:\n")
	json.NewEncoder(w).Encode(priorities)

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

func getParam(r *http.Request) (*dto.Item, error) {
	var err error
	if err = r.ParseForm(); err != nil {
		return nil, fmt.Errorf("ParseForm() err: %v", err)
	}

	queryParams := r.URL.Query()
	projectId := queryParams.Get("projectId")
	var projectIdNum, idNum, priorityNum int
	if projectId != "" {
		projectIdNum, err = strconv.Atoi(projectId)
		if err != nil {
			log.Info().Err(err).Msg("couldn't get projectId")
			return nil, fmt.Errorf("couldn't parse URL parameters")
		}
	}
	id := queryParams.Get("id")
	if id != "" {
		idNum, err = strconv.Atoi(id)
		if err != nil {
			log.Info().Err(err).Msg("couldn't get projectId")
			return nil, fmt.Errorf("couldn't parse URL parameters")
		}
	}

	var item dto.Item
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse request Body parameters, %s", err)
	}

	item.Id = idNum
	item.ProjectID = projectIdNum
	item.Priority = priorityNum

	return &item, nil
}
