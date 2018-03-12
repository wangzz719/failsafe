package failsafe

import (
	"reflect"
	"strings"
)

type DstFunc func() (interface{}, error)

type FailSafe struct {
}

func (failSafe *FailSafe) Safe(failRtn interface{}, open bool, errors []string, dstFunc DstFunc, rtn interface{}) error {
	rtnV := reflect.Indirect(reflect.ValueOf(rtn))
	rtnPtrV := reflect.ValueOf(rtn)
	if rtnPtrV.Kind() != reflect.Ptr {
		panic("failsafe: return must be a pointer")
	}

	failRtnV := reflect.ValueOf(failRtn)

	destFuncResult, err := dstFunc()
	destFuncResultV := reflect.ValueOf(destFuncResult)

	if !destFuncResultV.IsValid() || destFuncResultV.Kind() == reflect.Ptr && destFuncResultV.IsNil() {
		return nil
	}

	if !open {
		if err != nil {
			return err
		}
		rtnV.Set(destFuncResultV)
		return nil
	} else {
		if err != nil {
			errName := reflect.TypeOf(err).String()
			errName = formatError(errName)
			if isFailSafeErrors(errors, errName) {
				rtnV.Set(failRtnV)
				return nil
			}
			return err
		}

		left := recursiveIndirectType(rtnV.Type())
		right := recursiveIndirectType(destFuncResultV.Type())
		if left.Kind() != right.Kind() {
			rtnV.Set(failRtnV)
			return nil
		}
		rtnV.Set(destFuncResultV)
	}
	return nil
}

func formatError(err string) string {
	e := strings.Replace(err, "*", "", -1)
	eSplitN := strings.Split(e, ".")
	return eSplitN[len(eSplitN)-1]
}

func isFailSafeErrors(errors []string, err string) bool {
	for _, e := range errors {
		if e == err {
			return true
		}
	}
	return false
}
