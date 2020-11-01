package main

const (
	smtpCommandHelo      = "HELO"
	smtpCommandHelp      = "HELP"
	smtpCommandVerify    = "VRFY"
	smtpCommandQuit      = "QUIT"
	smtpCommandNoop      = "NOOP"
	smtpCommandMailFrom  = "MAIL FROM"
	smtpCommandRecipient = "RCPT TO"
	smtpCommandData      = "DATA"
)

var supportedCommands = [...]string{
	smtpCommandHelo,
	smtpCommandHelp,
	smtpCommandVerify,
	smtpCommandQuit,
	smtpCommandNoop,
	smtpCommandMailFrom,
	smtpCommandRecipient,
	smtpCommandData,
}

//SMTPRequest - SMTP request data
type SMTPRequest struct {
	Command string
	Body    string
}
