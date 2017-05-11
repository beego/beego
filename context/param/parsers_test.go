package param

import "testing"
import "reflect"
import "time"

type testDefinition struct {
	strValue       string
	expectedValue  interface{}
	expectedParser paramParser
}

func Test_Parsers(t *testing.T) {
	checkParser(testDefinition{"1", 1, intParser{}}, t)

	checkParser(testDefinition{"1.0", 1.0, floatParser{}}, t)

	checkParser(testDefinition{"1", "1", stringParser{}}, t)

	checkParser(testDefinition{"true", true, boolParser{}}, t)

	checkParser(testDefinition{"2017-05-30T13:54:53Z", time.Date(2017, 5, 30, 13, 54, 53, 0, time.UTC), timeParser{}}, t)

	checkParser(testDefinition{`{"X": 5}`, struct{ X int }{5}, jsonParser{}}, t)

	checkParser(testDefinition{`1,2`, []int{1, 2}, sliceParser(intParser{})}, t)
}

func checkParser(def testDefinition, t *testing.T) {
	toType := reflect.TypeOf(def.expectedValue)
	parser := getParser(&MethodParam{}, toType)

	if reflect.TypeOf(parser) != reflect.TypeOf(def.expectedParser) {
		t.Errorf("Invalid parser for value %v. Expected: %v, actual: %v", def.strValue, reflect.TypeOf(def.expectedParser).Name(), reflect.TypeOf(parser).Name())
		return
	}
	result, err := parser.parse(def.strValue, toType)
	if err != nil {
		t.Errorf("Parsing error for value %v. Expected result: %v, error: %v", def.strValue, def.expectedValue, err)
		return
	}
	if !reflect.DeepEqual(result, def.expectedValue) {
		t.Errorf("Parsing error for value %v. Expected result: %v, actual: %v", def.strValue, def.expectedValue, result)
	}
}
