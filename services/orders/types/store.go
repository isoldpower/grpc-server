package types

type ObjectStore[T any] interface {
	AddItem(*T) (error, bool)
	RemoveItem(string) (error, bool)
	GetItem(string) (*T, bool)
	ListItems(limit *uint64, offset *uint64) ([]*T, bool)
	UpdateItem(*T) (error, bool)
	OverwriteItem(*T) (error, bool)
}
