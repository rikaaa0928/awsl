package global

type constString string

const (
	//CTXReceiveAuth constString = "receiveAuth"
	//CTXSendAuth    constString = "sendAuth"
	CTXOutTag  constString = "outTag"
	CTXOutType constString = "outType"
	CTXInTag   constString = "inTag"
	CTXInType  constString = "inType"
	//CTXSendData    constString = "sendData"
	CTXSuperType constString = "superType"
	CTXSuperData constString = "superData"

	TransferAuth   = "auth"
	TransferAddr   = "addr"
	TransferSupper = "super"
)
