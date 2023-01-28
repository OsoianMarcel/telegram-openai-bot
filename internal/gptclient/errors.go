package gptclient

import "errors"

var (
	ErrRespNoChoices = errors.New("the response has not choices")
	ErrRespEmptyText = errors.New("the response choice text is empty")
)
