package main

//ParseRequest - parse string to SMTPRequest
func ParseRequest(rawRequest string) *SMTPRequest {
	for _, command := range supportedCommands {
		if len(rawRequest) < len(command) {
			continue
		}
		if command == rawRequest[:len(command)] {
			return &SMTPRequest{
				Command: command,
				Body:    rawRequest[len(command):],
			}
		}
	}

	return &SMTPRequest{
		Command: "",
		Body:    rawRequest,
	}
}
