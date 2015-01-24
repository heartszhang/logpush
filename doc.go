package main

type doc map[string]interface{}

func (this doc) merge(rhs doc) {
	for k, v := range rhs {
		this[k] = v
	}
}
