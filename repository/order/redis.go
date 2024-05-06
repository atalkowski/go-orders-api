package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dreamsofcode-io/orders-api/model"
	"github.com/dreamsofcode-io/orders-api/myutils"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}
	key := orderIDKey(order.OrderID)
	fmt.Println("order key created", key)
	txn := r.Client.TxPipeline() // Start a redis transaction as this has 2 update parts:
	res := txn.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set order: %w", err)
	}
	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add order to set: %w", err)
	}
	if _, err := txn.Exec(ctx); err != nil {
		// I guess the Exec does not require a subsequent Discard as it failed.
		return fmt.Errorf("failed to exec insert order txn: %w", err)
	}
	return nil
}

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	key := orderIDKey(id)
	value, err := r.Client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return model.Order{}, myutils.ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("get order by id: %w", err)
	}
	var order model.Order
	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode order json: %w", err)
	}
	return order, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := orderIDKey(id)
	txn := r.Client.TxPipeline() // Start a redis transaction
	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return myutils.ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("delete order: %w", err)
	}
	// The SRem takes context,  set key="orders" and key for the specific member
	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove order from set: %w", err)
	}
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec delete order txn: %w", err)
	}
	return nil
}

func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("update failed to encode order: %w", err)
	}
	key := orderIDKey(order.OrderID)
	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return myutils.ErrNotExist
	}
	if err != nil {
		return fmt.Errorf("failed to set order: %w", err)
	}
	return nil
}

type FindOrders struct {
	Size   uint64
	Offset uint64
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, search FindOrders) (FindResult, error) {
	scanCmd := r.Client.SScan(ctx, "orders", search.Offset, "*", int64(search.Size))
	keys, cursor, err := scanCmd.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to scan order ids: %w", err)
	}
	// MGet won't like a request with enplty list of keys ... so this check is important
	if len(keys) == 0 {
		return FindResult{
			Orders: []model.Order{}, // note the Cursor will be defaulted to uint64 0
		}, nil
	}
	// The keys will in general give orders in Random order ... this can be changed (TODO)
	// Note use of the variadic array of keys... here and result xs is []interface
	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to MGet order ids: %w", err)
	}
	// OK all good - so create and orders slice
	orders := make([]model.Order, len(xs))
	for i, x := range xs {
		x := x.(string)
		var order model.Order
		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to deserial order json:%w", err)
		}
		orders[i] = order
	}
	return FindResult{Orders: orders, Cursor: cursor}, nil
}
