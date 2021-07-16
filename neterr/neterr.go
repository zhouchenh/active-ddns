package neterr

import (
	"errors"
	"github.com/zhouchenh/active-ddns/logger"
	"io"
	"net"
	"syscall"
)

func LogError(err error) {
	isDisabled := isRepeatedLog(err.Error()) || isNetOpError(err, "use of closed network connection")
	isDebug := isDisabled || errors.Is(err, io.EOF)
	isWarning := isDebug || errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET)
	if isDisabled {
		return
	} else if isDebug {
		logger.Debug().Msg(keepLog(err.Error()))
	} else if isWarning {
		logger.Warning().Msg(keepLog(err.Error()))
	} else {
		logger.Error().Msg(keepLog(err.Error()))
	}
}

func isNetOpError(err error, target string) bool {
	opErr, ok := err.(*net.OpError)
	if !ok {
		return false
	}
	return opErr.Err.Error() == target
}

var lastLogContent string

func keepLog(logContent string) string {
	lastLogContent = logContent
	return logContent
}

func isRepeatedLog(logContent string) bool {
	return lastLogContent == logContent
}
