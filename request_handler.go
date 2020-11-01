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

	requestActionNext  = iota
	requestActionData  = iota
	requestActionClose = iota

	dataActionNextLine = iota
	dataActionDataEnd  = iota
	dataActionClose    = iota
)

//RequestHandler - SMTPRequest handler
type RequestHandler struct {
	curentStage int
	sender      string
	recipients  []string
	data        string
}

//NewRequestHandler - create new RequestHandler
func NewRequestHandler() *RequestHandler {
	return &RequestHandler{
		curentStage: connStageStart,
		sender:      "",
		recipients:  []string{},
	}
}

//HandleRequest - handle SMTPRequest and return response
func (handler *RequestHandler) HandleRequest(request *SMTPRequest) (string, int) {
	switch request.Command {
	case smtpCommandHelo:
		return handler.handleHelo(request)
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
		return "503 send HELO first", requestActionNext
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

//HandleData - handle data line by line
func (handler *RequestHandler) HandleData(data string) (string, int) {
	handler.data = data
	log.Println(data)
	return "250 OK", dataActionDataEnd
}
