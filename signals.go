package gmicro

import (
	"os"
	"syscall"
)

// InterruptSignals interrupt signals.
var InterruptSignals = []os.Signal{
	syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP,
	syscall.SIGSTOP, syscall.SIGQUIT,
}
