// Idea:
// Van Roy & Seif Haridi, Concepts, Techniques, and Models of Computer Programming, MIT Press, 2003.
// Adapted for Go from the Oz kernel language.
package main

import (
	"fmt"
	"strconv"
	"sync"
)

type Victim struct {
	ident int
	step  int
	last  *int
	succ  *Victim
	pred  *Victim
	alive bool
}

// Initialization functions: don't modify

// Set victim's successor (for setup only)
func (v *Victim) SetSucc(s *Victim) {
	v.succ = s
}

// Set victims predecessor (for setup only)
func (v *Victim) SetPred(p *Victim) {
	v.pred = p
}

// Build the ring of victims, and return it.
func makeVictimRing(n int, Last *int, k int) [](*Victim) {
	var victims [](*Victim)

	//TODO
	for i := 0; i < n; i++ { //Make circle of Victims
		victims = append(victims, &Victim{ident: i + 1, last: Last, step: k, alive: true})
	}
	for i := 1; i < n; i++ { //Set predecessors for all Victims in circle
		victims[i].SetPred(victims[i-1])
	}
	for i := 0; i < n-1; i++ { //Set successors for all Victims in circle
		victims[i].SetSucc(victims[i+1])
	}
	//Make the linked list circular by connecting the ends
	victims[0].SetPred(victims[n-1])
	victims[n-1].SetSucc(victims[0])

	return victims
}

// ------------------------------------
// 1. Sequential algorithm

// Set v's successor to s if v is alive; otherwise, set
// the first live precedessor of v (use recursion).
func (v *Victim) NewSucc(s *Victim) {
	// TODO
	if v.alive == true {
		v.SetSucc(s)
	} else {
		v.pred.NewSucc(s) //Since v is dead, call function again on v's predecessor
	}
}

// Set v's predecessor to p if v is alive; otherwise, set
// the first live successor of v (use recursion).
func (v *Victim) NewPred(p *Victim) {
	// TODO
	if v.alive == true {
		v.SetPred(p)
	} else {
		v.succ.NewPred(p) //Since v is dead, call function again on v's successor
	}
}

// Kill the victim (if alive), or one of its successors otherwise.
// It takes as parameters x, the number live objects traversed so far and
// s, the number of survivors (the current victim included in both counts).
func (v *Victim) Kill(x int, s int) {

	if v.alive {
		if s == 1 {
			*v.last = v.ident
		} else if x%v.step == 0 {
			v.alive = false
			v.pred.NewSucc(v.succ)
			v.succ.NewPred(v.pred)
			v.succ.Kill(x+1, s-1)
		} else {
			v.succ.Kill(x+1, s)
		}
	} else {
		v.succ.Kill(x, s)
	}
}

func Josephus(n int, k int) int {

	var Last int

	var victims = makeVictimRing(n, &Last, k)

	if len(victims) < n {
		fmt.Errorf("List of victims was not initialized properly: exiting.")
		return 0
	}

	victims[0].Kill(1, n)

	return Last
}

// -------------------------------
// 2. Concurrent algorithm: with synchronization counter

// Set v's successor to s if v is alive; otherwise, set
// the first live precedessor of v. This is the asynchronous version,
// that recurses on a goroutine.
func (v *Victim) NewSuccAsync(s *Victim, c *sync.WaitGroup) {

	defer c.Done()

	// TODO
	if v.alive == true {
		v.SetSucc(s)
	} else {
		c.Add(1)
		go v.pred.NewSuccAsync(s, c)
	}
}

// Set v's predecessor to s if v is alive; otherwise, set
// the first live successor of v. This is the asynchronous version,
// that recurses on a goroutine.
func (v *Victim) NewPredAsync(p *Victim, c *sync.WaitGroup) {

	defer c.Done()

	// TODO
	if v.alive == true {
		v.SetPred(p)
	} else {
		c.Add(1)
		go v.succ.NewPredAsync(p, c) //Since v is dead, call function again on v's successor
	}
}

// Kill the victim (if alive), or one of its successors otherwise.
// It takes as parameters x, the number live objects traversed so far and
// s, the number of survivors (the current victim included in both counts).
//
// This version modifies Kill(x, s) to make it asynchronous
func (v *Victim) KillAsync(x int, s int, c *sync.WaitGroup) {

	defer c.Done()

	// TODO
	// all calls to NewSuccAsync, NewPredAsync, and KillAsync should be goroutines.
	// (Use the sync counter to synchronize.)
	if v.alive {
		if s == 1 {
			*v.last = v.ident
		} else if x%v.step == 0 {
			v.alive = false
			c.Add(1)
			go v.pred.NewSuccAsync(v.succ, c)
			c.Add(1)
			go v.succ.NewPredAsync(v.pred, c)
			c.Add(1)
			go v.succ.KillAsync(x+1, s-1, c)
		} else {
			c.Add(1)
			v.succ.KillAsync(x+1, s, c)
		}
	} else {
		c.Add(1)
		v.succ.KillAsync(x, s, c)
	}
}

func JosephusAsync(n int, k int) int {

	var Last int

	var victims = makeVictimRing(n, &Last, k)

	if len(victims) < n {
		fmt.Errorf("List of victims was not initialized properly: exiting.")
		return 0
	}

	// TODO: call KillAsync asynchronously
	// (Use the sync counter to synchronize the 2 routines.)

	var counter sync.WaitGroup
	counter.Add(1)
	firstVictim := victims[0]
	go firstVictim.KillAsync(1, n, &counter)
	counter.Wait() //Wait for the asynchronous call to be finished before returning function

	return Last
}

// -------------------------------
// 3. Concurrent algorithm: with channels
//
// This version does not make use of a synchronization counter:
// It uses channels instead.

// The logic for the KillChan function is similar to the KillAsync function's (above).
// However:
// + It uses the synchronization counter only for the shortcut functions.
// + Every victim's removal results in the victim's id being sent on the channel.
// + A termination signal is sent on the channel when done.
func (v *Victim) KillChan(x int, s int, c *sync.WaitGroup, killed chan int) {

	// TODO, by modifying KillAsync:
	//
	// + all calls to NewSuccChan, NewPredChan, and KillChan should be goroutines.
	// + use a channel to synchronize and receive ids of victims from callees
	if v.alive {
		if s == 1 {
			*v.last = v.ident
			killed <- *v.last
			close(killed) //Since the last Victim is standing, we've reached our answer and can close our channel
		} else if x%v.step == 0 {
			v.alive = false
			killed <- v.ident
			c.Add(1)
			go v.pred.NewSuccAsync(v.succ, c)
			c.Add(1)
			go v.succ.NewPredAsync(v.pred, c)
			c.Add(1)
			go v.succ.KillChan(x+1, s-1, c, killed)
		} else {
			c.Add(1)
			go v.succ.KillChan(x+1, s, c, killed)
		}
	} else {
		c.Add(1)
		go v.succ.KillChan(x, s, c, killed)
	}
}

// The JosephusChan function reads the ids of the victims
// from a channel, ensuring that it does not terminates until
// the last victim is known.
// It returns the id of the last standing victim, as well as
// a string of victims, in the order they were killed:
// E.g. "3 6 1 5 2 8 4 7" (see tests and examples).
//
// This function the same helper routines as above: NewPredAsync, NewSuccAsync
func JosephusChan(n int, k int) (int, string) {

	var Last int
	//var c sync.WaitGroup

	var victims = makeVictimRing(n, &Last, k)

	if len(victims) < n {
		fmt.Errorf("List of victims was not initialized properly: exiting.")
		return 0, ""
	}

	// TODO:

	// Initialize a channel
	killed := make(chan int)
	// Call KillChan
	var counter sync.WaitGroup
	go victims[0].KillChan(1, n, &counter, killed)
	// Read the channel to construct the list of victims
	str := ""
	recentVictim := 0
	lastVictimTracker := 0
	for range victims {
		recentVictim = <-killed
		if recentVictim != 0 {
			lastVictimTracker = recentVictim //had to keep track of last victim separately because killed channel always returned 0 during last read.
		}
		//print(recentVictim)
		//print(" ")

		//Used following source for converting int to str: https://stackoverflow.com/questions/10105935/how-to-convert-an-int-value-to-string-in-go
		str += strconv.Itoa(recentVictim)
		str += " "
	}
	return lastVictimTracker, str
}
