package main

const (
	TagSys      = 1
	TagMessage  = 2
	TagFileName = 3
	TagFile     = 4
	TagTime     = 5
	TagName     = 6
)

const (
	// request:  should contain TagName with name
	// response: TagSys [SysLoginResponse, LoginStatus]
	SysLoginRequest   = 1
	LOGIN_OK          = 1
	NAME_USED         = 2
	NAME_WRONG_FORMAT = 3

	//request: 	None
	//response: TagSys [SysUserLoginNotiffication, type], TagName and TagTime
	SysUserLoginNotiffication = 3
	USER_CONNECTED            = 1
	USER_DISCONECTED          = 2

	//request: TagMessage TagText
	//response to sender: TagSys [SysMessage, message status], TagTime
	//response to others: TagSys [SysMessage], TagText, TagName, TagTime
	SysMessage           = 4
	MESSAGE_SENT         = 1
	MESSAGE_WRONG_FORMAT = 2
)
