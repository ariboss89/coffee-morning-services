package model

type InteractionRequest struct {
	ContentName string `db:"content_name,omitempty" json:"content_name"`
	Caption     string `db:"caption,omitempty" json:"caption"`
}
