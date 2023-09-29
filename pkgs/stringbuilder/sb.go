package stringbuilder

// copy by https://github.com/tredeske/u/blob/master/ustrings/builder.go

import (
	"bytes"
	"strconv"
)

// more convenient interface for building strings
type SB struct {
	Buff bytes.Buffer
}

func New(size int) (rv *SB) {
	rv = &SB{}
	rv.Buff.Grow(size)
	return
}

func (sb *SB) Append(s string) (self *SB) {
	sb.Buff.WriteString(s)
	return sb
}

func (sb *SB) AppendInt(i int) (self *SB) {
	sb.Buff.WriteString(strconv.Itoa(i))
	return sb
}

func (sb *SB) AppendArray(arr ...string) (self *SB) {
	for _, s := range arr {
		sb.Buff.WriteString(s)
	}
	return sb
}

func (sb *SB) String() (rv string) {
	return sb.Buff.String()
}

func (sb *SB) Reset() (self *SB) {
	sb.Buff.Reset()
	return sb
}

func (sb *SB) Truncate(to int) (self *SB) {
	sb.Buff.Truncate(to)
	return sb
}

func (sb *SB) Len() (rv int) {
	return sb.Buff.Len()
}
