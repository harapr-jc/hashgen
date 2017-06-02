// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package hashgen

import (
	"fmt"
	"testing"
	"time"
)

func TestAccumulate(t *testing.T) {

	var s Stats
	for i := 0; i < 20; i++ {
		start := time.Now()
		time.Sleep(100 * time.Millisecond)
		s.Accumulate(start)
	}
	// Because accumulate is in the background, let some time elapse
	// before looking at results
	time.Sleep(300 * time.Millisecond)
	s.RLock()
	count := s.requestCount
	elapsed := s.elapsed.Seconds()
	fmt.Printf("count: %d\n", s.requestCount)
	fmt.Printf("elapsed: %f\n", elapsed)
	s.RUnlock()

	if count != 20 {
		t.Errorf("Expected count = %d, actual = %d", 20, count)
	}
	// TODO: c'mon, you know this does not work for compare
	if elapsed < float64(1) {
		t.Errorf("Expected elapsed greater than 2.0, actual = %f", elapsed)
	}
}

func TestGetJson(t *testing.T) {

	var s Stats
	for i := 0; i < 20; i++ {
		start := time.Now()
		time.Sleep(100 * time.Millisecond)
		s.Accumulate(start)
	}
	// Because accumulate is in the background, let some time elapse
	// before looking at results
	time.Sleep(100 * time.Millisecond)

	b := s.GetJson()
	fmt.Println(string(b) + "\n")
}
