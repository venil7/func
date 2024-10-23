package task

import (
	"errors"
	"testing"
)

func TestTaskMap(t *testing.T) {
	r1 := (func() (int, error) {
		return 1, nil
	})
	r2 := Map(r1, func(v int) int {
		return v * 2
	})

	if res, err := r2(); res != 2 || err != nil {
		t.Errorf("Map failed")
	}
}

func TestTaskFlatMap(t *testing.T) {
	var t1 Task[int] = func() (int, error) {
		return 1, nil
	}

	t2 := FlatMap(t1, func(v int) Task[int] {
		return func() (int, error) {
			return v * 2, nil
		}
	})

	if res, err := t2(); res != 2 || err != nil {
		t.Errorf("FlatMap failed")
	}
}

func TestSequence1(t *testing.T) {
	seq := Sequence(Of(1), Of(2), Of(3))

	if res, err := seq(); len(res) != 3 || err != nil {
		t.Errorf("Sequence failed")
	}
}

func TestSequence2(t *testing.T) {
	seq := Sequence(Of(1), Of(2), Fail[int](errors.New("fail")))

	if res, err := seq(); res != nil || err == nil {
		t.Errorf("Sequence failed")
	}
}

func TestTraverse1(t *testing.T) {
	seq := []int{1, 2, 3, 4}
	trav := Traverse(seq, Of)

	if res, err := trav(); len(res) != len(seq) || err != nil {
		t.Errorf("Traverse failed")
	}
}

func TestTraverse2(t *testing.T) {
	seq := []int{1, 2, 3, 4}
	trav := Traverse(seq, func(i int) Task[int] {
		if i%2 == 0 {
			return Fail[int](errors.New("fail"))
		}
		return Of(i)
	})

	if res, err := trav(); res != nil || err == nil {
		t.Errorf("Traverse failed")
	}
}
