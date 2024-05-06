package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/dreamsofcode-io/orders-api/model"
	"github.com/dreamsofcode-io/orders-api/myutils"
	"github.com/dreamsofcode-io/orders-api/repository/order"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Order struct { // Slight question on naming convention for this Order Repo struct??
	Repo *order.RedisRepo
}

func (h *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create an order")
	var body struct { // anonymousinline body object : Notice this COOL way to set up a deserialize:
		CustomerID uuid.UUID        `json:"customer_id,omitempty"` // TODO omitempty is OK?!
		LineItems  []model.LineItem `json:"line_items,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("failed to decode order request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	order := model.Order{
		OrderID:    rand.Uint64(), // !!Not for PRODUCTION but OK for this demo code!!
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}
	err := h.Repo.Insert(r.Context(), order)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed to marshall response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (h *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all orders")
	cursor, err := myutils.GetUIntParam(r, "cursor", 0)
	if err != nil {
		fmt.Println("Invalid cursor value not numeric: ", myutils.GetQueryParam(r, "cursor", ""))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	const size = 50
	res, err := h.Repo.FindAll(r.Context(), order.FindOrders{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		fmt.Println("Failed to call repo FindAll:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Anonymous (inline) struct for the result:
	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty"`
	}
	response.Items = res.Orders
	response.Next = res.Cursor
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Failed to marshal FindAll result:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	//	w.WriteHeader(http.StatusOK) this is the default OK
}

func (h *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get and order by ID")
	orderId, err := myutils.StrToUInt(chi.URLParam(r, "id"))

	if err != nil {
		fmt.Println("Invalid orderId value")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := h.Repo.FindByID(r.Context(), orderId)
	if errors.Is(err, myutils.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id ", orderId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// NOte the logical equivalence of these two way s of doing this:
	if err := json.NewEncoder(w).Encode(value); err != nil {
		fmt.Println("Failed to marshal order result:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	// data, err := json.Marshal(value)
	// if err != nil {
	// 	fmt.Println("Failed to marshal order result:", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// w.Write(data)
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update an order by ID")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete an order by ID")
}
