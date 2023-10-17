package validate

import (
	"net/url"
)

type Var interface {
	Var(field interface{}, tag string) error
}

type QueryValidator struct {
	validate Var
}

func NewQueryValidator(v Var) *QueryValidator {
	return &QueryValidator{validate: v}
}

func (v QueryValidator) Validate(q url.Values, rules map[string]string) map[string]error {
	e := map[string]error{}
	for param, rule := range rules {
		if err := v.validate.Var(q.Get(param), rule); err != nil {
			e[param] = err
		}
	}
	return e
}
