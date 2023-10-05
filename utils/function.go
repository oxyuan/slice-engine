package utils

import (
	mapset "github.com/deckarep/golang-set"
	"reflect"
	"regexp"
	"slice-engine/entity"
)

type Predicate[T any] func(T) bool

func (p Predicate[T]) Test(val T) bool {
	return p(val)
}

func AnyMatch[T any](list *[]T, predicate Predicate[T]) bool {
	for _, v := range *list {
		if predicate(v) {
			return true
		}
	}
	return false
}

func AnySetMatch[T entity.PreSufTuple | string | *regexp.Regexp](set mapset.Set, predicate Predicate[T]) bool {
	for item := range set.Iter() {
		if item == nil {
			continue
		}
		var v T
		if reflect.TypeOf(item).Kind() == reflect.Ptr {
			v = *item.(*T)
		} else {
			v = item.(T)
		}
		if predicate(v) {
			return true
		}
	}
	return false
}
