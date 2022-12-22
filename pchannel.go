package pchannel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/vibhansa-msft/pchannel/alpha_sequence"
)

const (
	IndexFileName = "index.pch"
)

// PChannel represents a persistent channel
type PChannel struct {
	PChannelConfig        // Configuration for this channel
	persistPath    string // Path to directory of this channel
	alphaSeq       *alpha_sequence.AlphaSequence
}

// Init : Initializes the channel
func (pc *PChannel) Init(c PChannelConfig) error {
	pc.PChannelConfig = c

	// Validate the config
	err := validateConfig(&pc.PChannelConfig)
	if err != nil {
		return err
	}

	pc.persistPath = filepath.Join(pc.DiskPath, pc.PChannelID)
	pc.alphaSeq, err = alpha_sequence.CreateAlphaSequence(5)
	if err != nil {
		return err
	}

	if pc.PChannelID == "" {
		// This is a new channel getting created so allocate an ID
		pc.PChannelID = uuid.NewString()
		pc.persistPath = filepath.Join(pc.DiskPath, pc.PChannelID)

		// Create a new directory with uid to persist the messages of this channel
		err := os.Mkdir(pc.persistPath, 0447)
		if err != nil {
			return err
		}

	} else {
		// This is PChannel restore case so load existing data from disk
		err := dirExists(pc.persistPath)
		if err != nil {
			return err
		}

		err = pc.restoreChannel()
		if err != nil {
			return err
		}
	}

	return nil
}

// Destroy : Destroy the channel and discard all persisted messages
func (pc *PChannel) Destroy() error {
	// Delete all persisted message at this path and remove the uid directory as well
	err := os.RemoveAll(pc.persistPath)
	if err != nil {
		return err
	}

	return nil
}

// restoreChannel : Restore previously persisted messages for this channel
func (pc *PChannel) restoreChannel() error {
	// Get list of files from persist path
	files, err := os.ReadDir(pc.persistPath)
	if err != nil {
		return err
	}

	sortFileList(files)

	lastFile := files[len(files)-1]
	s := strings.Split(lastFile.Name(), "_")
	index := s[len(s)-1]
	fmt.Println(index)

	// TODO : Push some of these files to the channel in memory

	return nil
}

