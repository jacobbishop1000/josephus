package main

import (
	"fmt"
	"sync"
	"testing"
)

// Setup functions

func TestSetSucc(t *testing.T) {
	victim1 := &Victim{ident: 1, step: 2, last: new(int), alive: true}
	victim2 := &Victim{ident: 2, step: 2, last: new(int), alive: true}

	victim1.SetSucc(victim2)

	if victim1.succ != victim2 {
		t.Errorf("SetSucc should assign succ attribute.")
	}
}
func TestSetPred(t *testing.T) {
	victim1 := &Victim{ident: 1, step: 2, last: new(int), alive: true}
	victim2 := &Victim{ident: 2, step: 2, last: new(int), alive: true}

	victim1.SetPred(victim2)

	if victim1.pred != victim2 {
		t.Errorf("SetPred should assign succ attribute.")
	}
}

func TestVictimRingLength(t *testing.T) {
	var Last int
	victims := makeVictimRing(5, &Last, 2)

	if len(victims) != 5 {
		t.Errorf("makeVictimRing(5, _, _) should have length 5")
	}
}

func TestVictimRingIds(t *testing.T) {
	var Last int
	victims := makeVictimRing(5, &Last, 2)

	for i, v := range victims {
		if v.ident != i+1 {
			t.Errorf("makeVictimRing(5, _, _) elements should have numerical identifiers 1 through n.")
		}
	}
}

func TestVictimRingSuccLinks(t *testing.T) {
	var Last int
	victims := makeVictimRing(5, &Last, 2)

	for i, v := range victims {
		if i < 4 && v.succ != victims[i+1] {
			t.Errorf("makeVictimRing %dth element should have successor in %dth element.", i, i+1)
		} else if i == 4 && v.succ != victims[0] {
			t.Errorf("makeVictimRing last element should have successor in first element.")
		}
	}
}

func TestVictimRingPredLinks(t *testing.T) {
	var Last int
	victims := makeVictimRing(5, &Last, 2)

	for i, v := range victims {
		if i > 0 && v.pred != victims[i-1] {
			t.Errorf("makeVictimRing %dth element should have predecessor in %dth element", i, i-1)
		} else if i == 0 && v.pred != victims[4] {
			t.Errorf("makeVictimRing first element should have predecessor in last element.")
		}
	}
}

// Sequential functions

func TestNewSuccLiveVictim(t *testing.T) {
	var Last int
	var k = 2

	victims := makeVictimRing(4, &Last, k)

	if len(victims) < 4 {
		t.Errorf("Test cannot run without a correct makeVictimRing function.")
		return
	}

	// Successor should be set on victims[0]
	victims[1].NewSucc(victims[3])

	if victims[1].succ != victims[3] {
		t.Errorf("Live node can be assigned a new succcessor")
	}
}

func TestNewSuccDeadVictim1(t *testing.T) {
	var Last int
	var k = 2

	victims := makeVictimRing(4, &Last, k)
	if len(victims) < 4 {
		t.Errorf("Test cannot run without a correct makeVictimRing function.")
		return
	}

	victims[2].alive = false

	// Successor should be set on victims[1]
	victims[2].NewSucc(victims[3])

	if victims[1].succ != victims[3] {
		t.Errorf("Successor of a dead node should be set on its predecessor.")
	}
}

func TestNewSuccDeadVictim2(t *testing.T) {
	var Last int
	var k = 2

	victims := makeVictimRing(4, &Last, k)
	if len(victims) < 4 {
		t.Errorf("Test cannot run without a correct makeVictimRing function.")
		return
	}
	victims[1].alive = false
	victims[2].alive = false

	// Successor should be set on victims[0]
	victims[2].NewSucc(victims[3])

	if victims[0].succ != victims[3] {
		t.Errorf("Successor of a dead node should be set on its predecessor.")
	}
}

func TestNewPredLiveVictim(t *testing.T) {
	var Last int
	var k = 2

	victims := makeVictimRing(4, &Last, k)
	if len(victims) < 4 {
		t.Errorf("Test cannot run without a correct makeVictimRing function.")
		return
	}
	// Predecessor should be set on victims[2]
	victims[2].NewPred(victims[0])

	if victims[2].pred != victims[0] {
		t.Errorf("Live node can be assigned a new predcessor")
	}
}

func TestNewPredDeadVictim1(t *testing.T) {
	var Last int
	var k = 2

	victims := makeVictimRing(4, &Last, k)
	if len(victims) < 4 {
		t.Errorf("Test cannot run without a correct makeVictimRing function.")
		return
	}
	victims[1].alive = false

	// Predecessor should be set on victims[2]
	victims[1].NewPred(victims[0])

	if victims[2].pred != victims[0] {
		t.Errorf("Predecessor of a dead node should be set on its successor.")
	}
}

func TestNewPredDeadVictim2(t *testing.T) {
	var Last int
	var k = 2

	victims := makeVictimRing(4, &Last, k)
	if len(victims) < 4 {
		t.Errorf("Test cannot run without a correct makeVictimRing function.")
		return
	}
	victims[1].alive = false
	victims[2].alive = false

	// Predecessors should be set on victims[0]
	victims[1].NewPred(victims[0])

	if victims[3].pred != victims[0] {
		t.Errorf("Predecessor of a dead node should be set on its successor.")
	}
}

func TestJosephusStep2(t *testing.T) {
	var tests = []struct {
		n        int
		step     int
		expected int
	}{
		{1, 2, 1},
		{2, 2, 1},
		{3, 2, 3},
		{4, 2, 1},
		{5, 2, 3},
		{6, 2, 5},
		{7, 2, 7},
		{8, 2, 1},
		{9, 2, 3},
		{10, 2, 5},
		{11, 2, 7},
		{12, 2, 9},
		{13, 2, 11},
		{14, 2, 13},
		{15, 2, 15},
		{16, 2, 1},
	}
	for _, test := range tests {
		actual := Josephus(test.n, test.step)
		if actual != test.expected {
			t.Errorf("Flavius(%d, %d) != %d (actual: %d)", test.n, test.step, test.expected, actual)
		}
	}
}

func TestJosephusStep3(t *testing.T) {
	var tests = []struct {
		n        int
		step     int
		expected int
	}{
		{1, 3, 1},
		{2, 3, 2},
		{3, 3, 2},
		{4, 3, 1},
		{5, 3, 4},
		{6, 3, 1},
		{7, 3, 4},
		{8, 3, 7},
		{9, 3, 1},
		{10, 3, 4},
		{11, 3, 7},
		{12, 3, 10},
		{13, 3, 13},
		{14, 3, 2},
		{15, 3, 5},
		{16, 3, 8},
	}
	for _, test := range tests {
		actual := Josephus(test.n, test.step)
		if actual != test.expected {
			t.Errorf("Flavius(%d, %d) != %d (actual: %d)", test.n, test.step, test.expected, actual)
		}
	}
}

// Asynchronous functions: with synchronization counter

func TestNewSuccAsyncLiveVictim(t *testing.T) {
	var sc sync.WaitGroup
	var Last int
	var k = 2

	victims := makeVictimRing(4, &Last, k)
	if len(victims) < 4 {
		t.Errorf("Test cannot run without a correct makeVictimRing function.")
		return
	}
	sc.Add(1)
	// Successor should be set on victims[0]
	go victims[1].NewSuccAsync(victims[3], &sc)
	sc.Wait()

	if victims[1].succ != victims[3] {
		t.Errorf("Live node can be assigned a new succcessor")
	}

}

func TestNewSuccAsyncDeadVictim1(t *testing.T) {
	var sc sync.WaitGroup
	var Last int
	var k = 2

	victims := makeVictimRing(4, &Last, k)
	if len(victims) < 4 {
		t.Errorf("Test cannot run without a correct makeVictimRing function.")
		return
	}
	victims[2].alive = false

	sc.Add(1)
	// Successor should be set on victims[1]
	go victims[2].NewSuccAsync(victims[3], &sc)
	sc.Wait()

	if victims[1].succ != victims[3] {
		t.Errorf("Successor of a dead node should be set on its predecessor.")
	}

}

func TestNewSuccAsyncDeadVictim2(t *testing.T) {
	var sc sync.WaitGroup
	var Last int
	var k = 2

	victims := makeVictimRing(4, &Last, k)
	if len(victims) < 4 {
		t.Errorf("Test cannot run without a correct makeVictimRing function.")
		return
	}
	victims[1].alive = false
	victims[2].alive = false

	sc.Add(1)
	// Successor should be set on victims[0]
	go victims[2].NewSuccAsync(victims[3], &sc)
	sc.Wait()

	if victims[0].succ != victims[3] {
		t.Errorf("Successor of a dead node should be set on its predecessor.")
	}
}

func TestNewPredAsyncLiveVictim(t *testing.T) {
	var sc sync.WaitGroup
	var Last int
	var k = 2

	victims := makeVictimRing(4, &Last, k)
	if len(victims) < 4 {
		t.Errorf("Test cannot run without a correct makeVictimRing function.")
		return
	}
	sc.Add(1)
	// Predecessor should be set on victims[2]
	go victims[2].NewPredAsync(victims[0], &sc)
	sc.Wait()

	if victims[2].pred != victims[0] {
		t.Errorf("Live node can be assigned a new predcessor")
	}

}

func TestNewPredAsyncDeadVictim1(t *testing.T) {
	var sc sync.WaitGroup
	var Last int
	var k = 2

	victims := makeVictimRing(4, &Last, k)
	if len(victims) < 4 {
		t.Errorf("Test cannot run without a correct makeVictimRing function.")
		return
	}
	victims[1].alive = false

	sc.Add(1)
	// Predecessor should be set on victims[2]
	go victims[1].NewPredAsync(victims[0], &sc)
	sc.Wait()

	if victims[2].pred != victims[0] {
		t.Errorf("Predecessor of a dead node should be set on its successor.")
	}
}

func TestNewPredAsyncDeadVictim2(t *testing.T) {
	var sc sync.WaitGroup
	var Last int
	var k = 2

	victims := makeVictimRing(4, &Last, k)
	if len(victims) < 4 {
		t.Errorf("Test cannot run without a correct makeVictimRing function.")
		return
	}
	victims[1].alive = false
	victims[2].alive = false

	sc.Add(1)
	// Predecessors should be set on victims[0]
	go victims[1].NewPredAsync(victims[0], &sc)
	sc.Wait()

	if victims[3].pred != victims[0] {
		t.Errorf("Predecessor of a dead node should be set on its successor.")
	}
}

func TestJosephusAsyncStep2(t *testing.T) {
	var tests = []struct {
		n        int
		step     int
		expected int
	}{
		{1, 2, 1},
		{2, 2, 1},
		{3, 2, 3},
		{4, 2, 1},
		{5, 2, 3},
		{6, 2, 5},
		{7, 2, 7},
		{8, 2, 1},
		{9, 2, 3},
		{10, 2, 5},
		{11, 2, 7},
		{12, 2, 9},
		{13, 2, 11},
		{14, 2, 13},
		{15, 2, 15},
		{16, 2, 1},
	}
	for _, test := range tests {
		actual := JosephusAsync(test.n, test.step)
		if actual != test.expected {
			t.Errorf("Flavius(%d, %d) != %d (actual: %d)", test.n, test.step, test.expected, actual)
		}
	}
}

func TestJosephusAsyncStep3(t *testing.T) {
	var tests = []struct {
		n        int
		step     int
		expected int
	}{
		{1, 3, 1},
		{2, 3, 2},
		{3, 3, 2},
		{4, 3, 1},
		{5, 3, 4},
		{6, 3, 1},
		{7, 3, 4},
		{8, 3, 7},
		{9, 3, 1},
		{10, 3, 4},
		{11, 3, 7},
		{12, 3, 10},
		{13, 3, 13},
		{14, 3, 2},
		{15, 3, 5},
		{16, 3, 8},
	}
	for _, test := range tests {
		actual := JosephusAsync(test.n, test.step)
		if actual != test.expected {
			t.Errorf("Flavius(%d, %d) != %d (actual: %d)", test.n, test.step, test.expected, actual)
		}
	}
}

// Asynchronous functions: with channels

func TestJosephusChanStep2(t *testing.T) {
	var tests = []struct {
		n        int
		step     int
		expected int
	}{
		{1, 2, 1},
		{2, 2, 1},
		{3, 2, 3},
		{4, 2, 1},
		{5, 2, 3},
		{6, 2, 5},
		{7, 2, 7},
		{8, 2, 1},
		{9, 2, 3},
		{10, 2, 5},
		{11, 2, 7},
		{12, 2, 9},
		{13, 2, 11},
		{14, 2, 13},
		{15, 2, 15},
		{16, 2, 1},
	}
	for _, test := range tests {
		actual, _ := JosephusChan(test.n, test.step)
		if actual != test.expected {
			t.Errorf("Flavius(%d, %d) != %d (actual: %d)", test.n, test.step, test.expected, actual)
		}
	}
}

func TestJosephusChanStep3(t *testing.T) {
	var tests = []struct {
		n        int
		step     int
		expected int
	}{
		{1, 3, 1},
		{2, 3, 2},
		{3, 3, 2},
		{4, 3, 1},
		{5, 3, 4},
		{6, 3, 1},
		{7, 3, 4},
		{8, 3, 7},
		{9, 3, 1},
		{10, 3, 4},
		{11, 3, 7},
		{12, 3, 10},
		{13, 3, 13},
		{14, 3, 2},
		{15, 3, 5},
		{16, 3, 8},
	}
	for _, test := range tests {
		actual, _ := JosephusChan(test.n, test.step)
		if actual != test.expected {
			t.Errorf("Flavius(%d, %d) != %d (actual: %d)", test.n, test.step, test.expected, actual)
		}
	}
}

func ExampleJosephusChanStep2n3() {
	_, seq := JosephusChan(3, 2)
	fmt.Print(seq)
	// Output: 2 1 3
}

func ExampleJosephusChanStep2n16() {
	_, seq := JosephusChan(16, 2)
	fmt.Print(seq)
	// Output: 2 4 6 8 10 12 14 16 3 7 11 15 5 13 9 1
}

func ExampleJosephusChanStep3n3() {
	_, seq := JosephusChan(3, 3)
	fmt.Print(seq)
	// Output: 3 1 2
}

func ExampleJosephusChanStep3n16() {
	_, seq := JosephusChan(16, 3)
	fmt.Print(seq)
	// Output: 3 6 9 12 15 2 7 11 16 5 13 4 14 10 1 8
}
