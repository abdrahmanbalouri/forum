package config

import (
	"fmt"
	"regexp"
)

var (
	EmailExp    *regexp.Regexp
	UsernameExp *regexp.Regexp
)

func InitRegex() {
	var err error
	EmailExp, err = regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]+`)
	if err != nil {
		fmt.Println(err)
		return
	}

	UsernameExp, err = regexp.Compile(`^[a-zA-Z]+[a-zA-Z0-9._]+`)
	if err != nil {
		fmt.Println(err)
		return
	}
}
