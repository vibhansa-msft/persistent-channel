# persistent-channel

## Persistent channel for Go

Channels in go are non-persistent. If the application terminates then all data pushed to channel is lost.
To perserve the order of the messages pushed to channel across re-runs this persistent-channel is created.
Each persistent-channel has an unique ID assigned on creation. There will be a directory with the same name
created on the given path. Inside this directory a file will be created for each message pushed in the channel.
Once application is done processing the message it can release the message which deletes the file from disk.
In any condition application decides to terminates, messages which are not released will remain on the disk.
Next time when application starts it can provide the same ID of the channel as input and persistent-channel
will try to recover the messages from disk.

## Configuration

- PChannelID : Unique ID assigned to this channel. Use same ID for recovery. Passing this as empty string will create a new ID and respecitve directory
- MaxMsgCount : Maximum number of messages that can be stored in this channel.
- MaxCacheCount : Maximum number of messages to be kept in memory for faster access. These messages will be persisted on disk as well for recover. Rest of the messages will remain only on disk.
- DiskPath : Directory persistent-channel shall use for its book-keeping. This directory shall exists otherwise initilization will fail.

## Example

- Look at pchannel_test.go to understand how to use it.