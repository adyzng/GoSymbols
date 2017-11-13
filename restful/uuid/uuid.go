package uuid

import (
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	Size    = 16
	version = 0x20
	format  = "%08x-%04x-%04x-%04x-%012x"
)

var (
	ErrInvalidUUID = fmt.Errorf("invalid uuid string")
)

// UUID ...
type UUID [Size]byte

// Generate a random UUID
func Generate() (u UUID) {
	var (
		retry = 10
		count = 0
	)
	for ; retry >= 0; retry-- {
		n, err := io.ReadFull(rand.Reader, u[count:])
		if count == len(u) {
			break
		}

		if err != nil {
			if retry == 0 {
				panic(fmt.Sprintf("generate uuid failed with : %s", err))
			}
			count += n
		}
	}

	// set version
	u[4] = (u[4] & 0x0F) | version
	return
}

// NewUUID ...
func NewUUID() string {
	return Generate().RawString()
}

// Parse UUID from string
func Parse(s string) (u UUID, e error) {
	sLen := len(s)
	if sLen != 32 && sLen != 36 {
		return UUID{}, ErrInvalidUUID
	}

	if sLen == 36 {
		s = strings.Replace(s, "-", "", -1)
	}

	for i, n := 0, len(s); i < n; i += 2 {
		v, err := strconv.ParseUint(s[i:i+2], 16, 8)
		if err != nil {
			return UUID{}, err
		}

		u[i/2] = byte(v)
	}
	return
}

// String returns format 974AFFD3-1BCC-4475-8910-A967AFAE51FE
func (u UUID) String() string {
	return fmt.Sprintf(format, u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// RawString returns format 974AFFD31BCC44758910A967AFAE51FE
func (u UUID) RawString() string {
	return fmt.Sprintf("%x", u[:])
}

// Version return current version
func (u UUID) Version() byte {
	return (u[4] & 0xF0)
}
