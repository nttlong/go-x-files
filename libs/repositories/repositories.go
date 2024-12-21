package repositories

import "context"

// Entity represents a generic data entity. You might want to define
// a more specific interface or struct based on your application's needs.
type Entity interface {
	GetID() interface{} // Returns the ID of the entity
}

// Repository defines common database operations.
type Repository interface {
	Get(ctx context.Context, id interface{}) (Entity, error)
	List(ctx context.Context) ([]Entity, error)
	Create(ctx context.Context, entity Entity) error
	Update(ctx context.Context, entity Entity) error
	Delete(ctx context.Context, id interface{}) error
}
