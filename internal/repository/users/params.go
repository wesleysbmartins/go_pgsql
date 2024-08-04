package repository

import (
	"go_pgsql/internal/adapters/database/sql/operations"
)

type UserParams struct {
	Id        int
	Ids       []int
	Name      string
	Names     []string
	Username  string
	Usernames []string
	Email     string
	Emails    []string
	Password  string
	Token     string
	IsDeleted *bool
}

type IUserParams interface {
	handleWhere() *operations.Where
	handleSet() operations.Set
	handleSliceInterface(values []any) []interface{}
}

func (p *UserParams) handleWhere() *operations.Where {
	where := operations.Where{}
	equals := []operations.Equal{}

	where.IsDeleted = p.IsDeleted

	if p.Id != 0 {
		where.Equal = &operations.Equal{
			Column: "id",
			Value:  p.Id,
		}

		return &where
	}

	if p.Name != "" {
		equals = append(equals, operations.Equal{
			Column: "name",
			Value:  p.Name,
		})
	}

	if p.Username != "" {
		equals = append(equals, operations.Equal{
			Column: "username",
			Value:  p.Username,
		})
	}

	if p.Email != "" {
		equals = append(equals, operations.Equal{
			Column: "email",
			Value:  p.Email,
		})
	}

	if len(p.Ids) > 0 {
		where.In = &operations.In{
			Column: "id",
			Values: p.handleIntSliceInterface(p.Ids),
		}

		return &where
	}

	if len(p.Names) > 0 {
		where.In = &operations.In{
			Column: "name",
			Values: p.handleStringSliceInterface(p.Names),
		}

		return &where
	}

	if len(p.Usernames) > 0 {
		where.In = &operations.In{
			Column: "username",
			Values: p.handleStringSliceInterface(p.Usernames),
		}

		return &where
	}

	if len(p.Emails) > 0 {
		where.In = &operations.In{
			Column: "email",
			Values: p.handleStringSliceInterface(p.Emails),
		}

		return &where
	}

	if len(equals) == 1 {
		where.Equal = &equals[0]
	} else if len(equals) > 1 {
		where.Equals = &equals
	}

	return &where
}

func (p *UserParams) handleSet() operations.Set {
	set := operations.Set{}

	if p.Name != "" {
		set.Values = append(set.Values, operations.SetParam{Column: "name", Value: p.Name})
	}

	if p.Username != "" {
		set.Values = append(set.Values, operations.SetParam{Column: "username", Value: p.Username})
	}

	if p.Email != "" {
		set.Values = append(set.Values, operations.SetParam{Column: "email", Value: p.Email})
	}

	if p.Password != "" {
		set.Values = append(set.Values, operations.SetParam{Column: "password", Value: p.Password})
	}

	if p.Token != "" {
		set.Values = append(set.Values, operations.SetParam{Column: "token", Value: p.Token})
	}

	return set
}

func (p *UserParams) handleIntSliceInterface(values []int) []interface{} {
	interfaceSlice := make([]interface{}, len(values))

	for i, v := range values {
		interfaceSlice[i] = v
	}

	return interfaceSlice
}

func (p *UserParams) handleStringSliceInterface(values []string) []interface{} {
	interfaceSlice := make([]interface{}, len(values))

	for i, v := range values {
		interfaceSlice[i] = v
	}

	return interfaceSlice
}
