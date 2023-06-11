package model

import (
	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/model"
)

// SendDialogBody ...
type SendDialogBody struct {
	Text string `json:"text"`
}

// ListDialogResponse ...
type ListDialogResponse struct {
	From uuid.UUID `json:"from"`
	To   uuid.UUID `json:"to"`
	Text string    `json:"text"`
}

// NewDialogResponse ...
func NewDialogResponse(dialog model.Dialog) ListDialogResponse {
	return ListDialogResponse{
		From: dialog.From,
		To:   dialog.To,
		Text: dialog.Text,
	}
}

// NewDialogListResponse ...
func NewDialogListResponse(dialogs []model.Dialog) []ListDialogResponse {
	res := make([]ListDialogResponse, 0, len(dialogs))
	for _, dialog := range dialogs {
		res = append(res, NewDialogResponse(dialog))
	}

	return res
}
