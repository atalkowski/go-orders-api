package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID     uint64     `json:"order_id"`
	CustomerID  uuid.UUID  `json:"customer_id"`
	LineItems   []LineItem `json:"line_items"`
	CreatedAt   *time.Time `json:"created_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

// Note to add the json tags automatically ... installed the
// go install -v github.com/fatih/gomodifytags@v1.16.0 ....
// ... BUT from within this tool to make it work.
type LineItem struct {
	ItemID   uuid.UUID `json:"item_id"` // Can add ,omitempty if wished.
	Quantity uint      `json:"quantity"`
	Price    uint      `json:"price"`
}
