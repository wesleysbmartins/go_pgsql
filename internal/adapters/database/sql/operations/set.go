package operations

import "fmt"

type Set struct {
	Values []SetParam
}

type SetParam struct {
	Column string
	Value  interface{}
}

type ISet interface {
	handle() (str string, values []interface{})
}

func (s *Set) handle() (str string, values []interface{}) {
	for i, set := range s.Values {
		if str == "" {
			str = fmt.Sprintf("SET %s=$1", set.Column)
		} else {
			str = fmt.Sprintf("%s, %s=$%v", str, set.Column, i+1)
		}
		values = append(values, set.Value)
	}
	return
}
