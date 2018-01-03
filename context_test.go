package colly

import (
	"testing"
	"strconv"
	"fmt"
)

func TestContextIteration(t *testing.T) {
	ctx := NewContext()
	for i := 0; i < 10; i++ {
		ctx.Put(strconv.Itoa(i), i)
	}
	values := ctx.ForEach(func(k string, v interface{}) interface{} {
		return fmt.Sprintf("%s ==> %d", k, v.(int))
	})
	if len(values) != 10 {
		t.Fatal("fail to iterate context")
	}
	for _, s := range values {
		t.Log(s.(string))
	}
}
