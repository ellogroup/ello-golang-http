package validate

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/url"
	"testing"
)

type validateMock struct {
	mock.Mock
}

func (m *validateMock) Var(field interface{}, tag string) error {
	args := m.Called(field, tag)
	return args.Error(0)
}

func TestQueryValidator_Validate(t *testing.T) {

	type validateResponse struct {
		fieldInput interface{}
		tagInput   string
		errOut     error
	}

	type args struct {
		q     url.Values
		rules map[string]string
	}

	errValidation := errors.New("validation error")

	tests := []struct {
		name              string
		validateResponses []validateResponse
		args              args
		want              map[string]error
	}{
		{
			name: "One valid value, returns no errors",
			validateResponses: []validateResponse{
				{"test-value", "test-rule", nil},
			},
			args: args{
				q:     map[string][]string{"test-param": {"test-value"}},
				rules: map[string]string{"test-param": "test-rule"},
			},
			want: map[string]error{},
		},
		{
			name: "Two valid values, returns no errors",
			validateResponses: []validateResponse{
				{"test-value-1", "test-rule-1", nil},
				{"test-value-2", "test-rule-2", nil},
			},
			args: args{
				q:     map[string][]string{"test-param-1": {"test-value-1"}, "test-param-2": {"test-value-2"}},
				rules: map[string]string{"test-param-1": "test-rule-1", "test-param-2": "test-rule-2"},
			},
			want: map[string]error{},
		},
		{
			name: "One valid value, one invalid value, returns error for invalid value",
			validateResponses: []validateResponse{
				{"test-value-1", "test-rule-1", nil},
				{"test-value-2", "test-rule-2", errValidation},
			},
			args: args{
				q:     map[string][]string{"test-param-1": {"test-value-1"}, "test-param-2": {"test-value-2"}},
				rules: map[string]string{"test-param-1": "test-rule-1", "test-param-2": "test-rule-2"},
			},
			want: map[string]error{"test-param-2": errValidation},
		},
		{
			name: "One valid value, one missing value, returns error for missing value",
			validateResponses: []validateResponse{
				{"test-value-1", "test-rule-1", nil},
				{"", "test-rule-2", errValidation},
			},
			args: args{
				q:     map[string][]string{"test-param-1": {"test-value-1"}},
				rules: map[string]string{"test-param-1": "test-rule-1", "test-param-2": "test-rule-2"},
			},
			want: map[string]error{"test-param-2": errValidation},
		},
		{
			name: "Two missing values, returns error for both missing values",
			validateResponses: []validateResponse{
				{"", "test-rule-1", errValidation},
				{"", "test-rule-2", errValidation},
			},
			args: args{
				q:     map[string][]string{},
				rules: map[string]string{"test-param-1": "test-rule-1", "test-param-2": "test-rule-2"},
			},
			want: map[string]error{"test-param-1": errValidation, "test-param-2": errValidation},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := new(validateMock)
			for _, v := range tt.validateResponses {
				m.On("Var", v.fieldInput, v.tagInput).Return(v.errOut)
			}

			v := QueryValidator{
				validate: m,
			}
			assert.Equalf(t, tt.want, v.Validate(tt.args.q, tt.args.rules), "Validate(%v, %v)", tt.args.q, tt.args.rules)

			m.AssertExpectations(t)
		})
	}
}
