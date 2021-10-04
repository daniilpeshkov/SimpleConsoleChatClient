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

	//global OK return value in responses
	OK = 1

	// request:  should contain TagName with name
	// response: TagSys [SysLoginResponse, LoginStatus]
	SysLoginRequest   = 1
	LOGIN_OK          = OK
	NAME_USED         = 2
	NAME_WRONG_FORMAT = 3

	//request: 	None
	//response: TagSys [SysUserLoginNotiffication, type], TagName and TagTime
	SysUserLoginNotiffication = 3
	USER_CONNECTED            = OK
	USER_DISCONECTED          = 2

	//request: TagMessage TagText
	//response to sender: TagSys [SysMessage, message status], TagTime
	//response to others: TagSys [SysMessage], TagText, TagName, TagTime
	SysMessage           = 4
	MESSAGE_SENT         = OK
	MESSAGE_WRONG_FORMAT = 2

	//request TagFile TagFileName
	SysFile   = 5
	FILE_SENT = OK
	ERR       = 2
)
