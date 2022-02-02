package seedwork

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDocument interface {
	GetCollectionName() string
}

type Document struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	CreatedBy  string             `json:"createdBy" bson:"createdBy"`
	ModifiedAt *time.Time         `json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
	ModifiedBy *string            `json:"modifiedBy,omitempty" bson:"modifiedBy,omitempty"`
}
