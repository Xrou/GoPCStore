package model

type Computer struct {
	ID             string                 `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string                 `json:"name"`
	Price          float32                `json:"price"`
	Rating         float32                `json:"rating"`
	Specifications map[string]interface{} `json:"specifications"`
}
