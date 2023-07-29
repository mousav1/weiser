package helpers

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestUrlEncode(t *testing.T) {
	// Positive test case
	input := "hello world"
	expected := "hello%20world"
	result := urlEncode(input)
	if result != expected {
		t.Errorf("urlEncode(%s) = %s; expected %s", input, result, expected)
	}
	// Negative test case
	input = "hello@world"
	expected = "hello%40world"
	result = urlEncode(input)
	if result != expected {
		t.Errorf("urlEncode(%s) = %s; expected %s", input, result, expected)
	}
}
func TestArrayFilter(t *testing.T) {
	// Positive test case
	arr := []interface{}{1, 2, 3, 4, 5}
	expected := []interface{}{2, 4}
	result := arrayFilter(arr, func(v interface{}) bool {
		return v.(int)%2 == 0
	})
	if !isEqual(result, expected) {
		t.Errorf("arrayFilter(%v) = %v; expected %v", arr, result, expected)
	}
	// Negative test case
	arr = []interface{}{"a", "b", "c", "d", "e"}
	expected = []interface{}{}
	result = arrayFilter(arr, func(v interface{}) bool {
		return v.(string) == "z"
	})
	if !isEqual(result, expected) {
		t.Errorf("arrayFilter(%v) = %v; expected %v", arr, result, expected)
	}
}
func TestArrayMap(t *testing.T) {
	// Positive test case
	arr := []interface{}{1, 2, 3, 4, 5}
	expected := []interface{}{2, 4, 6, 8, 10}
	result := arrayMap(arr, func(v interface{}) interface{} {
		return v.(int) * 2
	})
	if !isEqual(result, expected) {
		t.Errorf("arrayMap(%v) = %v; expected %v", arr, result, expected)
	}
	// Negative test case
	arr = []interface{}{"a", "b", "c", "d", "e"}
	expected = []interface{}{"A", "B", "C", "D", "E"}
	result = arrayMap(arr, func(v interface{}) interface{} {
		return strings.ToUpper(v.(string))
	})
	if !isEqual(result, expected) {
		t.Errorf("arrayMap(%v) = %v; expected %v", arr, result, expected)
	}
}

// Helper function to compare two slices
func isEqual(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
func TestArrayReverse(t *testing.T) {
	// Positive test case
	arr := []interface{}{1, 2, 3, 4, 5}
	expectedResult := []interface{}{5, 4, 3, 2, 1}
	result := arrayReverse(arr)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("arrayReverse failed, expected %v but got %v", expectedResult, result)
	}
	// Negative test case
	arr = []interface{}{}
	expectedResult = []interface{}{}
	result = arrayReverse(arr)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("arrayReverse failed, expected %v but got %v", expectedResult, result)
	}
}

func TestArrayReduce(t *testing.T) {
	// Positive test case
	arr := []interface{}{1, 2, 3, 4, 5}
	f := func(x, y interface{}) interface{} {
		return x.(int) + y.(int)
	}
	expectedResult := 15
	result := arrayReduce(arr, f)
	if result != expectedResult {
		t.Errorf("arrayReduce failed, expected %v but got %v", expectedResult, result)
	}
	// Negative test case
	arr = []interface{}{}
	result = arrayReduce(arr, f)
	if result != nil {
		t.Errorf("arrayReduce failed, expected %v but got %v", nil, result)
	}
}

func TestArrayUnique(t *testing.T) {
	// Positive test case
	input := []interface{}{1, 2, 3, 2, 4, 3, 5}
	expected := []interface{}{1, 2, 3, 4, 5}
	result := arrayUnique(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	// Negative test case
	input = []interface{}{"a", "b", "c", "b", "d", "c", "e"}
	expected = []interface{}{"a", "b", "c", "d", "e"}
	result = arrayUnique(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestArrayChunk(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5, 6}
	size := 2
	expected := [][]interface{}{{1, 2}, {3, 4}, {5, 6}}
	result := arrayChunk(arr, size)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("arrayChunk() failed, expected %v, got %v", expected, result)
	}
}
func TestArrayDiff(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	arrs := [][]interface{}{{2, 4}, {3, 5}}
	expected := []interface{}{1}
	result := arrayDiff(arr, arrs...)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("arrayDiff() failed, expected %v, got %v", expected, result)
	}
}
func TestInArray(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	val := 3
	expected := true
	result := inArray(val, arr)
	if result != expected {
		t.Errorf("inArray() failed, expected %v, got %v", expected, result)
	}
}
func TestArrayKeys(t *testing.T) {
	arr := map[interface{}]interface{}{"a": 1, "b": 2, "c": 3}
	expected := []interface{}{"a", "b", "c"}
	result := arrayKeys(arr)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("arrayKeys() failed, expected %v, got %v", expected, result)
	}
}

func TestSliceShuffle(t *testing.T) {
	// Positive test case
	input := []interface{}{1, 2, 3, 4, 5}
	expected := make([]interface{}, len(input))
	copy(expected, input)
	sliceShuffle(input)
	if reflect.DeepEqual(input, expected) {
		t.Errorf("Slice was not shuffled")
	}
	// Negative test case
	input = []interface{}{}
	expected = []interface{}{}
	sliceShuffle(input)
	if !reflect.DeepEqual(input, expected) {
		t.Errorf("Empty slice was modified")
	}
}

func TestArrayValues(t *testing.T) {
	// Positive test case
	input := map[interface{}]interface{}{"key1": "value1", "key2": "value2"}
	expected := []interface{}{"value1", "value2"}
	result := arrayValues(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	// Negative test case
	input = map[interface{}]interface{}{}
	expected = []interface{}{}
	result = arrayValues(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
func TestArrayMerge(t *testing.T) {
	// Positive test case
	arr1 := []interface{}{1, 2, 3}
	arr2 := []interface{}{4, 5, 6}
	expected := []interface{}{1, 2, 3, 4, 5, 6}
	result := arrayMerge(arr1, arr2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	// Negative test case
	arr1 = []interface{}{}
	arr2 = []interface{}{}
	expected = []interface{}{}
	result = arrayMerge(arr1, arr2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
func TestArrayPad(t *testing.T) {
	// Positive test case
	arr := []interface{}{1, 2, 3}
	length := 5
	value := 0
	expected := []interface{}{1, 2, 3, 0, 0}
	result := arrayPad(arr, length, value)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	// Negative test case
	arr = []interface{}{1, 2, 3}
	length = 3
	value = 0
	expected = []interface{}{1, 2, 3}
	result = arrayPad(arr, length, value)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
func TestArraySearch(t *testing.T) {
	// Positive test case
	needle := "value2"
	haystack := []interface{}{"value1", "value2", "value3"}
	expected := 1
	result := arraySearch(needle, haystack)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	// Negative test case
	needle = "value4"
	haystack = []interface{}{"value1", "value2", "value3"}
	expected = -1
	result = arraySearch(needle, haystack)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
func TestArrayWalk(t *testing.T) {
	// Positive test case
	arr := []interface{}{1, 2, 3}
	expected := []interface{}{2, 3, 4}
	arrayWalk(arr, func(v interface{}) {
		vInt := v.(int)
		vInt++
		if vInt != 2 && vInt != 3 && vInt != 4 {
			t.Errorf("Expected %v, but got %v", expected, arr)
		}
	})
	// Negative test case
	arr = []interface{}{}
	expected = []interface{}{}
	arrayWalk(arr, func(v interface{}) {
		t.Errorf("This function should not be called for an empty array")
	})
}

func TestArrayColumn(t *testing.T) {
	arr := [][]interface{}{
		{1, 2, 3},
		{"a", "b", "c"},
		{true, false, true},
	}
	columnIndex := 1
	expected := []interface{}{2, "b", false}
	result := arrayColumn(arr, columnIndex)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	columnIndex = 3
	expected = []interface{}{}
	result = arrayColumn(arr, columnIndex)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestArrayKeyExists(t *testing.T) {
	arr := map[interface{}]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	key := "key2"
	expected := true
	result := arrayKeyExists(key, arr)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	key = "key4"
	expected = false
	result = arrayKeyExists(key, arr)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestArrayRand(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	num := 3
	result := arrayRand(arr, num)
	if len(result) != num {
		t.Errorf("Expected length of result to be %d, but got %d", num, len(result))
	}
	for _, v := range result {
		if !contains(arr, v) {
			t.Errorf("Result contains value that is not in the original array")
		}
	}
}
func TestArrayReplace(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	replace := map[interface{}]interface{}{
		2: "two",
		4: "four",
	}
	expected := []interface{}{1, "two", 3, "four", 5}
	result := arrayReplace(arr, replace)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
func TestArrayWalkRecursive(t *testing.T) {
	arr := []interface{}{
		1,
		[]interface{}{
			2,
			[]interface{}{
				3,
			},
		},
	}
	var result []interface{}
	f := func(v interface{}) {
		result = append(result, v)
	}
	arrayWalkRecursive(arr, f)
	expected := []interface{}{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
func TestArrayFlip(t *testing.T) {
	arr := map[interface{}]interface{}{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	expected := map[interface{}]interface{}{
		1: "one",
		2: "two",
		3: "three",
	}
	result := arrayFlip(arr)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
func TestArrayFill(t *testing.T) {
	length := 5
	value := "a"
	expected := []interface{}{"a", "a", "a", "a", "a"}
	result := arrayFill(length, value)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
func TestArrayIntersect(t *testing.T) {
	arr1 := []interface{}{1, 2, 3, 4, 5}
	arr2 := []interface{}{2, 3, 4, 5, 6}
	arr3 := []interface{}{3, 4, 5, 6, 7}
	expected := []interface{}{3, 4, 5}
	result := arrayIntersect(arr1, arr2, arr3)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
func TestArrayPush(t *testing.T) {
	arr := []interface{}{1, 2, 3}
	values := []interface{}{4, 5}
	expected := []interface{}{1, 2, 3, 4, 5}
	result := arrayPush(arr, values...)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
func TestArrayPop(t *testing.T) {
	arr := []interface{}{1, 2, 3}
	expectedValue := 3
	expectedResult := []interface{}{1, 2}
	value, result := arrayPop(arr)
	if value != expectedValue {
		t.Errorf("Expected value %v, but got %v", expectedValue, value)
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected result %v, but got %v", expectedResult, result)
	}
}
func TestArrayGet(t *testing.T) {
	arr := map[interface{}]interface{}{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	key := "two"
	defaultValue := "default"
	expected := 2
	result := arrayGet(arr, key, defaultValue)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	nonExistingKey := "four"
	result = arrayGet(arr, nonExistingKey, defaultValue)
	if result != defaultValue {
		t.Errorf("Expected %v, but got %v", defaultValue, result)
	}
}
func TestStringContains(t *testing.T) {
	s := "Hello, World!"
	substr := "World"
	if !stringContains(s, substr) {
		t.Errorf("Expected substring %s to be found in string %s", substr, s)
	}
	substr = "Goodbye"
	if stringContains(s, substr) {
		t.Errorf("Expected substring %s not to be found in string %s", substr, s)
	}
}
func TestStringStartsWith(t *testing.T) {
	s := "Hello, World!"
	prefix := "Hello"
	if !stringStartsWith(s, prefix) {
		t.Errorf("Expected string %s to start with prefix %s", s, prefix)
	}
	prefix = "Goodbye"
	if stringStartsWith(s, prefix) {
		t.Errorf("Expected string %s not to start with prefix %s", s, prefix)
	}
}
func TestStringEndsWith(t *testing.T) {
	s := "Hello, World!"
	suffix := "World!"
	if !stringEndsWith(s, suffix) {
		t.Errorf("Expected string %s to end with suffix %s", s, suffix)
	}
	suffix = "Goodbye"
	if stringEndsWith(s, suffix) {
		t.Errorf("Expected string %s not to end with suffix %s", s, suffix)
	}
}
func TestStringReverse(t *testing.T) {
	s := "Hello, World!"
	expected := "!dlroW ,olleH"
	result := stringReverse(s)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
func TestIntToString(t *testing.T) {
	i := 12345
	expected := "12345"
	result := intToString(i)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestSliceSum(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	expected := 15
	result := sliceSum(nums)
	if result != expected {
		t.Errorf("Expected %d but got %d", expected, result)
	}
}
func TestSliceFilter(t *testing.T) {
	s := []interface{}{1, "2", 3, "4", 5}
	expected := []interface{}{"2", "4"}
	result := sliceFilter(s, func(v interface{}) bool {
		_, ok := v.(string)
		return ok
	})
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}
func TestSliceMap(t *testing.T) {
	s := []interface{}{1, 2, 3, 4, 5}
	expected := []interface{}{2, 4, 6, 8, 10}
	result := sliceMap(s, func(v interface{}) interface{} {
		return v.(int) * 2
	})
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}
func TestMapFilter(t *testing.T) {
	m := map[interface{}]interface{}{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
	}
	expected := map[interface{}]interface{}{
		"two":  2,
		"four": 4,
	}
	result := mapFilter(m, func(k, v interface{}) bool {
		_, ok := k.(string)
		return ok && v.(int)%2 == 0
	})
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}
func TestMapMap(t *testing.T) {
	m := map[interface{}]interface{}{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
	}
	expected := map[interface{}]interface{}{
		"ONE":   1,
		"TWO":   2,
		"THREE": 3,
		"FOUR":  4,
		"FIVE":  5,
	}
	result := mapMap(m, func(k, v interface{}) (interface{}, interface{}) {
		return strings.ToUpper(k.(string)), v
	})
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}
func TestFileExists(t *testing.T) {
	// Positive test case
	expected := true
	result := fileExists("testdata/testfile.txt")
	if result != expected {
		t.Errorf("Expected %t but got %t", expected, result)
	}
	// Negative test case
	expected = false
	result = fileExists("testdata/nonexistentfile.txt")
	if result != expected {
		t.Errorf("Expected %t but got %t", expected, result)
	}
}

// dd function cannot be tested as it is a debugging function and is intended to cause a panic.

func TestStrContains(t *testing.T) {
	s := "Hello, World!"
	substr := "Hello"
	if !strContains(s, substr) {
		t.Errorf("Expected true, got false")
	}
	substr = "foo"
	if strContains(s, substr) {
		t.Errorf("Expected false, got true")
	}
}
func TestStrReplace(t *testing.T) {
	s := "Hello, World!"
	old := "Hello"
	new := "Hi"
	expected := "Hi, World!"
	result := strReplace(s, old, new)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
func TestIsEmpty(t *testing.T) {
	s := ""
	if !isEmpty(s) {
		t.Errorf("Expected true, got false")
	}
	s = "Hello"
	if isEmpty(s) {
		t.Errorf("Expected false, got true")
	}
}
func TestIsNumeric(t *testing.T) {
	s := "123.45"
	if !isNumeric(s) {
		t.Errorf("Expected true, got false")
	}
	s = "abc"
	if isNumeric(s) {
		t.Errorf("Expected false, got true")
	}
}
func TestStringToInt(t *testing.T) {
	s := "123"
	defaultValue := 0
	expected := 123
	result := stringToInt(s, defaultValue)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
	s = "abc"
	expected = defaultValue
	result = stringToInt(s, defaultValue)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
func TestArraySlice(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	start := 1
	length := 3
	expected := []interface{}{2, 3, 4}
	result, err := arraySlice(arr, start, length)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if !compareSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	start = 5
	length = 1
	_, err = arraySlice(arr, start, length)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
func compareSlices(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFileReadLines(t *testing.T) {
	filePath := "testfile.txt"
	// Positive test case
	expected := []string{"line1", "line2", "line3"}
	err := fileWriteLines(filePath, expected)
	if err != nil {
		t.Errorf("Error writing lines to file: %v", err)
	}
	lines, err := fileReadLines(filePath)
	if err != nil {
		t.Errorf("Error reading lines from file: %v", err)
	}
	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines, but got %d", len(expected), len(lines))
	}
	for i := range lines {
		if lines[i] != expected[i] {
			t.Errorf("Expected line %d to be %q, but got %q", i+1, expected[i], lines[i])
		}
	}
	// Negative test case
	_, err = fileReadLines("nonexistentfile.txt")
	if err == nil {
		t.Error("Expected error reading lines from nonexistent file, but got nil")
	}
}
func TestFileWriteLines(t *testing.T) {
	filePath := "testfile.txt"
	lines := []string{"line1", "line2", "line3"}
	// Positive test case
	err := fileWriteLines(filePath, lines)
	if err != nil {
		t.Errorf("Error writing lines to file: %v", err)
	}
	// Negative test case
	err = fileWriteLines("/root/testfile.txt", lines)
	if err == nil {
		t.Error("Expected error writing lines to inaccessible file, but got nil")
	}
	// Clean up the test file
	err = os.Remove(filePath)
	if err != nil {
		t.Errorf("Error removing test file: %v", err)
	}
}
func TestGenerateRandomString(t *testing.T) {
	length := 10
	// Positive test case
	randomString := generateRandomString(length)
	if len(randomString) != length {
		t.Errorf("Expected random string of length %d, but got length %d", length, len(randomString))
	}
	// Negative test case
	randomString = generateRandomString(-1)
	if len(randomString) != 0 {
		t.Errorf("Expected empty random string for negative length, but got length %d", len(randomString))
	}
}

func TestIsValidUrl(t *testing.T) {
	// Positive test case
	validUrl := "https://www.example.com"
	if !isValidUrl(validUrl) {
		t.Errorf("isValidUrl(%s) = false, expected true", validUrl)
	}
	// Negative test case
	invalidUrl := "not a valid url"
	if isValidUrl(invalidUrl) {
		t.Errorf("isValidUrl(%s) = true, expected false", invalidUrl)
	}
}
