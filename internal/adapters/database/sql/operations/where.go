package operations

import "fmt"

type Where struct {
	Equal     *Equal
	Equals    *[]Equal
	In        *In
	IsDeleted *bool
}

type In struct {
	Column string
	Values []interface{}
}

type Equal struct {
	Column string
	Value  interface{}
}

type IWhere interface {
	Handle() (string, []interface{})
}

func (w *Where) Handle() (where string, args []interface{}) {
	if w != nil {
		if w.Equal != nil {
			where, args = w.Equal.handle()
		} else if w.Equals != nil {
			where, args = w.Equal.handleMany(*w.Equals)
		} else if w.In != nil {
			where, args = w.In.handle()
		}

		if w.IsDeleted == nil || !*w.IsDeleted {
			where = fmt.Sprintf("%s AND deletedAt IS NULL", where)
		} else {
			where = fmt.Sprintf("%s AND deletedAt IS NOT NULL", where)
		}

		return
	}
	return "", nil
}

func (e *Equal) handle() (str string, value []interface{}) {
	str = fmt.Sprintf("WHERE %s=$1", e.Column)
	value = []any{e.Value}
	return
}

func (e *Where) handleUpdateWhere(index int) (str string, value []interface{}) {
	if e != nil && e.Equals != nil {
		for _, where := range *e.Equals {
			if str == "" {
				str = fmt.Sprintf("WHERE %s=$%v", where.Column, index+1)
			} else {
				str = fmt.Sprintf("%s, %s=$%v", str, where.Column, index+1)
			}
			value = append(value, where.Value)
		}
	}
	return
}

func (e *Equal) handleMany(equals []Equal) (str string, values []interface{}) {
	for i, equal := range equals {
		if str == "" {
			str = fmt.Sprintf("WHERE %s=$1", e.Column)
		} else {
			str = fmt.Sprintf("%s AND %s=$%v", str, e.Column, i+1)
		}
		values = append(values, equal.Value)
	}
	return
}

func (in *In) handle() (str string, values []interface{}) {
	for i, value := range in.Values {
		if str == "" {
			str = fmt.Sprintf("WHERE %s IN ($1", in.Column)
		} else {
			str = fmt.Sprintf("%s, $%v", str, i+1)
		}
		values = append(values, value)
	}

	str = fmt.Sprintf("%s)", str)

	return
}
