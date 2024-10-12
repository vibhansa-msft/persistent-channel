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

const diskPath = "./temp"

func (suite *pchannelTestSuite) TestCreatePChannelSuccess() {
	p := &PChannel[string]{}
	err := p.Init(PChannelConfig{
		PChannelID:    "",
		MaxMsgCount:   10,
		MaxCacheCount: 5,
		DiskPath:      diskPath,
	}, func(s string) []byte { return []byte(s) },
		func(s []byte) string { return string(s) },
	)

	suite.assert.Nil(err)

	err = p.Destroy()
	suite.assert.Nil(err)
}

func (suite *pchannelTestSuite) TestCreatePChannelWithNewDiskath() {
	p := &PChannel[string]{}
	err := p.Init(PChannelConfig{
		PChannelID:    "",
		MaxMsgCount:   10,
		MaxCacheCount: 5,
		DiskPath:      diskPath + "/new",
	}, func(s string) []byte { return []byte(s) },
		func(s []byte) string { return string(s) },
	)

	suite.assert.Nil(err)

	err = p.Destroy()
	suite.assert.Nil(err)
}

func (suite *pchannelTestSuite) TestCreatePChannelWrongSize() {
	p := &PChannel[string]{}
	err := p.Init(PChannelConfig{
		PChannelID:    "",
		MaxMsgCount:   5,
		MaxCacheCount: 10,
		DiskPath:      diskPath,
	}, func(s string) []byte { return []byte(s) },
		func(s []byte) string { return string(s) },
	)

	suite.assert.NotNil(err)
}

func (suite *pchannelTestSuite) TestCreatePChannelWrongPath() {
	p := &PChannel[string]{}
	err := p.Init(PChannelConfig{
		PChannelID:    "",
		MaxMsgCount:   5,
		MaxCacheCount: 10,
		DiskPath:      diskPath,
	}, func(s string) []byte { return []byte(s) },
		func(s []byte) string { return string(s) },
	)

	suite.assert.NotNil(err)
}

func (suite *pchannelTestSuite) TestCreatePChannelPutAndGet() {
	p := &PChannel[string]{}
	err := p.Init(PChannelConfig{
		PChannelID:    "",
		MaxMsgCount:   10,
		MaxCacheCount: 8,
		DiskPath:      diskPath,
		IdLen:         5,
	}, func(s string) []byte { return []byte(s) },
		func(s []byte) string { return string(s) },
	)
	suite.assert.Nil(err)

	err = p.PutMessage("123")
	suite.assert.Nil(err)

	data, id := p.GetMessage()
	suite.assert.Equal(id, "aaaaa")
	suite.assert.Equal(data, "123")

	err = p.ReleseMessage(id)
	suite.assert.Nil(err)

	err = p.Destroy()
	suite.assert.Nil(err)
}

func (suite *pchannelTestSuite) TestOnlyPersistMessage() {
	p := &PChannel[string]{}
	err := p.Init(PChannelConfig{
		PChannelID:    "",
		MaxMsgCount:   5,
		MaxCacheCount: 1,
		DiskPath:      diskPath,
		IdLen:         5,
	}, func(s string) []byte { return []byte(s) },
		func(s []byte) string { return string(s) },
	)
	suite.assert.Nil(err)

	err = p.PutMessage("123")
	suite.assert.Nil(err)

	err = p.PutMessage("456")
	suite.assert.Nil(err)

	err = p.PutMessage("789")
	suite.assert.Nil(err)

	data, id := p.GetMessage()
	suite.assert.Equal(id, "aaaaa")
	suite.assert.Equal(data, "123")
	err = p.ReleseMessage(id)
	suite.assert.Nil(err)

	data, id = p.GetMessage()
	suite.assert.Equal(id, "aaaab")
	suite.assert.Equal(data, "456")
	err = p.ReleseMessage(id)
	suite.assert.Nil(err)

	err = p.Destroy()
	suite.assert.Nil(err)
}

func (suite *pchannelTestSuite) TestRestoreChannel() {
	p := &PChannel[string]{}
	err := p.Init(PChannelConfig{
		PChannelID:    "",
		MaxMsgCount:   5,
		MaxCacheCount: 1,
		DiskPath:      diskPath,
		IdLen:         5,
	}, func(s string) []byte { return []byte(s) },
		func(s []byte) string { return string(s) },
	)
	suite.assert.Nil(err)

	err = p.PutMessage("123")
	suite.assert.Nil(err)

	err = p.PutMessage("456")
	suite.assert.Nil(err)

	err = p.PutMessage("789")
	suite.assert.Nil(err)

	err = p.PutMessage("901")
	suite.assert.Nil(err)

	data, id := p.GetMessage()
	suite.assert.Equal(id, "aaaaa")
	suite.assert.Equal(data, "123")
	err = p.ReleseMessage(id)
	suite.assert.Nil(err)

	// Without closing the pchannel lets recreate the object
	chId := p.PChannelID
	p = &PChannel[string]{}
	err = p.Init(PChannelConfig{
		PChannelID:    chId,
		MaxMsgCount:   5,
		MaxCacheCount: 1,
		DiskPath:      diskPath,
		IdLen:         5,
	}, func(s string) []byte { return []byte(s) },
		func(s []byte) string { return string(s) },
	)

	data, id = p.GetMessage()
	suite.assert.Equal(id, "aaaab")
	suite.assert.Equal(data, "456")
	err = p.ReleseMessage(id)
	suite.assert.Nil(err)

	data, id = p.GetMessage()
	suite.assert.Equal(id, "aaaac")
	suite.assert.Equal(data, "789")
	err = p.ReleseMessage(id)
	suite.assert.Nil(err)

	err = p.Destroy()
	suite.assert.Nil(err)
}

func TestAlphaSequence(t *testing.T) {
	suite.Run(t, new(pchannelTestSuite))
}
