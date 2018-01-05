package colly

import (
	"strconv"
	"testing"
)

func TestContextIteration(t *testing.T) {
	ctx := NewContext()
	for i := 0; i < 10; i++ {
		ctx.Put(strconv.Itoa(i), i)
	}
	values := ctx.ForEach(func(k string, v interface{}) interface{} {
		return v.(int)
	})
	if len(values) != 10 {
		t.Fatal("fail to iterate context")
	}
	for _, i := range values {
		v := i.(int)
		if v != ctx.GetAny(strconv.Itoa(v)).(int) {
			t.Fatal("value not equal")
		}
	}
}
