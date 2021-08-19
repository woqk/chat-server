package hub

import "time"

var (
	writeWait = 10 * time.Second
)

const (
	TypeError         = "error"
	TypeLogin         = "login"
	TypePostMessage   = "postMessage"
	TypeMessage       = "message"
	TypeChannelJoin   = "channel.join"
	TypeChannelJoined = "channel.joined"
	TypeChannelLeave  = "channel.leave"
	TypeChannelLeaved = "channel.leaved"
	TypeChannelList   = "channel.list"
	TypeChannelResult = "channel.result"
	TypeFriendList    = "friend.list"
	TypeFriendResult  = "friend.result"
)
