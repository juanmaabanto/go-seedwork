package seedwork

import (
	"time"
)

type IDocument interface {
	GetCollectionName() string
}

type Document struct {
	Id         int64      `json:"id" bson:"_id,omitempty"`
	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
	CreatedBy  string     `json:"createdBy" bson:"createdBy"`
	ModifiedAt *time.Time `json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
	ModifiedBy *string    `json:"modifiedBy,omitempty" bson:"modifiedBy,omitempty"`
}
