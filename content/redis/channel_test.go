package redis

import "fmt"

type AnyJSON struct {
	UID int32 `json:"uid"`
}

func (aj *AnyJSON) Topic() string {
	return "AnyJSON"
}

func ExampleChannel() {
	channel := &Channel{REDISADDR}
	fmt.Println("Sub")
	retchan, err := channel.Subscribe(&AnyJSON{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Pub")
	err = channel.Publish(&AnyJSON{3})
	if err != nil {
		fmt.Println(err)
	}
	ret := <-retchan
	fmt.Println(ret.(*AnyJSON).UID)
	//Output: Sub
	//Pub
	//3
}
