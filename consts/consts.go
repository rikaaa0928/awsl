package consts

type constString string

const (
	CTXReceiveAuth constString = "receiveAuth"
	CTXSendAuth    constString = "sendAuth"
	CTXRoute       constString = "route"
	CTXInTag       constString = "inTag"
	//CTXSendData    constString = "sendData"
	CTXSuperType constString = "superType"
	CTXSuperData constString = "superData"

	TransferAuth   = "auth"
	TransferAddr   = "addr"
	TransferSupper = "super"
)
