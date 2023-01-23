package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Comparable interface for custom types for testing
type Comparable[T any] interface {
	Equal(T) bool
}

// assertEqual fails if the two values are not equal
func assertEqual[T comparable](t testing.TB, got, want T) error {
	t.Helper()

	if got != want {
		return fmt.Errorf("got: %v != want: %v", got, want)
	}

	return nil
}

// assertEqualComparable is the same as assertEqual but using Comparable interface for generic T type
func assertEqualComparable[T Comparable[T]](t testing.TB, got, want T) error {
	t.Helper()

	if !got.Equal(want) {
		return fmt.Errorf("got: %v != want: %v", got, want)
	}

	return nil
}

// checkVar checks the internal fields of Var
func checkVar[T comparable](t testing.TB, v Var[T], wantSet, wantValid bool, wantValue T) error {
	t.Helper()

	errSet := assertEqual(t, v.IsSet(), wantSet)
	errValid := assertEqual(t, v.Valid(), wantValid)
	errValue := assertEqual(t, v.Val(), wantValue)

	finalErr := []string{}
	if errSet != nil {
		finalErr = append(finalErr, fmt.Sprintf("[set] %v", errSet))
	}
	if errValid != nil {
		finalErr = append(finalErr, fmt.Sprintf("[valid] %v", errValid))
	}
	if errValue != nil {
		finalErr = append(finalErr, fmt.Sprintf("[value] %v", errValue))

	}

	if len(finalErr) == 0 {
		return nil
	}

	return errors.New(strings.Join(finalErr, " "))
}

// checkVarComparable is the same as checkVar but using Comparable interface for generic T type
func checkVarComparable[T Comparable[T]](t testing.TB, v Var[T], wantSet, wantValid bool, wantValue T) error {
	t.Helper()

	errSet := assertEqual(t, v.IsSet(), wantSet)
	errValid := assertEqual(t, v.Valid(), wantValid)
	errValue := assertEqualComparable(t, v.Val(), wantValue)

	finalErr := []string{}
	if errSet != nil {
		finalErr = append(finalErr, fmt.Sprintf("[set] %v", errSet))
	}
	if errValid != nil {
		finalErr = append(finalErr, fmt.Sprintf("[valid] %v", errValid))
	}
	if errValue != nil {
		finalErr = append(finalErr, fmt.Sprintf("[value] %v", errValue))

	}

	if len(finalErr) == 0 {
		return nil
	}

	return errors.New(strings.Join(finalErr, " "))
}

// customDefinedInt64 is a custom int64
type customDefinedInt64 int64

// customDefinedString is a custom string
type customDefinedString string

// customDefinedSlice is the same as userDefinedSlice in convert_test.go
// except it has the sql.Scanner interface implemented. It was only necessary
// for both tests to pass, because the tests in convert_test.go are checking for
// an error, since userDefinedSlice is missing the sql.Scanner interface implementation
type customDefinedSlice []int

// Equal implements the Comparable interface for customDefinedSlice
func (u customDefinedSlice) Equal(o customDefinedSlice) bool {
	if len(u) != len(o) {
		return false
	}

	for i := 0; i < len(u); i++ {
		if u[i] != o[i] {
			return false
		}
	}

	return true
}

// Scan implements the sql.Scanner interface
func (u *customDefinedSlice) Scan(src any) error {
	switch val := src.(type) {
	case []byte:
		return json.Unmarshal(val, u)
	case nil:
		return nil
	default:
		return fmt.Errorf("incompatible type for %T", u)
	}
}

// Value implements the driver.Valuer interface
func (u customDefinedSlice) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// customDefinedStruct is a struct to test struct as an underlying value for Var[T]
type customDefinedStruct struct {
	Str   string `json:"str"`
	Int64 int64  `json:"int64"`
	Slice []int  `json:"slice"`
}

// Equal implements the Comparable interface for customDefinedStruct
func (u customDefinedStruct) Equal(o customDefinedStruct) bool {
	if len(u.Slice) != len(o.Slice) {
		return false
	}

	for i := 0; i < len(u.Slice); i++ {
		if u.Slice[i] != o.Slice[i] {
			return false
		}
	}

	return u.Str == o.Str && u.Int64 == o.Int64
}

// Scan implements the sql.Scanner interface
func (u *customDefinedStruct) Scan(src any) error {
	switch val := src.(type) {
	case []byte:
		return json.Unmarshal(val, u)
	case nil:
		return nil
	default:
		return fmt.Errorf("incompatible type for %T", u)
	}
}

// customDefinedMap is a map to test a map as an underlying value for Var[T]
type customDefinedMap[K, V comparable] map[K]V

func (u customDefinedMap[K, V]) Equal(o customDefinedMap[K, V]) bool {
	if len(u) != len(o) {
		return false
	}

	for k := range u {
		if _, ok := o[k]; ok {
			if u[k] != o[k] {
				return false
			}
		}
	}

	return true
}

// Scan implements the sql.Scanner interface
func (u *customDefinedMap[K, V]) Scan(src any) error {
	switch val := src.(type) {
	case []byte:
		return json.Unmarshal(val, u)
	case nil:
		return nil
	default:
		return fmt.Errorf("incompatible type for %T", u)
	}
}

// testCase is a struct for creating test cases
type testCase struct {

	// string
	str_v            Var[string]
	str_set_for_test bool
	str_scan         any
	str_expect       string

	// ------------------------------------------------------------------------

	// int
	int_v            Var[int]
	int_set_for_test bool
	int_scan         any
	int_expect       int

	// int8
	int8_v            Var[int8]
	int8_set_for_test bool
	int8_scan         any
	int8_expect       int8

	// int16
	int16_v            Var[int16]
	int16_set_for_test bool
	int16_scan         any
	int16_expect       int16

	// int32
	int32_v            Var[int32]
	int32_set_for_test bool
	int32_scan         any
	int32_expect       int32

	// int64
	int64_v            Var[int64]
	int64_set_for_test bool
	int64_scan         any
	int64_expect       int64

	// ------------------------------------------------------------------------

	// uint
	uint_v            Var[uint]
	uint_set_for_test bool
	uint_scan         any
	uint_expect       uint

	// uint8
	uint8_v            Var[uint8]
	uint8_set_for_test bool
	uint8_scan         any
	uint8_expect       uint8

	// uint16
	uint16_v            Var[uint16]
	uint16_set_for_test bool
	uint16_scan         any
	uint16_expect       uint16

	// uint32
	uint32_v            Var[uint32]
	uint32_set_for_test bool
	uint32_scan         any
	uint32_expect       uint32

	// uint64
	uint64_v            Var[uint64]
	uint64_set_for_test bool
	uint64_scan         any
	uint64_expect       uint64

	// ------------------------------------------------------------------------

	// float32
	float32_v            Var[float32]
	float32_set_for_test bool
	float32_scan         any
	float32_expect       float32

	// float64
	float64_v            Var[float64]
	float64_set_for_test bool
	float64_scan         any
	float64_expect       float64

	// ------------------------------------------------------------------------

	// time.Time
	time_v            Var[time.Time]
	time_set_for_test bool
	time_scan         any
	time_expect       time.Time

	// ------------------------------------------------------------------------

	// customDefinedInt64
	customDefinedInt64_v            Var[customDefinedInt64]
	customDefinedInt64_set_for_test bool
	customDefinedInt64_scan         any
	customDefinedInt64_expect       customDefinedInt64

	// customDefinedSlice
	customDefinedSlice_v            Var[customDefinedSlice]
	customDefinedSlice_set_for_test bool
	customDefinedSlice_scan         any
	customDefinedSlice_expect       customDefinedSlice

	// customDefinedString
	customDefinedString_v            Var[customDefinedString]
	customDefinedString_set_for_test bool
	customDefinedString_scan         any
	customDefinedString_expect       customDefinedString

	// customDefinedStruct
	customDefinedStruct_v            Var[customDefinedStruct]
	customDefinedStruct_set_for_test bool
	customDefinedStruct_scan         any
	customDefinedStruct_expect       customDefinedStruct

	// customDefinedMap
	customDefinedMap_v            Var[customDefinedMap[string, string]]
	customDefinedMap_set_for_test bool
	customDefinedMap_scan         any
	customDefinedMap_expect       customDefinedMap[string, string]

	// ------------------------------------------------------------------------

}

// testCases return a slice of test cases
func testCases() []testCase {
	// scan values can be the following according to the database/sql package docs:
	//
	//	int64
	//  float64
	//  bool
	//  []byte
	//  string
	//  time.Time

	t := time.Unix(1672531261, 0)

	uds := customDefinedStruct{
		Str:   "foo",
		Int64: 123,
		Slice: []int{1, 2, 3},
	}
	udsj, _ := json.Marshal(uds)

	udm := customDefinedMap[string, string]{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	udmj, _ := json.Marshal(udm)

	return []testCase{
		{str_v: Var[string]{}, str_set_for_test: true, str_scan: []byte("foo"), str_expect: "foo"},

		{int_v: Var[int]{}, int_set_for_test: true, int_scan: int64(1), int_expect: 1},
		{int8_v: Var[int8]{}, int8_set_for_test: true, int8_scan: int64(1), int8_expect: 1},
		{int16_v: Var[int16]{}, int16_set_for_test: true, int16_scan: int64(1), int16_expect: 1},
		{int32_v: Var[int32]{}, int32_set_for_test: true, int32_scan: int64(1), int32_expect: 1},
		{int64_v: Var[int64]{}, int64_set_for_test: true, int64_scan: int64(1), int64_expect: 1},

		{uint_v: Var[uint]{}, uint_set_for_test: true, uint_scan: int64(1), uint_expect: 1},
		{uint8_v: Var[uint8]{}, uint8_set_for_test: true, uint8_scan: int64(1), uint8_expect: 1},
		{uint16_v: Var[uint16]{}, uint16_set_for_test: true, uint16_scan: int64(1), uint16_expect: 1},
		{uint32_v: Var[uint32]{}, uint32_set_for_test: true, uint32_scan: int64(1), uint32_expect: 1},
		{uint64_v: Var[uint64]{}, uint64_set_for_test: true, uint64_scan: int64(1), uint64_expect: 1},

		{float32_v: Var[float32]{}, float32_set_for_test: true, float32_scan: float64(1.5), float32_expect: 1.5},
		{float64_v: Var[float64]{}, float64_set_for_test: true, float64_scan: float64(1.5), float64_expect: 1.5},

		{time_v: Var[time.Time]{}, time_set_for_test: true, time_scan: t, time_expect: t},

		{customDefinedInt64_v: Var[customDefinedInt64]{}, customDefinedInt64_set_for_test: true, customDefinedInt64_scan: int64(1), customDefinedInt64_expect: 1},
		{customDefinedString_v: Var[customDefinedString]{}, customDefinedString_set_for_test: true, customDefinedString_scan: []byte("foo"), customDefinedString_expect: "foo"},
		{customDefinedSlice_v: Var[customDefinedSlice]{}, customDefinedSlice_set_for_test: true, customDefinedSlice_scan: []byte(`[1,2,3]`), customDefinedSlice_expect: customDefinedSlice{1, 2, 3}},
		{customDefinedStruct_v: Var[customDefinedStruct]{}, customDefinedStruct_set_for_test: true, customDefinedStruct_scan: udsj, customDefinedStruct_expect: uds},
		{customDefinedMap_v: Var[customDefinedMap[string, string]]{}, customDefinedMap_set_for_test: true, customDefinedMap_scan: udmj, customDefinedMap_expect: udm},
	}
}

// testSingleCase is a helper function that helps to test a single test case
func testSingleCase[T comparable](t testing.TB, v Var[T], scan any, expect T) error {
	t.Helper()

	jsonv, _ := json.Marshal(expect)

	var defExpect T
	err := checkVar(t, v, false, false, defExpect)
	if err != nil {
		return fmt.Errorf("[default] %w", err)
	}

	v.Set(expect)
	err = checkVar(t, v, true, true, expect)
	if err != nil {
		return fmt.Errorf("[set value] %w", err)
	}

	j, jErr := json.Marshal(v)
	if jErr != nil {
		return fmt.Errorf("[json.Marshal - value] %w", jErr)
	}

	err = assertEqual(t, bytes.Equal(j, jsonv), true)
	if err != nil {
		return fmt.Errorf("[json.Marshal - value] %w", err)
	}

	v.SetNil()
	err = checkVar(t, v, true, false, defExpect)
	if err != nil {
		return fmt.Errorf("[set nil] %w", err)
	}

	j, jErr = json.Marshal(v)
	if jErr != nil {
		return fmt.Errorf("[json.Marshal - nil] %w", jErr)
	}
	err = assertEqual(t, bytes.Equal(j, nullBytes), true)
	if err != nil {
		return fmt.Errorf("[json.Marshal - nil] %w", err)
	}

	v.Unset()
	err = checkVar(t, v, false, false, defExpect)
	if err != nil {
		return fmt.Errorf("[unset] %w", err)
	}

	jErr = json.Unmarshal(jsonv, &v)
	if jErr != nil {
		return fmt.Errorf("[json.Unmarshal - value] %w", err)
	}
	err = checkVar(t, v, true, true, expect)
	if err != nil {
		return fmt.Errorf("[json.Unmarshal - value] %w", err)
	}

	v.Unset()
	jErr = json.Unmarshal(nullBytes, &v)
	if jErr != nil {
		return fmt.Errorf("[json.Unmarshal - nil] %w", err)
	}
	err = checkVar(t, v, true, false, defExpect)
	if err != nil {
		return fmt.Errorf("[json.Unmarshal - nil] %w", err)
	}

	err = convertAssign(&v, scan)
	if err != nil {
		return fmt.Errorf("[convertAssign - value] %w", err)
	}
	err = checkVar(t, v, true, true, expect)
	if err != nil {
		return fmt.Errorf("[convertAssign - value] %w", err)
	}

	v.Unset()
	err = convertAssign(&v, nil)
	if err != nil {
		return fmt.Errorf("[convertAssign - nil] %w", err)
	}

	err = checkVar(t, v, true, false, defExpect)
	if err != nil {
		return fmt.Errorf("[convertAssign - nil] %w", err)
	}

	_, err = v.Value()
	if err != nil {
		return fmt.Errorf("[driver.Value - nil] %w", err)
	}

	v.Set(expect)
	_, err = v.Value()
	if err != nil {
		return fmt.Errorf("[driver.Value - value] %w", err)
	}

	return nil
}

// testSingleCaseComparable is a helper function that helps to test a single test case
func testSingleCaseComparable[T Comparable[T]](t testing.TB, v Var[T], scan any, expect T) error {
	t.Helper()

	jsonv, _ := json.Marshal(expect)

	var defExpect T
	err := checkVarComparable(t, v, false, false, defExpect)
	if err != nil {
		return fmt.Errorf("[default] %w", err)
	}

	v.Set(expect)
	err = checkVarComparable(t, v, true, true, expect)
	if err != nil {
		return fmt.Errorf("[set value] %w", err)
	}

	j, jErr := json.Marshal(v)
	if jErr != nil {
		return fmt.Errorf("[json.Marshal - value] %w", jErr)
	}

	err = assertEqual(t, bytes.Equal(j, jsonv), true)
	if err != nil {
		return fmt.Errorf("[json.Marshal - value] %w", err)
	}

	v.SetNil()
	err = checkVarComparable(t, v, true, false, defExpect)
	if err != nil {
		return fmt.Errorf("[set nil] %w", err)
	}

	j, jErr = json.Marshal(v)
	if jErr != nil {
		return fmt.Errorf("[json.Marshal - nil] %w", jErr)
	}
	err = assertEqual(t, bytes.Equal(j, nullBytes), true)
	if err != nil {
		return fmt.Errorf("[json.Marshal - nil] %w", err)
	}

	v.Unset()
	err = checkVarComparable(t, v, false, false, defExpect)
	if err != nil {
		return fmt.Errorf("[unset] %w", err)
	}

	jErr = json.Unmarshal(jsonv, &v)
	if jErr != nil {
		return fmt.Errorf("[json.Unmarshal - value] %w", err)
	}
	err = checkVarComparable(t, v, true, true, expect)
	if err != nil {
		return fmt.Errorf("[json.Unmarshal - value] %w", err)
	}

	v.Unset()
	jErr = json.Unmarshal(nullBytes, &v)
	if jErr != nil {
		return fmt.Errorf("[json.Unmarshal - nil] %w", err)
	}
	err = checkVarComparable(t, v, true, false, defExpect)
	if err != nil {
		return fmt.Errorf("[json.Unmarshal - nil] %w", err)
	}

	err = convertAssign(&v, scan)
	if err != nil {
		return fmt.Errorf("[convertAssign - value] %w", err)
	}
	err = checkVarComparable(t, v, true, true, expect)
	if err != nil {
		return fmt.Errorf("[convertAssign - value] %w", err)
	}

	v.Unset()
	err = convertAssign(&v, nil)
	if err != nil {
		return fmt.Errorf("[convertAssign - nil] %w", err)
	}

	err = checkVarComparable(t, v, true, false, defExpect)
	if err != nil {
		return fmt.Errorf("[convertAssign - nil] %w", err)
	}

	_, err = v.Value()
	if err != nil {
		return fmt.Errorf("[driver.Value - nil] %w", err)
	}

	v.Set(expect)
	_, err = v.Value()
	if err != nil {
		return fmt.Errorf("[driver.Value - value] %w", err)
	}

	return nil
}

// TestNullVar tests the Var functionality with string type
func TestNullVar(t *testing.T) {
	errF := func(n int, err error) {
		t.Errorf("testCase #%v: %v", n, err)
	}

	for n, tc := range testCases() {

		// string
		if tc.str_set_for_test {
			err := testSingleCase(t, tc.str_v, tc.str_scan, tc.str_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// int
		if tc.int_set_for_test {
			err := testSingleCase(t, tc.int_v, tc.int_scan, tc.int_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// int8
		if tc.int8_set_for_test {
			err := testSingleCase(t, tc.int8_v, tc.int8_scan, tc.int8_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// int16
		if tc.int16_set_for_test {
			err := testSingleCase(t, tc.int16_v, tc.int16_scan, tc.int16_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// int32
		if tc.int32_set_for_test {
			err := testSingleCase(t, tc.int32_v, tc.int32_scan, tc.int32_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// int64
		if tc.int64_set_for_test {
			err := testSingleCase(t, tc.int64_v, tc.int64_scan, tc.int64_expect)
			if err != nil {
				errF(n, err)
			}
		}
		// uint
		if tc.uint_set_for_test {
			err := testSingleCase(t, tc.uint_v, tc.uint_scan, tc.uint_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// uint8
		if tc.uint8_set_for_test {
			err := testSingleCase(t, tc.uint8_v, tc.uint8_scan, tc.uint8_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// uint16
		if tc.uint16_set_for_test {
			err := testSingleCase(t, tc.uint16_v, tc.uint16_scan, tc.uint16_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// uint32
		if tc.uint32_set_for_test {
			err := testSingleCase(t, tc.uint32_v, tc.uint32_scan, tc.uint32_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// uint64
		if tc.uint64_set_for_test {
			err := testSingleCase(t, tc.uint64_v, tc.uint64_scan, tc.uint64_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// float32
		if tc.float32_set_for_test {
			err := testSingleCase(t, tc.float32_v, tc.float32_scan, tc.float32_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// float64
		if tc.float64_set_for_test {
			err := testSingleCase(t, tc.float64_v, tc.float64_scan, tc.float64_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// time.Time
		if tc.time_set_for_test {
			err := testSingleCase(t, tc.time_v, tc.time_scan, tc.time_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// customDefinedInt64
		if tc.customDefinedInt64_set_for_test {

			err := testSingleCase(t, tc.customDefinedInt64_v, tc.customDefinedInt64_scan, tc.customDefinedInt64_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// customDefinedString
		if tc.customDefinedString_set_for_test {
			err := testSingleCase(t, tc.customDefinedString_v, tc.customDefinedString_scan, tc.customDefinedString_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// customDefinedSlice
		if tc.customDefinedSlice_set_for_test {

			err := testSingleCaseComparable(t, tc.customDefinedSlice_v, tc.customDefinedSlice_scan, tc.customDefinedSlice_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// customDefinedStruct
		if tc.customDefinedStruct_set_for_test {

			err := testSingleCaseComparable(t, tc.customDefinedStruct_v, tc.customDefinedStruct_scan, tc.customDefinedStruct_expect)
			if err != nil {
				errF(n, err)
			}
		}

		// customDefinedMap
		if tc.customDefinedMap_set_for_test {

			err := testSingleCaseComparable(t, tc.customDefinedMap_v, tc.customDefinedMap_scan, tc.customDefinedMap_expect)
			if err != nil {
				errF(n, err)
			}
		}
	}

}
