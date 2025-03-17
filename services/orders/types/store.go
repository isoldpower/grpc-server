package types

type ObjectStore[T any] interface {
	AddItem(*T) (error, bool)
	RemoveItem(string) (error, bool)
	GetItem(string) (*T, bool)
	ListItems(from int, to int) ([]*T, bool)
	UpdateItem(*T) (error, bool)
	OverwriteItem(*T) (error, bool)
}
