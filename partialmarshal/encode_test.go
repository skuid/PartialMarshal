package partialmarshal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleMarshal() {
	// A struct type with partialmarshal.Extra included as an embedded type
	type examplestruct struct {
		ExampleFieldOne string
		ExampleFieldTwo string `json:"example_field_two"`
		Extra
	}

	// An instance of that type
	source := examplestruct{
		"Value 1",
		"Value 2",
		Extra{
			"some_other_field": "some other value",
		},
	}

	// Marshaling into JSON-formatted string.
	JSONData, _ := Marshal(source)
	fmt.Println(string(JSONData))

	// Output:
	// {"ExampleFieldOne":"Value 1","example_field_two":"Value 2","some_other_field":"some other value"}
}

func TestMarshal(t *testing.T) {
	testCases := []struct {
		testDescription string
		inStruct        interface{}
		outData         []byte
		outErrMsg       string
	}{
		// Happy Path Cases
		{
			"should marshal extra field into top-level keys in payload from struct pointer",
			&struct {
				FieldOne string
				Extra
			}{
				"value one",
				Extra{
					"field_two": "value two",
				},
			},
			[]byte(`{"FieldOne":"value one","field_two":"value two"}`),
			"",
		},
		{
			"should marshal extra field into top-level keys in payload from non-pointer struct",
			struct {
				FieldOne string
				Extra
			}{
				"value one",
				Extra{
					"field_two": "value two",
				},
			},
			[]byte(`{"FieldOne":"value one","field_two":"value two"}`),
			"",
		},
		{
			"should marshal struct fields using json tags",
			&struct {
				FieldOne string `json:"field_one"`
				Extra
			}{
				"value one",
				Extra{
					"field_two": "value two",
				},
			},
			[]byte(`{"field_one":"value one","field_two":"value two"}`),
			"",
		},
		// Sad Path Cases
		{
			"should return error when no partialmarshal.Extra embedded type present",
			&struct {
				FieldOne string `json:"field_one"`
			}{
				"value one",
			},
			nil,
			"no partialmarshal.Extra embedded type found in provided struct",
		},
		{
			"should return error when provided value not struct/struct pointer",
			"",
			nil,
			"value must be of type struct",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDescription, func(t *testing.T) {

			result, err := Marshal(tc.inStruct)
			if tc.outErrMsg != "" {
				assert.EqualError(t, err, tc.outErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.outData, result)
			}
		})
	}
}
