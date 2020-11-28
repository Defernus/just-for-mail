package main

const (
	//SMTP commands
	smtpCommandHelo      = "HELO"
	smtpCommandEhlo      = "EHLO"
	smtpCommandHelp      = "HELP"
	smtpCommandVerify    = "VRFY"
	smtpCommandQuit      = "QUIT"
	smtpCommandNoop      = "NOOP"
	smtpCommandMailFrom  = "MAIL FROM"
	smtpCommandRecipient = "RCPT TO"
	smtpCommandData      = "DATA"

	//ESMTP commands
	esmtpCommandStartTLS = "STARTTLS" //not yet implemented
	esmtpCommandAuth     = "AUTH" //not yet implemented
	esmtpCommandBdat     = "BDAT" //not yet implemented
)

var supportedCommands = [...]string{
	smtpCommandHelo,
	smtpCommandEhlo,
	smtpCommandHelp,
	smtpCommandVerify,
	smtpCommandQuit,
	smtpCommandNoop,
	smtpCommandMailFrom,
	smtpCommandRecipient,
	smtpCommandData,
	esmtpCommandStartTLS,
	esmtpCommandAuth,
}

//SMTPRequest - SMTP request data
type SMTPRequest struct {
	Command string
	Body    string
}
