package alpha_sequence

import (
	"fmt"
	"strings"
	"sync"
)

/*
	AlphaSequence : Represents a alphabet based sequence
	if user has given length as 3 then it starts with 'aaa' and on each nextSequence call it increments least index like
	'aaa', 'aab', 'aac'.... 'aaz', 'aba'.....'zzz', 'aaa'....
*/

type AlphaSequence struct {
	sequence   string
	maxLength  int
	mtx        sync.Mutex
	capitalize bool
}

// CreateAlphaSequence : Create object of AlphaSequence with a given length string set to 'aa....'
func CreateAlphaSequence(len int) (*AlphaSequence, error) {
	if len <= 0 {
		return nil, fmt.Errorf("length should be > 0")
	}

	as := &AlphaSequence{
		maxLength:  len,
		capitalize: false,
	}

	as.reset()
	return as, nil
}

// CreateAlphaSequenceCaps : Create object of AlphaSequence with a given length string set to 'AA....'
func CreateAlphaSequenceCaps(len int) (*AlphaSequence, error) {
	if len <= 0 {
		return nil, fmt.Errorf("length should be > 0")
	}

	as := &AlphaSequence{
		maxLength:  len,
		capitalize: true,
	}

	as.reset()
	return as, nil
}

// Reset : Reset the sequence to start position which is 'aa...'
func (as *AlphaSequence) Reset() {
	as.mtx.Lock()
	defer as.mtx.Unlock()
	as.reset()
}

// Set : Set the sequence to start at particular index
func (as *AlphaSequence) Set(index int) {
	as.mtx.Lock()
	defer as.mtx.Unlock()

	as.reset()
	for i := 1; i < index; i++ {
		as.increment()
	}
}

// SetString : Set the sequence to start with given value
func (as *AlphaSequence) SetString(val string) {
	as.mtx.Lock()
	defer as.mtx.Unlock()
	as.sequence = val
}

// GetNextSequence : Generate the next sequence number and return
func (as *AlphaSequence) Next() string {
	as.mtx.Lock()
	defer as.mtx.Unlock()

	val := as.sequence
	as.increment()
	return val
}

// GetNextSequence : Generate the prev sequence number and return
func (as *AlphaSequence) Prev() string {
	as.mtx.Lock()
	defer as.mtx.Unlock()

	val := as.sequence
	as.decrement()
	return val
}

// reset : Reset the value to first sequence
func (as *AlphaSequence) reset() {
	if as.capitalize {
		as.sequence = strings.Repeat("A", as.maxLength)
	} else {
		as.sequence = strings.Repeat("a", as.maxLength)
	}
}

// increment : Increment the sequence to next alphabet value (+1)
func (as *AlphaSequence) increment() {
	for i := as.maxLength - 1; i >= 0; i-- {
		if as.sequence[i] < 'z' ||
			(as.capitalize && as.sequence[i] < 'Z') {
			// There is still scope to increment at this index
			as.setChar(i, rune(as.sequence[i]+1))
			break
		} else {
			// This index is rolling over so move to next index and set this index to 'a'
			if as.capitalize {
				as.setChar(i, 'A')
			} else {
				as.setChar(i, 'a')
			}
		}
	}
}

// decrement : Decrement the sequence to next alphabet value (+1)
func (as *AlphaSequence) decrement() {
	for i := as.maxLength - 1; i >= 0; i-- {
		if as.sequence[i] > 'a' ||
			(as.capitalize && as.sequence[i] < 'A') {
			// There is still scope to increment at this index
			as.setChar(i, rune(as.sequence[i]-1))
			break
		} else {
			// This index is rolling over so move to next index and set this index to 'a'
			if as.capitalize {
				as.setChar(i, 'Z')
			} else {
				as.setChar(i, 'z')
			}
		}
	}
}

// setChar : Set a character at particular index in sequence
func (as *AlphaSequence) setChar(i int, c rune) {
	as.sequence = as.sequence[:i] + string(c) + as.sequence[i+1:]
}
