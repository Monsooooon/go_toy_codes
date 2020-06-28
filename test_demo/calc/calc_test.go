package main

import (
	"fmt"
	"os"
	"testing"
)

func setup() {
	fmt.Println("Before all tests")
}

func teardown() {
	fmt.Println("After all tests")
}

func TestAdd(t *testing.T) {
	cases := []struct {
		x, y, expect int
	}{
		{1, 2, 3},
		{-1, 1, 0},
		{-1, -1, -2},
	}

	for _, tc := range cases {
		if val := Add(tc.x, tc.y); val != tc.expect {
			t.Errorf("Add(%d, %d): expect %d, got %d", tc.x, tc.y, tc.expect, val)
		}
	}
}

func TestMul(t *testing.T) {
	cases := []struct {
		name         string
		x, y, expect int
	}{
		{"pos", 1, 2, 2},
		{"neg", -1, 1, -1},
		{"zero", -1, 0, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if ans := Mul(tc.x, tc.y); ans != tc.expect {
				t.Fatalf("Mul(%d, %d): expect %d, but got %d",
					tc.x, tc.y, tc.expect, ans)
			}
		})
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
