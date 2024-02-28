package dto

import (
	"fmt"
	"time"
)

var (
	QueryExecuteErorr = fmt.Errorf("could not execute query")
)

type Meta struct {
	Total  int `json:"total,omitempty" db:"total"`
	Remove int `json:"remove,omitempty" db:"remove"`
	Limit  int `json:"limit,omitempty" db:"limit"`
	Offset int `json:"offset,omitempty" db:"offset"`
}

type Item struct {
	Id          int       `json:"id" db:"id"`
	ProjectID   int       `json:"projectId" db:"projectId"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"Description"`
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
