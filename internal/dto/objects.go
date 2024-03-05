package dto

import (
	"fmt"
	"time"
)

var (
	ErrQueryExecute = fmt.Errorf("could not execute query")
	ErrNotFound     = fmt.Errorf("errors.good.notFound")
)

type Meta struct {
	Total   int `json:"total" db:"total"`
	Removed int `json:"removed" db:"removed"`
	Limit   int `json:"limit" db:"limit"`
	Offset  int `json:"offset" db:"offset"`
}

type Item struct {
	Id          int       `json:"id" db:"id"`
	ProjectID   int       `json:"projectId" db:"projectId"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Priority    int       `json:"priority" db:"priority"`
	Removed     bool      `json:"removed" db:"removed"`
	CreatedAt   time.Time `json:"createdAt" db:"createdAt"`
}

type ItemShort struct {
	Id        int  `json:"id" db:"id"`
	ProjectID int  `json:"projectId" db:"projectId"`
	Removed   bool `json:"removed" db:"removed"`
}

type GetResponse struct {
	Meta  Meta
	Goods []Item
}

type Priority struct {
	Id       int `json:"id" db:"id"`
	Priority int `json:"priority" db:"priority"`
}

type ReprResponse struct {
	Priorities []Priority `json:"priorities" db:"priorities"`
}
