package query

import (
	"github.com/go-playground/validator"
	"net/url"
)

type Var interface {
	Var(field interface{}, tag string) error
}

type Validator struct {
	validate Var
}

func NewValidator() *Validator {
	return &Validator{validate: validator.New()}
}

func (v Validator) Validate(q url.Values, rules map[string]string) map[string]error {
	e := map[string]error{}
	for param, rule := range rules {
		if err := v.validate.Var(q.Get(param), rule); err != nil {
			e[param] = err
		}
	}
	return e
}
