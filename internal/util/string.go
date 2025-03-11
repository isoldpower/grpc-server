package util

import (
	"fmt"
	"strings"
)

type PrefixApplier[T string] struct {
	prefix string
}

func NewPrefixApplier[T string](prefix string) *PrefixApplier[T] {
	return &PrefixApplier[T]{
		prefix: prefix,
	}
}

func (pa *PrefixApplier[T]) WithPrefix(key T) string {
	if strings.HasPrefix(string(key), pa.prefix) {
		return string(key)
	}

	return fmt.Sprintf("%s-%s", pa.prefix, string(key))
}
