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
	MinError = errors.New("min value can not be bigger than max value")
	MaxError = errors.New("max value can not be smaller than min value")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type rp  struct {
	sync.Mutex
	MinLength          int
	MaxLength          int
	CurrentLength      int
	CharLength         int
	CurrentPossibility int
	MaxPossibility     int
	Letters            string
	Bitset             *bitset.BitSet

}

type RandomPermutation struct {
	rp
}
func max(a, b int) (int) {
	if a < b {
		return b
	}
	return a
}

func NewPermutor(s string, minLength int, maxLength int) (*RandomPermutation, error) {
	// seed rand

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
	s.Lock()
	defer s.Unlock()
	s.reset()
}

func (s *RandomPermutation) reset() {
	s.MaxPossibility = int(math.Pow(float64(s.CharLength), float64(s.CurrentLength)))
	s.CurrentPossibility = 0
	s.Bitset = bitset.New(uint(s.MaxPossibility))
}

func (s *RandomPermutation) Generate() (string, int)  {
	s.rp.Lock()
	defer s.rp.Unlock()

	i := 0
	for {
		s.CurrentPossibility = rand.Intn(max(s.MaxPossibility, s.CurrentPossibility))
		if !s.Bitset.Test(uint(s.CurrentPossibility)) {
			s.Bitset.Set(uint(s.CurrentPossibility))
			break
		}
		// when letters length is big, checking all bits getting slower so if we cant find in 8 loop then check all bits
		if i % 20 == 0 {
			if (s.Bitset.All()) {
				if (s.CurrentLength < s.MaxLength) {
					s.CurrentLength++
					s.reset()
					continue
				}
				return "", 0
			}
		}
		i++
	}
	pos := s.permute()
	return pos, s.CurrentPossibility
}

func (s *RandomPermutation) marshal() ([]byte, error) {
	return json.Marshal(&s.rp)
}

func (s *RandomPermutation) MarshalJSON() ([]byte, error) {
	s.Lock()
	defer s.Unlock()
	return s.marshal()
}

func (s *RandomPermutation) unmarshal(b []byte) error {
	return json.Unmarshal(b, &s.rp)
}

func (s *RandomPermutation) UnmarshalJSON(b []byte) error {
	s.Lock()
	defer s.Unlock()
	return s.unmarshal(b)
}

func (s *RandomPermutation) SaveTo(fname string) (error) {
	s.Lock()
	defer s.Unlock()
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


