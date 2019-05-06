package utils

import "hash/fnv"

func HashFnv1a(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}
