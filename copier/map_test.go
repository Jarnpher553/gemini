package copier

import "testing"

type A struct {
	Name  string
	Age   int
	Other string
}

type B struct {
	Name  string
	Age   int
	Other *string
}

func TestCopy(t *testing.T) {
	a1 := A{
		Name: "lijianfeng",
		Age:  15,
	}

	a2 := A{
		Name:  "lijianfeng",
		Age:   15,
		Other: "other",
	}

	var b1 B

	err := Copy(&a1, &b1)
	if err != nil {
		t.FailNow()
	}

	if a1.Name != b1.Name {
		t.FailNow()
	}
	if a1.Age != b1.Age {
		t.FailNow()
	}
	if b1.Other != nil {
		t.FailNow()
	}

	var b2 B

	err = Copy(&a2, &b2)
	if err != nil {
		t.FailNow()
	}

	if a2.Name != b2.Name {
		t.FailNow()
	}
	if a2.Age != b2.Age {
		t.FailNow()
	}
	if b2.Other == nil {
		t.FailNow()
	}

	aSlice := []A{
		a1,
		a2,
	}

	var bSlice []B

	err = Copy(&aSlice, &bSlice)
	if err != nil {
		t.Error(`slice`)
	}

	t.Log(bSlice)

	if bSlice[0].Other != nil {
		t.Error(`index 0`)
	}
	if bSlice[1].Other == nil {
		t.Error(`index 1`)
	}
}
