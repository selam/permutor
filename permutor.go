package permutor

import (
	"math"
	"github.com/willf/bitset"
	"math/rand"
	"sync"
	"errors"
	"encoding/json"
	"io/ioutil"
	"time"
)

var (
	MinError = errors.New("min value can not bigger than max value")
	MaxError = errors.New("max value can not bigger than max value")
)

type rp  struct {
	MinLength          int
	MaxLength          int
	CurrentLength      int
	CharLength         int
	CurrentPossibility int
	MaxPossibility     int
	Letters            string
	Bitset             *bitset.BitSet
	l                  sync.Mutex
}

type RandomPermutation struct {
	rp
}

func NewPermutor(s string, minLength int, maxLength int) (*RandomPermutation, error) {
	if minLength > maxLength {
		return nil, MinError
	}

	if maxLength < minLength {
		return nil, MaxError
	}


	maxPosibility := int(math.Pow(float64(len(s)), float64(minLength)))
	return &RandomPermutation{
		rp{
			MinLength:minLength,
			CharLength: len(s),
			MaxLength:maxLength,
			CurrentLength: minLength,
			CurrentPossibility:0,
			MaxPossibility: maxPosibility,
			Letters: s,
			Bitset: bitset.New(uint(maxPosibility)),
		},
	}, nil
}

func (s *RandomPermutation) permute() (string) {
	index := s.rp.CurrentPossibility % s.rp.CharLength
	items := string(s.rp.Letters[index])
	i := s.rp.CurrentPossibility
	for n := 1; n < s.rp.CurrentLength; n++ {
		i = i / s.rp.CharLength
		index = i % s.rp.CharLength
		items += string(s.rp.Letters[index])
	}
	return items

}

func (s *RandomPermutation) Reset() {
	s.l.Lock()
	defer s.l.Unlock()
	s.reset()
}

func (s *RandomPermutation) reset() {
	s.MaxPossibility = int(math.Pow(float64(s.CharLength), float64(s.CurrentLength)))
	s.CurrentPossibility = 0
	s.Bitset = bitset.New(uint(s.MaxPossibility))
}

func (s *RandomPermutation) Generate() (string) {
	s.l.Lock()
	defer s.l.Unlock()
	for {
		if (s.Bitset.All()) {
			if (s.CurrentLength < s.MaxLength) {
				s.CurrentLength++
				s.reset()
				continue
			}
			return ""
		}
		rand.Seed(time.Now().UnixNano())
		s.CurrentPossibility = rand.Intn(s.MaxPossibility)
		if !s.Bitset.Test(uint(s.CurrentPossibility)) {
			s.Bitset.Set(uint(s.CurrentPossibility))
			break
		}

	}
	return s.permute()
}

func (s *RandomPermutation) marshal() ([]byte, error) {
	return json.Marshal(&s.rp)
}

func (s *RandomPermutation) MarshalJSON() ([]byte, error) {
	s.rp.l.Lock()
	defer s.rp.l.Unlock()
	return s.marshal()
}

func (s *RandomPermutation) unmarshal(b []byte) error {
	return json.Unmarshal(b, &s.rp)
}

func (s *RandomPermutation) UnmarshalJSON(b []byte) error {
	s.rp.l.Lock()
	defer s.rp.l.Unlock()
	return s.unmarshal(b)
}

func (s *RandomPermutation) SaveTo(fname string) (error) {
	s.rp.l.Lock()
	defer s.rp.l.Unlock()
	b, err := s.marshal()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fname, b, 0644)
}

func LoadFrom(fname string) (*RandomPermutation, error) {
	var err error
	var d []byte
	d, err = ioutil.ReadFile(fname)
	var s RandomPermutation
	if err != nil {
		return nil, err
	}

	err = s.unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &s, nil
}


