package public

// CmdCode
var (
	CmdUse     = uint8(0)
	CmdSend    = uint8(1)
	CmdWatch   = uint8(2)
	CmdIgnore  = uint8(3)
	CmdReserve = uint8(4)
)

// ResCode
var (
	ResSuccess   = uint8(0)
	ResSystemErr = uint8(1)
	ResUnknowCmd = uint8(2)

	ResFail = uint8(255)
)
