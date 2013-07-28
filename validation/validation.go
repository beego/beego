package validation

import (
	"fmt"
	"reflect"
	"regexp"
)

type ValidationError struct {
	Message, Key string
}

// Returns the Message.
func (e *ValidationError) String() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// A Validation context manages data validation and error messages.
type Validation struct {
	Errors []*ValidationError
}

func (v *Validation) Clear() {
	v.Errors = []*ValidationError{}
}

func (v *Validation) HasErrors() bool {
	return len(v.Errors) > 0
}

// Return the errors mapped by key.
// If there are multiple validation errors associated with a single key, the
// first one "wins".  (Typically the first validation will be the more basic).
func (v *Validation) ErrorMap() map[string]*ValidationError {
	m := map[string]*ValidationError{}
	for _, e := range v.Errors {
		if _, ok := m[e.Key]; !ok {
			m[e.Key] = e
		}
	}
	return m
}

// Add an error to the validation context.
func (v *Validation) Error(message string, args ...interface{}) *ValidationResult {
	result := (&ValidationResult{
		Ok:    false,
		Error: &ValidationError{},
	}).Message(message, args...)
	v.Errors = append(v.Errors, result.Error)
	return result
}

// A ValidationResult is returned from every validation method.
// It provides an indication of success, and a pointer to the Error (if any).
type ValidationResult struct {
	Error *ValidationError
	Ok    bool
}

func (r *ValidationResult) Key(key string) *ValidationResult {
	if r.Error != nil {
		r.Error.Key = key
	}
	return r
}

func (r *ValidationResult) Message(message string, args ...interface{}) *ValidationResult {
	if r.Error != nil {
		if len(args) == 0 {
			r.Error.Message = message
		} else {
			r.Error.Message = fmt.Sprintf(message, args...)
		}
	}
	return r
}

// Test that the argument is non-nil and non-empty (if string or list)
func (v *Validation) Required(obj interface{}, key string) *ValidationResult {
	return v.apply(Required{key}, obj)
}

// Test that the obj is greater than min if obj's type is int
func (v *Validation) Min(obj interface{}, min int, key string) *ValidationResult {
	return v.apply(Min{min, key}, obj)
}

// Test that the obj is less than max if obj's type is int
func (v *Validation) Max(obj interface{}, max int, key string) *ValidationResult {
	return v.apply(Max{max, key}, obj)
}

// Test that the obj is between mni and max if obj's type is int
func (v *Validation) Range(obj interface{}, min, max int, key string) *ValidationResult {
	return v.apply(Range{Min{Min: min}, Max{Max: max}, key}, obj)
}

func (v *Validation) MinSize(obj interface{}, min int, key string) *ValidationResult {
	return v.apply(MinSize{min, key}, obj)
}

func (v *Validation) MaxSize(obj interface{}, max int, key string) *ValidationResult {
	return v.apply(MaxSize{max, key}, obj)
}

func (v *Validation) Length(obj interface{}, n int, key string) *ValidationResult {
	return v.apply(Length{n, key}, obj)
}

func (v *Validation) Alpha(obj interface{}, key string) *ValidationResult {
	return v.apply(Alpha{key}, obj)
}

func (v *Validation) Numeric(obj interface{}, key string) *ValidationResult {
	return v.apply(Numeric{key}, obj)
}

func (v *Validation) AlphaNumeric(obj interface{}, key string) *ValidationResult {
	return v.apply(AlphaNumeric{key}, obj)
}

func (v *Validation) Match(obj interface{}, regex *regexp.Regexp, key string) *ValidationResult {
	return v.apply(Match{regex, key}, obj)
}

func (v *Validation) NoMatch(obj interface{}, regex *regexp.Regexp, key string) *ValidationResult {
	return v.apply(NoMatch{Match{Regexp: regex}, key}, obj)
}

func (v *Validation) AlphaDash(obj interface{}, key string) *ValidationResult {
	return v.apply(AlphaDash{NoMatch{Match: Match{Regexp: alphaDashPattern}}, key}, obj)
}

func (v *Validation) Email(obj interface{}, key string) *ValidationResult {
	return v.apply(Email{Match{Regexp: emailPattern}, key}, obj)
}

func (v *Validation) IP(obj interface{}, key string) *ValidationResult {
	return v.apply(IP{Match{Regexp: ipPattern}, key}, obj)
}

func (v *Validation) Base64(obj interface{}, key string) *ValidationResult {
	return v.apply(Base64{Match{Regexp: base64Pattern}, key}, obj)
}

func (v *Validation) Mobile(obj interface{}, key string) *ValidationResult {
	return v.apply(Mobile{Match{Regexp: mobilePattern}, key}, obj)
}

func (v *Validation) Tel(obj interface{}, key string) *ValidationResult {
	return v.apply(Tel{Match{Regexp: telPattern}, key}, obj)
}

func (v *Validation) Phone(obj interface{}, key string) *ValidationResult {
	return v.apply(Phone{Mobile{Match: Match{Regexp: mobilePattern}},
		Tel{Match: Match{Regexp: telPattern}}, key}, obj)
}

func (v *Validation) ZipCode(obj interface{}, key string) *ValidationResult {
	return v.apply(ZipCode{Match{Regexp: zipCodePattern}, key}, obj)
}

func (v *Validation) apply(chk Validator, obj interface{}) *ValidationResult {
	if chk.IsSatisfied(obj) {
		return &ValidationResult{Ok: true}
	}

	// Add the error to the validation context.
	err := &ValidationError{
		Message: chk.DefaultMessage(),
		Key:     chk.GetKey(),
	}
	v.Errors = append(v.Errors, err)

	// Also return it in the result.
	return &ValidationResult{
		Ok:    false,
		Error: err,
	}
}

// Apply a group of validators to a field, in order, and return the
// ValidationResult from the first one that fails, or the last one that
// succeeds.
func (v *Validation) Check(obj interface{}, checks ...Validator) *ValidationResult {
	var result *ValidationResult
	for _, check := range checks {
		result = v.apply(check, obj)
		if !result.Ok {
			return result
		}
	}
	return result
}

// the obj parameter must be a struct or a struct pointer
func (v *Validation) Valid(obj interface{}) (b bool, err error) {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	switch {
	case isStruct(objT):
	case isStructPtr(objT):
		objT = objT.Elem()
		objV = objV.Elem()
	default:
		err = fmt.Errorf("%v must be a struct or a struct pointer", obj)
		return
	}

	for i := 0; i < objT.NumField(); i++ {
		var vfs []ValidFunc
		if vfs, err = getValidFuncs(objT.Field(i)); err != nil {
			return
		}
		for _, vf := range vfs {
			if _, err = funcs.Call(vf.Name,
				mergeParam(v, objV.Field(i).Interface(), vf.Params)...); err != nil {
				return
			}
		}
	}
	return !v.HasErrors(), nil
}
