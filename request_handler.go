package main

import (
	"log"
)

const (
	connStageStart             = iota
	connStageHeloReseived      = iota
	connStageMailFromReseived  = iota
	connStageRecipientReseived = iota
	connStageDataReseived      = iota

	connTypeSMTP  = iota
	connTypeESMTP = iota

	requestActionNext     = iota
	requestActionData     = iota
	requestActionStartTLS = iota
	requestActionClose    = iota

	dataActionNext    = iota
	dataActionDataEnd = iota
)

//RequestHandler - SMTPRequest handler
type RequestHandler struct {
	connType       int
	curentStage    int
	authorizedUser string
	sender         string
	recipients     []string
	data           string
}

//NewRequestHandler - create new RequestHandler
func NewRequestHandler() *RequestHandler {
	return &RequestHandler{
		connType:       connTypeSMTP,
		curentStage:    connStageStart,
		sender:         "",
		authorizedUser: "",
		recipients:     []string{},
	}
}

//HandleRequest - handle SMTPRequest and return response
func (handler *RequestHandler) HandleRequest(request *SMTPRequest) (string, int) {
	log.Printf("Command: %s, Body: %s\n", request.Command, request.Body)
	switch request.Command {
	case smtpCommandHelo:
		return handler.handleHelo(request)
	case smtpCommandEhlo:
		return handler.handleEhlo(request)
	case smtpCommandHelp:
		return handler.handleHelp(request)
	case smtpCommandVerify:
		return handler.handleVerify(request)
	case smtpCommandQuit:
		return handler.handleQuit(request)
	case smtpCommandNoop:
		return handler.handleNoop(request)
	case smtpCommandMailFrom:
		return handler.handleMailFrom(request)
	case smtpCommandRecipient:
		return handler.handleRecipient(request)
	case smtpCommandData:
		return handler.handleData(request)
	default:
		return "502", requestActionNext
	}
}

func (handler *RequestHandler) handleHelo(request *SMTPRequest) (string, int) {
	if handler.curentStage == connStageStart {
		handler.curentStage = connStageHeloReseived
	}
	return "250 OK", requestActionNext
}

func (handler *RequestHandler) handleEhlo(request *SMTPRequest) (string, int) {
	if handler.curentStage == connStageStart {
		handler.connType = connTypeESMTP
		handler.curentStage = connStageHeloReseived
	}
	return "250 OK", requestActionNext
}

func (handler *RequestHandler) handleHelp(request *SMTPRequest) (string, int) {
	return "214 https://tools.ietf.org/html/rfc5321", requestActionNext
}
func (handler *RequestHandler) handleVerify(request *SMTPRequest) (string, int) {
	return "252 just try to send mail", requestActionNext
}
func (handler *RequestHandler) handleQuit(request *SMTPRequest) (string, int) {
	return "221 see you", requestActionClose
}

func (handler *RequestHandler) handleNoop(request *SMTPRequest) (string, int) {
	return "250 OK", requestActionNext
}

func (handler *RequestHandler) handleMailFrom(request *SMTPRequest) (string, int) {
	if handler.curentStage == connStageStart {
		return "503 send HELO/EHLO first", requestActionNext
	}

	email := getEmailFromRequestBody(request.Body)
	if email == "" {
		return "555 wrong email format", requestActionNext
	}

	handler.curentStage = connStageMailFromReseived
	handler.sender = email
	handler.recipients = []string{}
	handler.data = ""

	return "250 OK", requestActionNext
}

func (handler *RequestHandler) handleRecipient(request *SMTPRequest) (string, int) {
	if handler.curentStage != connStageMailFromReseived {
		return "503", requestActionNext
	}

	email := getEmailFromRequestBody(request.Body)
	if email == "" {
		return "555 wrong email format", requestActionNext
	}

	handler.curentStage = connStageRecipientReseived
	handler.recipients = append(handler.recipients, email)

	return "250 OK", requestActionNext
}

func (handler *RequestHandler) handleData(request *SMTPRequest) (string, int) {
	if handler.curentStage != connStageRecipientReseived {
		return "503 send RCPT first", requestActionNext
	}

	handler.curentStage = connStageDataReseived
	return "354", requestActionData
}

//HandleData - handle data
func (handler *RequestHandler) HandleData(data string) (string, int) {
	if len(handler.data) != 0 {
		handler.data += "\n"
	}
	handler.data += data

	if data[len(data)-2:] == "\n." {
		for _, r := range handler.recipients {
			sendMessage(r, handler.sender, handler.data)
		}
		return "250 OK", dataActionDataEnd
	}
	return "", dataActionNext
}
