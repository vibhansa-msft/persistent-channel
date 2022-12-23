package pchannel

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/vibhansa-msft/pchannel/alpha_sequence"
)

const (
	IndexFileName = "index.pch"
)

type PMessage[T any] struct {
	id   string
	data T
}

type Serialize[T any] func(T) []byte
type Deserialize[T any] func([]byte) T

// PChannel represents a persistent channel
type PChannel[T any] struct {
	PChannelConfig                // Configuration for this channel
	persistPath    string         // Path to directory of this channel
	serialize      Serialize[T]   // Method to serialize the data so that it can be dumped to disk
	deserialize    Deserialize[T] // Method to deserialize the data read from disk

	mtx sync.Mutex // Mutex to make shared variables sfae

	alphaSeq      *alpha_sequence.AlphaSequence // Alphabet sequence of all items pushed so far
	alphaSeqCache *alpha_sequence.AlphaSequence // Alphabet sequence of all items in cache
	cacheCount    uint64                        // Number of messages currently in channel

	channel chan (PMessage[T]) // Channel that will hold the data in memory
}

// Init : Initializes the channel
func (pc *PChannel[T]) Init(c PChannelConfig, s Serialize[T], d Deserialize[T]) error {
	pc.PChannelConfig = c

	// Validate the config
	err := validateConfig(&pc.PChannelConfig)
	if err != nil {
		return err
	}

	if s == nil || d == nil {
		return fmt.Errorf("Serialize and Deserialize methods can not be null")
	}

	pc.serialize = s
	pc.deserialize = d

	pc.persistPath = filepath.Join(pc.DiskPath, pc.PChannelID)

	pc.alphaSeq, err = alpha_sequence.CreateAlphaSequence(5)
	if err != nil {
		return err
	}

	pc.alphaSeqCache, err = alpha_sequence.CreateAlphaSequence(5)
	if err != nil {
		return err
	}

	pc.channel = make(chan PMessage[T], pc.MaxCacheCount)
	pc.cacheCount = 0

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
func (pc *PChannel[T]) Destroy() error {
	close(pc.channel)

	// Delete all persisted message at this path and remove the uid directory as well
	err := os.RemoveAll(pc.persistPath)
	if err != nil {
		return err
	}

	return nil
}

// restoreChannel : Restore previously persisted messages for this channel
func (pc *PChannel[T]) restoreChannel() error {
	// Get list of files from persist path
	files, err := os.ReadDir(pc.persistPath)
	if err != nil {
		return err
	}

	sortFileList(files)

	if len(files) != 0 {
		// There are items to be recovered

		pc.alphaSeq.SetString(path.Base(files[len(files)-1].Name()))

		for _, file := range files {
			if pc.cacheCount >= pc.MaxCacheCount {
				break
			}

			err := pc.readAndQueue(path.Base(file.Name()), filepath.Join(pc.persistPath, file.Name()))
			if err != nil {
				return err
			}
			pc.alphaSeqCache.SetString(path.Base(file.Name()))
		}
	}

	return nil
}

func (pc *PChannel[T]) PutMessage(data T) error {
	id := pc.alphaSeq.Next()
	fname := filepath.Join(pc.persistPath, id)

	err := os.WriteFile(fname, pc.serialize(data), 0447)
	if err != nil {
		return err
	}
	pc.mtx.Lock()
	defer pc.mtx.Unlock()

	if pc.cacheCount < pc.MaxCacheCount {
		// This message goes to cache as well
		msg := PMessage[T]{
			data: data,
			id:   id,
		}

		pc.channel <- msg
		pc.cacheCount++
		_ = pc.alphaSeqCache.Next()
	}

	return nil
}

func (pc *PChannel[T]) GetMessage() (T, string) {
	msg := <-pc.channel

	pc.mtx.Lock()
	defer pc.mtx.Unlock()
	pc.cacheCount--

	if pc.cacheCount <= pc.MaxCacheCount {
		// There is space of atleast one more message in cache so read from disk and fill cache
		if pc.alphaSeq.Get() > pc.alphaSeqCache.Get() {
			// There is atleast one message on disk that can be read
			id := pc.alphaSeqCache.Next()
			fname := filepath.Join(pc.persistPath, id)
			err := pc.readAndQueue(id, fname)
			if err != nil {
				var x T
				return x, err.Error()
			}
		}
	}

	return msg.data, msg.id
}

func (pc *PChannel[T]) ReleseMessage(id string) error {
	fname := filepath.Join(pc.persistPath, id)
	return os.Remove(fname)
}

func (pc *PChannel[T]) readAndQueue(id string, fname string) error {
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	msg := PMessage[T]{
		data: pc.deserialize(data),
		id:   id,
	}

	pc.channel <- msg
	pc.cacheCount++
	return nil
}
