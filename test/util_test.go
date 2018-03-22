package test

import (
	"crd/pkg/apis/pmd/v1"
	"testing"
)

type InnerInnerStruct struct {
	f float64
}

type InnerStruct struct {
	d uint32
	e []uint16
	f InnerInnerStruct
}

type TestStruct struct {
	a int
	b int
	c InnerStruct
}

func TestCompare(t *testing.T) {
	val1 := TestStruct{1, 2,
		InnerStruct{4, []uint16{1, 2},
			InnerInnerStruct{64.0},
		},
	}
	val2 := TestStruct{1, 2, InnerStruct{4, []uint16{1, 2}, InnerInnerStruct{64.0}}}
	set := v1.CompareObjAndDiff(val1, val2)
	t.Log("dump the attrset", set)

}
