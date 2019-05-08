package utils

import "hash/fnv"

// HashFnv1a 简化fvn1a算法调用
func HashFnv1a(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}
