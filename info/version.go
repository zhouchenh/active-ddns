package info

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	name    = "active-ddns"
	version = "0.0.1"
	build   = "Beta "
	intro   = "A simple DDNS tool"
)

func Name() string {
	return name
}

func Version() string {
	return version
}

func VersionStatement() []string {
	return []string{
		concatenate(Name(), " ", Version(), " ", build, "(", runtime.GOOS, "/", runtime.GOARCH, ")"),
		intro,
	}
}

func concatenate(a ...interface{}) string {
	builder := strings.Builder{}
	for _, value := range a {
		builder.WriteString(fmt.Sprint(value))
	}
	return builder.String()
}
