package utils

import (
	"fmt"
	"testing"
)

func TestInterfaceIsNil(t *testing.T) {
	fmt.Println("InterfaceIsNil(nil)", InterfaceIsNil(nil))
	var ip *int
	fmt.Println("InterfaceIsNil(ip)", InterfaceIsNil(ip))
	var i int
	fmt.Println("InterfaceIsNil(i):", InterfaceIsNil(i))
	ip = &i
	fmt.Println("InterfaceIsNil(ip):", InterfaceIsNil(ip))
}

func ExampleHashFnv1a() {
	fmt.Printf("%s:%d\n", "", HashFnv1a(""))
	fmt.Printf("%s:%d\n", " ", HashFnv1a(" "))
	fmt.Printf("%s:%d\n", "yhc", HashFnv1a("yhc"))
	//Output:
	//:2166136261
	//  :621580159
	//yhc:1723544021
}

func TestPublishTCPAddr(t *testing.T) {
	fmt.Println(":80", "=>", PublishTCPAddr(":80"))
	fmt.Println("24.2342.23.0:80", "=>", PublishTCPAddr("24.2342.23.0:80"))
	fmt.Println("127.0.0.1:80", "=>", PublishTCPAddr("127.0.0.1:80"))
	fmt.Println("[::]:80", "=>", PublishTCPAddr("[::]:80"))
}

func TestID32GenWithSalt(t *testing.T) {
	fmt.Println(":", ID32GenWithSalt(""))
	fmt.Println(" :", ID32GenWithSalt(" "))
	fmt.Println("yhc:", ID32GenWithSalt("yhc"))
	fmt.Println("yhc:", ID32GenWithSalt("yhc"))
}

func TestID32Gen(t *testing.T) {
	res := make(map[uint32]int)
	for i := 0; i < 10000000; i++ {
		res[ID32Gen()]++
	}
	var collison int
	for r, v := range res {
		if v > 1 {
			fmt.Printf("%d:%d\n", r, v)
			collison++
		}
	}
	fmt.Println("collison:", collison)
}
