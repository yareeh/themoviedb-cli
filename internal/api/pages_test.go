package api

import (
	"errors"
	"reflect"
	"testing"
)

func TestAllPagesConcatenatesPages(t *testing.T) {
	var calls []int
	got, err := allPages(func(page int) ([]int, int, error) {
		calls = append(calls, page)
		return []int{page * 10, page*10 + 1}, 3, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{10, 11, 20, 21, 30, 31}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("results = %v, want %v", got, want)
	}
	if !reflect.DeepEqual(calls, []int{1, 2, 3}) {
		t.Errorf("pages fetched = %v, want [1 2 3]", calls)
	}
}

func TestAllPagesZeroTotal(t *testing.T) {
	calls := 0
	got, err := allPages(func(page int) ([]string, int, error) {
		calls++
		return nil, 0, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("results = %v, want empty", got)
	}
	if calls != 1 {
		t.Errorf("fetch called %d times, want 1", calls)
	}
}

func TestAllPagesErrorOnPage2(t *testing.T) {
	boom := errors.New("boom")
	got, err := allPages(func(page int) ([]int, int, error) {
		if page == 2 {
			return nil, 0, boom
		}
		return []int{page}, 3, nil
	})
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want %v", err, boom)
	}
	if got != nil {
		t.Errorf("results = %v, want nil on error", got)
	}
}

func TestAllPagesCapsRunawayTotal(t *testing.T) {
	calls := 0
	_, err := allPages(func(page int) ([]int, int, error) {
		calls++
		return []int{page}, 1_000_000, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != maxPages {
		t.Errorf("fetch called %d times, want cap %d", calls, maxPages)
	}
}
