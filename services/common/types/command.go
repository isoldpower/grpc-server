package types

type Service interface {
	ExecuteExternal()
	Execute(ready chan<- bool) <-chan bool
}
