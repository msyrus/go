package matcher

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/mock/gomock"
)

type lenMatcher struct {
	x int
}

func (l lenMatcher) Matches(x interface{}) bool {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Chan, reflect.Map, reflect.Array, reflect.Slice, reflect.String:
		return v.Len() == l.x
	}
	return false
}

func (l lenMatcher) String() string {
	return fmt.Sprintf("is equal lenth to %d", l.x)
}

type andMatcher struct {
	ms []gomock.Matcher
}

func (a andMatcher) Matches(x interface{}) bool {
	for _, m := range a.ms {
		if !m.Matches(x) {
			return false
		}
	}
	return true
}

func (a andMatcher) String() string {
	strs := []string{}
	for _, m := range a.ms {
		strs = append(strs, m.String())
	}
	return strings.Join(strs, " && ")
}

type orMatcher struct {
	ms []gomock.Matcher
}

func (o orMatcher) Matches(x interface{}) bool {
	for _, m := range o.ms {
		if m.Matches(x) {
			return true
		}
	}
	return false
}

func (o orMatcher) String() string {
	strs := []string{}
	for _, m := range o.ms {
		strs = append(strs, m.String())
	}
	return strings.Join(strs, " || ")
}

// And matches all the matchers
func And(m gomock.Matcher, ms ...gomock.Matcher) gomock.Matcher {
	mcs := append([]gomock.Matcher{m}, ms...)
	return andMatcher{mcs}
}

// Or matches any of the matchers
func Or(m gomock.Matcher, ms ...gomock.Matcher) gomock.Matcher {
	mcs := append([]gomock.Matcher{m}, ms...)
	return orMatcher{mcs}
}

// Len matches the length of chan, map, array, slice or string
func Len(l int) gomock.Matcher {
	return lenMatcher{l}
}
