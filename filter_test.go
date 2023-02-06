package null

import (
	"fmt"
	"testing"
)

// assertEqualTerminateTest makes the test fail if the two values are not equal
func assertEqualTerminateTest[T comparable](t testing.TB, got, want T) {
	t.Helper()

	if got != want {
		t.Errorf("got: %v != want: %v", got, want)
	}
}

func TestFilterStruct(t *testing.T) {
	type S3 struct {
		Filterable

		SomeField string `json:"some_field"`
	}

	type S2 struct {
		Filterable

		NoTagNullableStr      Var[string]
		unexportedNullableStr Var[string]
		OtherS2Field          string `json:"s2_other_str" custom_tag:"s2_other_str"`
	}

	type S1 struct {
		Filterable

		NullFieldStr             Var[string] `json:"str" custom_tag:"str"`
		unexportedNullFieldInt64 Var[int64]  `json:"unexported_nullfield_int64" custom_tag:"unexported_nullfield_int64"`
		OtherField               string      `json:"other_str" custom_tag:"other_str"`
		StructField              S2          `json:"s2" custom_tag:"s2"`
		AnotherStructField       S3          `json:"s3"`
		S3
	}

	_, err := FilterStruct(nil)
	assertEqualTerminateTest(t, err.Error(), "input cannot be nil")

	_, err = FilterStruct(int64(1))
	assertEqualTerminateTest(t, err.Error(), "invalid type int64. input must be a struct")

	expectDef := map[string]any{
		"other_str":  "",
		"some_field": "",
		"s2": map[string]any{
			"s2_other_str": "",
		},
		"s3": map[string]any{
			"some_field": "",
		},
	}
	def := S1{}
	filteredDef, err := FilterStruct(def)
	assertEqualTerminateTest(t, err == nil, true)
	assertEqualTerminateTest(t, fmt.Sprintf("%+v", expectDef), fmt.Sprintf("%+v", filteredDef))

	expectedDefCustomTag := map[string]any{
		"other_str": "",
		"s2": map[string]any{
			"s2_other_str": "",
		},
	}
	filteredDefCustomTag, err := FilterStruct(def, UseTag("custom_tag"))
	assertEqualTerminateTest(t, err == nil, true)
	assertEqualTerminateTest(t, fmt.Sprintf("%+v", expectedDefCustomTag), fmt.Sprintf("%+v", filteredDefCustomTag))

	// ------------------------------------------------------------------------

	type S4 struct {
		Filterable

		C Var[int64] `json:"c"`
		D Var[int64] `json:"d"`
	}

	type S6 struct {
		K string `json:"k"`
	}

	type S5 struct {
		A   Var[string]       `json:"a"`
		B   string            `json:"b"`
		M   map[string]any    `json:"m"`
		M2  map[string]string `json:"m2"`
		M3  map[string]any    `json:"m3"`
		S   S4                `json:"s"`
		S4  `json:"ss"`
		US6 S6 `json:"s6_1"`
		S6  `json:"s6_2"`
	}

	expectedDef2 := map[string]any{
		"a": nil,
		"b": "",
		"m": map[string]any{
			"z":   "z",
			"zzz": 5,
			"mm": map[string]any{
				"q":   "",
				"qq":  0.2,
				"qqq": nil,
			},
		},
		"m2": map[string]any{
			"m2":  "222",
			"m22": "2222",
		},
		"s": map[string]any{
			"d": nil,
		},
		"s6_1": S6{},
		"s6_2": S6{},
	}

	def2 := S5{
		A: Var[string]{
			set:   true,
			valid: false,
			value: "",
		},
		M: map[string]any{
			"z": Var[string]{
				set:   true,
				valid: true,
				value: "z",
			},
			"zz":  Var[string]{},
			"zzz": 5,
			"mm": map[string]any{
				"q": "",
				"qq": Var[float64]{
					set:   true,
					valid: true,
					value: 0.2,
				},
				"qqq": Var[float64]{
					set:   true,
					valid: false,
					value: 0,
				},
			},
		},
		M2: map[string]string{
			"m2":  "222",
			"m22": "2222",
		},
		M3: map[string]any{
			"m3": Var[string]{},
		},
		S: S4{
			C: Var[int64]{
				set:   false,
				valid: false,
				value: 0,
			},
			D: Var[int64]{
				set:   true,
				valid: false,
				value: 0,
			},
		},
		S4: S4{},
	}
	filteredDef2, err := FilterStruct(def2)
	assertEqualTerminateTest(t, err == nil, true)
	assertEqualTerminateTest(t, fmt.Sprintf("%+v", expectedDef2), fmt.Sprintf("%+v", filteredDef2))

	filteredDef2, err = FilterStruct(def2, UseTag(""))
	assertEqualTerminateTest(t, err == nil, true)
	assertEqualTerminateTest(t, fmt.Sprintf("%+v", expectedDef2), fmt.Sprintf("%+v", filteredDef2))
}

func TestFilterMap(t *testing.T) {
	_, err := FilterMap(nil)
	assertEqualTerminateTest(t, err.Error(), "input cannot be nil")

	m := map[string]any{
		"a": Var[string]{
			set:   true,
			valid: false,
			value: "",
		},
		"b": Var[string]{
			set:   true,
			valid: true,
			value: "B",
		},
		"c": Var[string]{},
		"d": "D",
		"e": map[string]string{
			"ee": "ee",
		},
		"f": map[string]any{
			"a": Var[string]{
				set:   true,
				valid: false,
				value: "",
			},
			"b": Var[string]{
				set:   true,
				valid: true,
				value: "B",
			},
		},
		"g": map[string]any{
			"a": Var[string]{},
		},
	}

	mExpect := map[string]any{
		"a": nil,
		"b": "B",
		"d": "D",
		"e": map[string]string{
			"ee": "ee",
		},
		"f": map[string]any{
			"a": nil,
			"b": "B",
		},
	}

	mFiltered, err := FilterMap(m)
	assertEqualTerminateTest(t, err == nil, true)
	assertEqualTerminateTest(t, fmt.Sprintf("%+v", mExpect), fmt.Sprintf("%+v", mFiltered))
}
