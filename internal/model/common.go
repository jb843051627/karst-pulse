package model

import "time"

type ID int64

type ListFilter struct {
	SpringID ID
	SensorID ID
	EventID  ID
	Status   string
	From     time.Time
	To       time.Time
	Limit    int
	Offset   int
}

func (f ListFilter) Normalized() ListFilter {
	if f.Limit <= 0 {
		f.Limit = 50
	}
	if f.Limit > 200 {
		f.Limit = 200
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	return f
}

type PageInfo struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

type APIList[T any] struct {
	Items []T      `json:"items"`
	Page  PageInfo `json:"page"`
}

type Health struct {
	Status    string    `json:"status"`
	Database  string    `json:"database"`
	StartedAt time.Time `json:"started_at"`
	Now       time.Time `json:"now"`
}
