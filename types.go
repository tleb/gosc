package gosc

import (
	"fmt"
	"time"
)

// PackageType is used to create comparable constants.
type PackageType string

// Constants representing the different types of packages that exist.
const (
	PackageTypeMessage = PackageType("message")
	PackageTypeBundle  = PackageType("bundle")
)

// Immediately is a specific Timetag representing immediate execution of a Bundle.
const Immediately = Timetag(1)
const timeTo1970 = 2208988800                         // Source: RFC 868
const nsPerFraction = float64(0.23283064365386962891) // 1e9/(2^32)

// Package is the generalization of both package types.
type Package interface {
	GetType() PackageType
}

// Message is the data structure for OSC message packets.
type Message struct {
	// The Address is a '/' separated string as per the specification
	Address string
	// Arguments is the array of that that is written when the package is sent.
	// Only data-types with defined writers is supported.
	Arguments []any
}

// GetType returns the package type for Messages
func (m *Message) GetType() PackageType {
	return PackageTypeMessage
}

func (m *Message) String() string {
	return fmt.Sprintf("%s: %v", m.Address, m.Arguments)
}

// Timetag represents the time since 1900-01-01 00:00. From the spec:
//
// Time tags are represented by a 64 bit fixed point number. The first 32 bits
// specify the number of seconds since midnight on January 1, 1900, and the
// last 32 bits specify fractional parts of a second to a precision of about
// 200 picoseconds. This is the representation used by Internet NTP
// timestamps. The time tag value consisting of 63 zero bits followed by a one
// in the least signifigant bit is a special case meaning “immediately.”
type Timetag uint64

// Bundle is the data structure for OSC bundle packets.
type Bundle struct {
	// Timetag for execution of the messages in this Bundle
	Timetag Timetag
	// List of messages to execute at Timetag. Messages are expected to be
	// handled atomically.
	Messages []*Message
	// Bundles can contain bundles, bundles in bundles are not handled
	// atomically.
	Bundles []*Bundle
	// Name is the name of the packet after the '#' when encoding. If omitted
	// this is set to 'bundle'.
	Name string
}

// GetType returns the package type for Bundle
func (b *Bundle) GetType() PackageType {
	return PackageTypeBundle
}

func getPadBytes(length int) int {
	return (4 - (length % 4)) % 4
}

// Fractions will return the fractions of the Timetag with picoseconds
// resolution.
func (tt Timetag) Fractions() uint32 {
	return uint32(tt)
}

// Seconds will return the seconds since the Timetag beginning.
func (tt Timetag) Seconds() uint32 {
	return uint32(tt >> 32)
}

// Time will return the Timetag as a Golang time.Time type.
func (tt Timetag) Time() time.Time {
	if uint64(tt) == 1 {
		// Means "immediately". It cannot occur otherwise as timetag == 0 gets
		// converted to January 1, 1900 while time.Time{} means year 1 in Go.
		// Use the time.Time.IsZero() method to detect it.
		return time.Time{}
	}

	return time.Unix(
		int64(tt.Seconds())-timeTo1970,
		int64(float64(tt.Fractions())*nsPerFraction),
	).UTC()
}
