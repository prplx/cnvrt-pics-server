package validator

type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.addError(key, message)
	}
}

func (v *Validator) addError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}
