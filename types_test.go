package gosc

import (
	"testing"
	"time"
)

func tryTimetag(t *testing.T, a Timetag, b time.Time) {
	if a.Time() != b {
		t.Fatalf("failed timetag %x %v %v", a, a.Time(), b)
	}
}

func TestTimetagToTime(t *testing.T) {
	// 0*0.233ns = 0ns
	tryTimetag(t, Timetag(0), time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC))
	// "Immediately" special case
	tryTimetag(t, Timetag(1), time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC))
	// 2*0.233ns ~= 0ns
	tryTimetag(t, Timetag(2), time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC))
	// 5*0.233ns ~= 1ns
	tryTimetag(t, Timetag(5), time.Date(1900, 1, 1, 0, 0, 0, 1, time.UTC))
	// 1s
	tryTimetag(t, Timetag(1<<32+0), time.Date(1900, 1, 1, 0, 0, 1, 0, time.UTC))
}
