package models

import "github.com/google/uuid"

type Order struct {
    ID        uuid.UUID `json:"id"`
    ProductID int       `json:"product_id"`
    UserID    string    `json:"user_id"`
}