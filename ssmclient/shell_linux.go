//go:build linux
// +build linux

package ssmclient

import (
	"golang.org/x/sys/unix"
	"os"
)

const (
	ioctlReadTermios  uint = unix.TCGETS
	ioctlWriteTermios uint = unix.TCSETSF
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
	newTermios.Iflag = origTermios.Iflag | unix.IUTF8
	newTermios.Lflag = origTermios.Lflag ^ unix.ICANON ^ unix.ECHO ^ unix.ISIG

	return unix.IoctlSetTermios(int(os.Stdin.Fd()), ioctlWriteTermios, &newTermios)
}
