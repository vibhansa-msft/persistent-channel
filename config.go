package pchannel

import (
	"fmt"
	"os"
	"sort"
)

// PChannelConfig defines configuration for PChannel
type PChannelConfig struct {
	PChannelID    string // Unique ID of this channel
	MaxMsgCount   uint64 // Maximum number of messages this channel can store
	MaxCacheCount uint64 // Maximum number of messages this channel can keep in memory, rest will be saved on disk
	DiskPath      string // Path where messages will be persisted
}

func validateConfig(c *PChannelConfig) error {
	if c.MaxMsgCount <= c.MaxCacheCount {
		return fmt.Errorf("MaxMsgCount shall be higher than MaxCacheCount")
	}

	err := dirExists(c.DiskPath)
	if err != nil {
		return err
	}

	return nil
}

func dirExists(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path %s does not exists", path)
	}

	if !info.IsDir() {
		return fmt.Errorf("path %s is not a directory", path)
	}

	return nil
}

func sortFileList(files []os.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		f1info, _ := files[i].Info()
		f2info, _ := files[j].Info()

		if f1info.ModTime() == f2info.ModTime() {
			return files[i].Name() < files[j].Name()
		} else {
			return f1info.ModTime().Unix() < f2info.ModTime().Unix()
		}
	})
}
