package command

type CmdError struct {
	Stderr []byte
	Err    error
}
