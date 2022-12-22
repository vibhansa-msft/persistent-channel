package pchannel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type pchannelTestSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *pchannelTestSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
}

const diskPath = "C:\\Users\\vibhansa\\Documents\\Projects\\persistent-channel"

func (suite *pchannelTestSuite) TestCreatePChannelSuccess() {
	p := &PChannel{}
	err := p.Init(PChannelConfig{
		PChannelID:    "",
		MaxMsgCount:   10,
		MaxCacheCount: 5,
		DiskPath:      diskPath,
	})

	suite.assert.Nil(err)

	err = p.Destroy()
	suite.assert.Nil(err)
}

func (suite *pchannelTestSuite) TestCreatePChannelWrongSize() {
	p := &PChannel{}
	err := p.Init(PChannelConfig{
		PChannelID:    "",
		MaxMsgCount:   5,
		MaxCacheCount: 10,
		DiskPath:      diskPath,
	})

	suite.assert.NotNil(err)
}

func (suite *pchannelTestSuite) TestCreatePChannelWrongPath() {
	p := &PChannel{}
	err := p.Init(PChannelConfig{
		PChannelID:    "",
		MaxMsgCount:   5,
		MaxCacheCount: 10,
		DiskPath:      "D:\\",
	})

	suite.assert.NotNil(err)
}

func TestAlphaSequence(t *testing.T) {
	suite.Run(t, new(pchannelTestSuite))
}
