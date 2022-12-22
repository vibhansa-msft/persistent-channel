package alpha_sequence

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type alphaSequenceTestSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *alphaSequenceTestSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
}

func (suite *alphaSequenceTestSuite) TestCreateAlphaSequence() {
	as, err := CreateAlphaSequence(3)
	suite.assert.NotNil(as)
	suite.assert.Nil(err)

	as, err = CreateAlphaSequence(0)
	suite.assert.Nil(as)
	suite.assert.NotNil(err)

	as, err = CreateAlphaSequence(-5)
	suite.assert.Nil(as)
	suite.assert.NotNil(err)

	as, err = CreateAlphaSequence(5)
	suite.assert.NotNil(as)
	suite.assert.Nil(err)
}

func (suite *alphaSequenceTestSuite) TestAlphaSequenceValueCheck() {
	as, err := CreateAlphaSequence(5)
	suite.assert.NotNil(as)
	suite.assert.Nil(err)
	suite.assert.Equal(as.sequence, strings.Repeat("a", 5))
}

func (suite *alphaSequenceTestSuite) TestAlphaSequenceNext() {
	as, err := CreateAlphaSequence(5)
	suite.assert.NotNil(as)
	suite.assert.Nil(err)
	suite.assert.Equal(as.sequence, strings.Repeat("a", 5))

	suite.assert.Equal(as.Next(), strings.Repeat("a", 5))
	suite.assert.Equal(as.Next(), "aaaab")
	suite.assert.Equal(as.Next(), "aaaac")
}

func (suite *alphaSequenceTestSuite) TestAlphaSequenceLoop() {
	as, err := CreateAlphaSequence(5)
	suite.assert.NotNil(as)
	suite.assert.Nil(err)

	for i := 0; i < 25; i++ {
		_ = as.Next()
	}

	suite.assert.Equal(as.Next(), "aaaaz")

	as.Next()
	suite.assert.Equal(as.Next(), "aaabb")
	for i := 0; i < 25; i++ {
		_ = as.Next()
	}
	suite.assert.Equal(as.Next(), "aaacb")
}

func (suite *alphaSequenceTestSuite) TestAlphaSequenceRollover() {
	as, err := CreateAlphaSequence(2)
	suite.assert.NotNil(as)
	suite.assert.Nil(err)

	for j := 0; j < 26; j++ {
		for i := 0; i < 25; i++ {
			_ = as.Next()
		}
		_ = as.Next()
	}
	suite.assert.Equal(as.Next(), "aa")
}

func (suite *alphaSequenceTestSuite) TestAlphaSequenceSet() {
	as, err := CreateAlphaSequence(2)
	suite.assert.NotNil(as)
	suite.assert.Nil(err)

	as.Set(3)
	suite.assert.Equal(as.Next(), "ac")
}

func (suite *alphaSequenceTestSuite) TestAlphaSequenceSetCaptialize() {
	as, err := CreateAlphaSequenceCaps(2)
	suite.assert.NotNil(as)
	suite.assert.Nil(err)

	as.Set(3)
	suite.assert.Equal(as.Next(), "AC")
	suite.assert.NotEqual(as.Next(), "ac")
}

func (suite *alphaSequenceTestSuite) TestAlphaSequenceSetString() {
	as, err := CreateAlphaSequence(2)
	suite.assert.NotNil(as)
	suite.assert.Nil(err)

	as.SetString("bd")
	suite.assert.Equal(as.Next(), "bd")
	suite.assert.Equal(as.Next(), "be")
	suite.assert.Equal(as.Next(), "bf")
}

func (suite *alphaSequenceTestSuite) TestAlphaSequenceDecrement() {
	as, err := CreateAlphaSequence(2)
	suite.assert.NotNil(as)
	suite.assert.Nil(err)

	as.SetString("bd")
	suite.assert.Equal(as.Prev(), "bd")
	suite.assert.Equal(as.Prev(), "bc")
	suite.assert.Equal(as.Prev(), "bb")

	as.SetString("aa")
	_ = as.Prev()
	suite.assert.Equal(as.Prev(), "zz")
}

func TestAlphaSequence(t *testing.T) {
	suite.Run(t, new(alphaSequenceTestSuite))
}
