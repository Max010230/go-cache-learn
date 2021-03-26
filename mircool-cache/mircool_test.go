package mircool_cache

import (
	"errors"
	"fmt"
	"testing"
)

func TestGetter(t *testing.T) {
	//var f Getter = GetterFunc(func(key string) ([]byte, error) {
	//	return []byte(key),nil
	//})
	//expect := []byte("aa")
	//if bytes, _ := f.Get("aa");!reflect.DeepEqual(bytes,expect){
	//	t.Errorf("callback faild")
	//}
	fmt.Println(DoTheThing(true))
	fmt.Println(DoTheThing(false))
}

var ErrDidNotWork = errors.New("did not work")

func DoTheThing(reallyDoIt bool) (err error) {
	if reallyDoIt {
		result, err := tryTheThing()
		if err != nil || result != "it worked" {
			return err
		}
	}
	return err
}

func tryTheThing() (string, error) {
	return "", ErrDidNotWork
}
