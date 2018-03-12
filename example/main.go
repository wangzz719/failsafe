package main

import (
	"github.com/wangzz719/failsafe/failsafe"
	"fmt"
	"errors"
)

func main() {
	failSafe := &failsafe.FailSafe{}
	defaultRtn := int(10)
	result := int(0)
	err := failSafe.Safe(defaultRtn, false, []string{"errorString"}, func() (interface{}, error) {
		i := int(0)
		return i, errors.New("fail safe not work")
	}, &result)
	fmt.Println(result)
	fmt.Println(err)

	err = failSafe.Safe(defaultRtn, true, []string{"errorString"}, func() (interface{}, error) {
		i := int(100)
		return i, errors.New("fail safe work")
	}, &result)
	fmt.Println(err)
	fmt.Println(result)
}
