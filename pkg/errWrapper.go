package myerr

import (
	"fmt"
	"runtime"
)

type err struct {
	msg string
}

func New(e string) err {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return err{msg: fmt.Sprintf("%s >> %s", details.Name(), e)}
	} else {
		return err{}
	}
}

func (thisErr err) Error() string {
	return thisErr.msg
}
