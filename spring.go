package beego

import (
	"errors"
	"github.com/astaxie/beego/context"
	"reflect"
	"regexp"
	"strings"
)

const (
	TagName = "autowired"

	PreHandling  = "preHandling"
	PostHandling = "postHandling"

	InterfaceType      = "InterfaceType"
	ImplementationType = "ImplementationType"
	BeanName           = "BeanName"

	PanicOnError   = "panicOnError"
	NilOnError     = "nilOnError"
	SaveErrorToCtx = "saveErrorToCtx"
)

var (
	TypeMismatch         = errors.New("type mismatch")
	FactoryNotFound      = errors.New("factory not found")
	ErrorHandlerNotFound = errors.New("error handler not found")
	BeanNameNotSpecified = errors.New("bean name not specified")

	AspectType = aspectType()

	beanNamePattern     = regexp.MustCompile(`beanName\((?P<Value>.+)\)`)
	bindMethodPattern   = regexp.MustCompile(`bindMethod\((?P<Value>.+)\)`)
	errorHandlerPattern = regexp.MustCompile(`errorHandler\((?P<Value>.+)\)`)
	occasionPattern     = regexp.MustCompile(`occasion\((?P<Value>.+)\)`)

	kernel    = make(map[string]Factory)
	container = make(map[string]reflect.Type)

	panicOnError = new(panicIfError)
	nilOnError   = new(nilIfError)
	toCtxOnerror = new(errorToCtx)

	errorHandlers = map[string]ErrorResolver{
		PanicOnError:   panicOnError,
		NilOnError:     nilOnError,
		SaveErrorToCtx: toCtxOnerror,
	}

	ignoredTypes = []reflect.Type{reflect.TypeOf(new(context.Context))}
)

type Aspect interface {
}

type ErrorResolver interface {
	Handle(*context.Context, error) (breakOperation bool)
}

type Factory interface {
	New(*context.Context) (interface{}, error)
}

func Register(beanName string, factory Factory) {
	kernel[beanName] = factory
}

func Resolve(beanName string) Factory {
	return kernel[beanName]
}

func RegisterErrorHandler(name string, handler ErrorResolver) {
	errorHandlers[name] = handler
}

func RegisterIgnoreType(obj interface{}) {
	t := reflect.TypeOf(obj)
	for ; t.Kind() == reflect.Ptr; t = t.Elem() {
	}
	ignoredTypes = append(ignoredTypes, t)
}

func RegisterInterfaceTypeBinding(binding interface{}) {
	bv := reflect.ValueOf(binding)
	for ; bv.Kind() == reflect.Ptr; bv = bv.Elem() {
	}
	bt := bv.Type()
	fit, exist := bt.FieldByName(InterfaceType)
	if !exist {
		panic("interface type must be defined")
	}
	fii, exist := bt.FieldByName(ImplementationType)
	if !exist {
		panic("implementation type must be defined")
	}
	fiit := fii.Type
	for ; fiit.Kind() == reflect.Ptr; fiit = fiit.Elem() {
	}
	vbk := bv.FieldByName(BeanName)
	bindingKey := "global"
	if vbk.IsValid() && vbk.CanInterface() {
		bindingKey, _ = vbk.Interface().(string)
	}
	id := fit.Type.Name() + bindingKey
	container[id] = fiit
}

func injectPre(input *context.Context, target interface{}, bindingMethod string) {
	if target == nil {
		return
	}
	inject(input, reflect.ValueOf(target), bindingMethod, PreHandling, 0)
}

func injectPost(input *context.Context, target interface{}, bindingMethod string) {
	if target == nil {
		return
	}
	inject(input, reflect.ValueOf(target), bindingMethod, PostHandling, 0)
}

type dummyAspectFieldStruct struct {
	Dummy Aspect
}

func aspectType() reflect.Type {
	if df, exist := reflect.TypeOf(dummyAspectFieldStruct{}).FieldByName("Dummy"); exist {
		return df.Type
	}
	panic("filed of aspect not found")
}

func inject(input *context.Context, target reflect.Value, bindingMethod, occasion string, level int) {
	if level > 5 {
		return
	}
	type_, value, injectable := getUnderlayType(target)
	if !injectable || !value.IsValid() {
		return
	}
	count := type_.NumField()
	for i := 0; i < count; i++ {
		field := type_.Field(i)
		if ignored(field.Type) {
			continue
		}
		if breakOperation := injectOneField(&field, &value, i, input, bindingMethod, occasion, level); breakOperation {
			break
		}
		child := value.FieldByName(field.Name)
		inject(input, child, bindingMethod, occasion, level+1)
	}
}

func createAndSetDefaultValue(value *reflect.Value, fieldIndex int) (target reflect.Value) {
	fv := value.Field(fieldIndex)
	if ignored(fv.Type()) {
		return
	}
	fvk := fv.Kind()
	if fvk == reflect.Struct {
		target = fv
		return
	}
	if fvk == reflect.Interface {
		id := fv.Type().Name() + "global"
		if ct, exist := container[id]; exist {
			cv := reflect.New(ct)
			if fv.CanSet() && cv.IsValid() {
				fv.Set(cv)
				target = cv
				return
			}
		}
	}
	if !fv.CanSet() {
		return
	}
	ft := fv.Type()
	if ft.Kind() != reflect.Ptr || !fv.IsNil() {
		return
	}

	for ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
		fv = fv.Elem()
	}
	if !fv.CanSet() || ft.Kind() != reflect.Struct {
		return
	}
	target = reflect.New(ft).Elem()
	fv.Set(target)
	return
}

func ignored(t reflect.Type) bool {
	for ; t.Kind() == reflect.Ptr; t = t.Elem() {
	}
	for _, tt := range ignoredTypes {
		if t == tt {
			return true
		}
	}
	return false
}

func getUnderlayType(value reflect.Value) (underlayType reflect.Type, underlayValue reflect.Value, injectable bool) {
	if !value.IsValid() {
		return
	}
	t := value.Type()
	k := t.Kind()
	for {
		if k == reflect.Struct {
			underlayType = t
			underlayValue = value
			injectable = true
			return
		} else if k == reflect.Ptr {
			t = t.Elem()
			value = value.Elem()
			k = t.Kind()
		} else {
			break
		}
	}
	return
}

func getBeanConfig(tagValue, bindingMethod, occasion string) (beanName string, errorHandler ErrorResolver, notForThisMethod bool) {
	configs := strings.SplitN(tagValue, ";", -1)
	occ := PreHandling
	for _, config := range configs {
		c := strings.TrimSpace(config)
		if o := occasionPattern.FindStringSubmatch(c); o != nil {
			occ = strings.TrimSpace(o[1])
			if occ != occasion {
				notForThisMethod = true
				return
			}
		}
		if bm := bindMethodPattern.FindStringSubmatch(c); bm != nil {
			rawBM := bm[1]
			bindingMethods := strings.SplitN(rawBM, ",", -1)
			found := false
			for _, bmm := range bindingMethods {
				if strings.TrimSpace(bmm) == bindingMethod {
					found = true
					break
				}
			}
			if !found {
				notForThisMethod = true
				return
			}
		}
		if bn := beanNamePattern.FindStringSubmatch(c); bn != nil {
			beanName = bn[1]
			continue
		}
		if eh := errorHandlerPattern.FindStringSubmatch(c); eh != nil {
			ehName := eh[1]
			if ehh, exist := errorHandlers[ehName]; exist {
				errorHandler = ehh
				continue
			}
			panic(ErrorHandlerNotFound)
		}
	}
	if occ != occasion {
		notForThisMethod = true
		return
	}
	if errorHandler == nil {
		errorHandler = nilOnError
	}
	return
}

func injectOneField(field *reflect.StructField, value *reflect.Value, fieldIndex int, input *context.Context, bindingMethod, occasion string, level int) (breakOperation bool) {
	if !value.IsValid() {
		return
	}
	beanConfig, exist := field.Tag.Lookup(TagName)
	if !exist {
		if target := createAndSetDefaultValue(value, fieldIndex); target.IsValid() {
			inject(input, target, bindingMethod, occasion, level+1)
		}
		return
	}
	beanName, errorHandler, notForThisMethod := getBeanConfig(beanConfig, bindingMethod, occasion)
	if notForThisMethod {
		return
	}
	var instance interface{}
	var err error
	if beanName == "" {
		beanName = "global"
	}
	created := false
	factory, exist := kernel[beanName]
	if !exist {
		ft := field.Type
		if ft.Kind() == reflect.Interface {
			id := ft.Name() + beanName
			if fit, exist := container[id]; exist {
				if fiv := reflect.New(fit); fiv.CanInterface() {
					instance = fiv.Interface()
					created = true
				}
			} else {
				breakOperation = errorHandler.Handle(input, FactoryNotFound)
				return
			}
		} else {
			breakOperation = errorHandler.Handle(input, FactoryNotFound)
			return
		}
	} else {
		if instance, err = factory.New(input); err != nil {
			breakOperation = errorHandler.Handle(input, err)
			return
		}
	}
	if (beanName == "" || beanName == "global") && !created {
		panic(BeanNameNotSpecified)
	}
	situation := resolveSituation(reflect.TypeOf(instance), field.Type)
	if situation == "mismatch" {
		breakOperation = errorHandler.Handle(input, TypeMismatch)
		return
	}
	if situation == "aspect" {
		return
	}
	setValue(field, value, instance, situation)
	return
}

func setValue(field *reflect.StructField, target *reflect.Value, instance interface{}, situation string) {
	var finalTarget reflect.Value
	var finalValue reflect.Value
	if situation == "ok" {
		finalTarget = target.FieldByName(field.Name)
		finalValue = reflect.ValueOf(instance)
	} else if situation == "ptrToValue" {
		finalTarget = target.FieldByName(field.Name)
		finalValue = reflect.ValueOf(instance).Elem()
	} else if situation == "valueToPtr" {
		finalTarget = target.FieldByName(field.Name).Elem()
		finalValue = reflect.ValueOf(instance)
	}
	if finalTarget.CanSet() {
		finalTarget.Set(finalValue)
	}
}

func resolveSituation(gotType, expectType reflect.Type) string {
	if expectType == AspectType {
		return "aspect"
	}
	if gotType.AssignableTo(expectType) {
		return "ok"
	}
	situation := "mismatch"
	if gotType != expectType {
		if gotType.Kind() == reflect.Ptr && expectType.Kind() != reflect.Ptr {
			if gotType.Elem() == expectType {
				situation = "ptrToValue"
			}
		} else if gotType.Kind() != reflect.Ptr && expectType.Kind() == reflect.Ptr {
			if gotType == expectType.Elem() {
				situation = "valueToPtr"
			}
		}
	}
	return situation
}

type panicIfError struct {
}

func (panicIfError) Handle(_ *context.Context, err error) (breakOperation bool) {
	breakOperation = true
	panic(err)
}

type nilIfError struct {
}

func (nilIfError) Handle(_ *context.Context, _ error) (breakOperation bool) {
	breakOperation = false
	return
}

type errorToCtx struct {
}

func (errorToCtx) Handle(ctx *context.Context, err error) (breakOperation bool) {
	breakOperation = false
	ctx.Input.SetData("error", err)
	return
}
