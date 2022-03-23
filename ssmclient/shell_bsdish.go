//go:build darwin || netbsd || freebsd || openbsd || dragonfly
// +build darwin netbsd freebsd openbsd dragonfly

package ssmclient

import (
	"golang.org/x/sys/unix"
	"os"
)

const (
	ioctlReadTermios  uint = unix.TIOCGETA
	ioctlWriteTermios uint = unix.TIOCSETAF
)

// see also: https://pkg.go.dev/golang.org/x/term?utm_source=godoc#MakeRaw.
func configureStdin() (err error) {
	origTermios, err = unix.IoctlGetTermios(int(os.Stdin.Fd()), ioctlReadTermios)
	if err != nil {
		return err
	}

	// unsetting ISIG means that this process will no longer respond to the INT, QUIT, SUSP
	// signals (they go downstream to the instance session, which is desirable).  Which means
	// those signals are unavailable for shutting down this process
	newTermios := *origTermios
	newTermios.Lflag = origTermios.Lflag ^ unix.ICANON ^ unix.ECHO ^ unix.ISIG

	return unix.IoctlSetTermios(int(os.Stdin.Fd()), ioctlWriteTermios, &newTermios)
}
