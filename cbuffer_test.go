package cbuffer

import "testing"

type D struct {
	value int
}

func (d D) Less(other D) bool {
	return d.value < other.value
}

func (d D) Equal(other D) bool {
	return d == other
}

func TestOrderedIntCB(t *testing.T) {
	cb := NewOrderedCircuitBuffer[D](3)

	t.Logf("CircuitBuffer cap: %v", cb.Cap())
	t.Logf("CircuitBuffer len: %v", cb.Len())

	cb.Add(D{1})
	cb.Add(D{2})
	cb.Add(D{3})
	cb.Add(D{4})

	zeroItem := D{2}

	if cb.GetItem(0) != zeroItem {
		t.Fatalf("%v != %v", cb.GetItem(0), zeroItem)
	} else {
		t.Logf("item at index 0 == %v", cb.GetItem(0))
	}
}

func TestOCBIteration(t *testing.T) {
	cb := NewOrderedCircuitBuffer[D](3)

	cb.Add(D{1})
	cb.Add(D{2})
	cb.Add(D{3})
	cb.Add(D{4})

	index := 0
	check_list := []D{{2}, {3}, {4}}

	// Check that no iteration object related to this OCB
	if cb.iter != nil {
		t.Fatalf("OCB %v must not contain an iter object", cb)
	}

	// Check iteration function
	for item := range cb.Iter() {
		t.Logf("Element %v: %v", index, *item)

		// Check that no iteration object created and related
		if cb.iter == nil && item.value != 4 {
			t.Fatalf("OCB %v must contain an iter object, but got nil", cb)
		}

		if !item.Equal(check_list[index]) {
			t.Fatalf("Index %v: %v != %v (%v)", index, *item, check_list[index], cb)
		}

		index++
	}

	// Check that iteration object removed correctly
	if cb.iter != nil {
		t.Fatalf("OCB %v must not contain an iter object after end of iteration", cb)
	}

}

func TestCBBreakIteration(t *testing.T) {
	cb := NewCircuitBuffer[D](3)

	for i := range []int{1, 3, 5, 4} {
		cb.Add(D{i})
	}

	index := 0

	// Check iteration function
	for item := range cb.Iter() {
		t.Logf("Element %v: %v", index, *item)

		// Check that no iteration object created and related
		if item.value != 3 {
			cb.Break()
			break
		}

		index++
	}

	// Check that iteration object removed correctly
	if cb.iter != nil {
		t.Fatalf("OCB %v must not contain an iter object after Break of iteration", cb)
	}
}

func TestOCBWrongBreakIteration(t *testing.T) {
	ocb := NewOrderedCircuitBuffer[D](3)

	for i := range []int{1, 3, 5, 4} {
		err := ocb.Add(D{i})

		if err != nil && i != 4 {
			t.Fatalf("%v can't be inserted into ocb %v", i, ocb)
		}
	}

}

func TestOCBSearch(t *testing.T) {
	ocb := NewOrderedCircuitBuffer[D](100)

	for i := 0; i < 100; i++ {
		err := ocb.Add(D{i})

		if err != nil {
			t.Fatalf("%v can't be inserted into ocb %s", i, ocb)
		}
	}

	expected := 55
	index, found := ocb.Search(D{expected})
	if !found {
		t.Fatalf("%v can't be found in ocb %s", expected, ocb)
	}

	if index != expected {
		t.Fatalf("%v have wrong index %v (expected 55) founded in ocb %s", expected, index, ocb)
	}

}

func TestOCBSearchNotFound(t *testing.T) {
	ocb := NewOrderedCircuitBuffer[D](50)

	for i := 0; i < 100; i++ {
		err := ocb.Add(D{i})

		if err != nil {
			t.Fatalf("%v can't be inserted into ocb %s", i, ocb)
		}
	}

	expected := 20
	index, found := ocb.Search(D{expected})
	if found {
		t.Fatalf("%v couldn't be found in ocb %s", expected, ocb)
	}

	if index != -1 {
		t.Fatalf("Expected index -1, got %v, while searching %v in %s", index, expected, ocb)
	}

}
