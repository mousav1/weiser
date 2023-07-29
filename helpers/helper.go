package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func urlEncode(s string) string {
	return strings.Replace(url.QueryEscape(s), "+", "%20", -1)
}

func arrayFilter(arr []interface{}, f func(interface{}) bool) []interface{} {
	var result []interface{}
	for _, v := range arr {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
func arrayMap(arr []interface{}, f func(interface{}) interface{}) []interface{} {
	var result []interface{}
	for _, v := range arr {
		result = append(result, f(v))
	}
	return result
}
func arrayReduce(arr []interface{}, f func(interface{}, interface{}) interface{}) interface{} {
	if len(arr) == 0 {
		return nil
	}
	result := arr[0]
	for i := 1; i < len(arr); i++ {
		result = f(result, arr[i])
	}
	return result
}
func arrayReverse(arr []interface{}) []interface{} {
	for i := 0; i < len(arr)/2; i++ {
		j := len(arr) - i - 1
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
func arrayUnique(arr []interface{}) []interface{} {
	var result []interface{}
	encountered := map[interface{}]bool{}
	for _, v := range arr {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}
func arrayChunk(arr []interface{}, size int) [][]interface{} {
	var result [][]interface{}
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		result = append(result, arr[i:end])
	}
	return result
}

func arrayDiff(arr []interface{}, arrs ...[]interface{}) []interface{} {
	m := make(map[interface{}]bool)
	for _, v := range arrs {
		for _, x := range v {
			m[x] = true
		}
	}
	var result []interface{}
	for _, x := range arr {
		if !m[x] {
			result = append(result, x)
		}
	}
	return result
}

func inArray(val interface{}, arr []interface{}) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
func arrayKeys(arr map[interface{}]interface{}) []interface{} {
	var keys []interface{}
	for k := range arr {
		keys = append(keys, k)
	}
	return keys
}

func arrayValues(arr map[interface{}]interface{}) []interface{} {
	var values []interface{}
	if len(arr) == 0 {
		return []interface{}{}
	}
	for _, v := range arr {
		values = append(values, v)
	}
	return values
}

func arrayMerge(arrs ...[]interface{}) []interface{} {
	var result []interface{}
	allEmpty := true
	for _, arr := range arrs {
		if len(arr) > 0 {
			allEmpty = false
			result = append(result, arr...)
		}
	}
	if allEmpty {
		return []interface{}{}
	}
	return result
}
func arrayPad(arr []interface{}, length int, value interface{}) []interface{} {
	if len(arr) >= length {
		return arr
	}
	diff := length - len(arr)
	padding := make([]interface{}, diff)
	for i := range padding {
		padding[i] = value
	}
	return append(arr, padding...)
}
func arraySearch(needle interface{}, haystack []interface{}) int {
	for i, v := range haystack {
		if v == needle {
			return i
		}
	}
	return -1
}

func arrayWalk(arr []interface{}, f func(interface{})) {
	for _, v := range arr {
		f(v)
	}
}

func arrayColumn(arr [][]interface{}, columnIndex int) []interface{} {
	var result []interface{}
	found := false
	for _, row := range arr {
		if columnIndex >= 0 && columnIndex < len(row) {
			result = append(result, row[columnIndex])
			found = true
		}
	}
	if !found {
		return []interface{}{}
	}
	return result
}

func arrayKeyExists(key interface{}, arr map[interface{}]interface{}) bool {
	_, ok := arr[key]
	return ok
}
func arrayRand(arr []interface{}, num int) []interface{} {
	if num >= len(arr) {
		return arr
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr[:num]
}
func arrayReplace(arr []interface{}, replace map[interface{}]interface{}) []interface{} {
	for i, v := range arr {
		if newVal, ok := replace[v]; ok {
			arr[i] = newVal
		}
	}
	return arr
}
func arrayWalkRecursive(arr []interface{}, f func(interface{})) {
	for _, v := range arr {
		switch val := v.(type) {
		case []interface{}:
			arrayWalkRecursive(val, f)
		default:
			f(v)
		}
	}
}
func arrayFlip(arr map[interface{}]interface{}) map[interface{}]interface{} {
	result := make(map[interface{}]interface{})
	for k, v := range arr {
		result[v] = k
	}
	return result
}
func arrayFill(length int, value interface{}) []interface{} {
	result := make([]interface{}, length)
	for i := range result {
		result[i] = value
	}
	return result
}

func arrayIntersect(arrs ...[]interface{}) []interface{} {
	if len(arrs) == 0 {
		return nil
	}
	if len(arrs) == 1 {
		return arrs[0]
	}
	result := make([]interface{}, 0)
	for _, v := range arrs[0] {
		inAll := true
		for _, a := range arrs[1:] {
			if !contains(a, v) {
				inAll = false
				break
			}
		}
		if inAll {
			result = append(result, v)
		}
	}
	return result
}

func contains(arr []interface{}, val interface{}) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func arrayPush(arr []interface{}, values ...interface{}) []interface{} {
	return append(arr, values...)
}

func arrayPop(arr []interface{}) (interface{}, []interface{}) {
	if len(arr) == 0 {
		return nil, nil
	}
	return arr[len(arr)-1], arr[:len(arr)-1]
}

func arrayGet(arr map[interface{}]interface{}, key interface{}, defaultValue interface{}) interface{} {
	if val, ok := arr[key]; ok {
		return val
	}
	return defaultValue
}

func stringContains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func stringStartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func stringEndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

func stringReverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func intToString(i int) string {
	return strconv.Itoa(i)
}

func sliceShuffle(s []interface{}) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

func sliceSum(nums []int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

func sliceFilter(s []interface{}, f func(interface{}) bool) []interface{} {
	filtered := make([]interface{}, 0)
	for _, v := range s {
		if f(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func sliceMap(s []interface{}, f func(interface{}) interface{}) []interface{} {
	mapped := make([]interface{}, len(s))
	for i, v := range s {
		mapped[i] = f(v)
	}
	return mapped
}

func mapFilter(m map[interface{}]interface{}, f func(interface{}, interface{}) bool) map[interface{}]interface{} {
	filtered := make(map[interface{}]interface{})
	for k, v := range m {
		if f(k, v) {
			filtered[k] = v
		}
	}
	return filtered
}

func mapMap(m map[interface{}]interface{}, f func(interface{}, interface{}) (interface{}, interface{})) map[interface{}]interface{} {
	mapped := make(map[interface{}]interface{})
	for k, v := range m {
		newKey, newVal := f(k, v)
		mapped[newKey] = newVal
	}
	return mapped
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	fmt.Println("Error:", err) // Print the error for debugging
	return err == nil
}

func dd(vars ...interface{}) {
	for _, v := range vars {
		fmt.Printf("%v (%v)\n", v, reflect.TypeOf(v))
	}
	panic("dd")
}

func strContains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func strReplace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil && strings.TrimSpace(s) != ""
}

func stringToInt(s string, defaultValue int) int {
	if !isNumeric(s) {
		return defaultValue
	}
	i, _ := strconv.Atoi(s)
	return i
}

func arraySlice(arr []interface{}, start, length int) ([]interface{}, error) {
	if start < 0 || start >= len(arr) || length < 0 || start+length > len(arr) {
		return nil, errors.New("invalid start or length")
	}
	return arr[start : start+length], nil
}

func fileReadLines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func fileWriteLines(filePath string, lines []string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func generateRandomString(length int) string {
	if length <= 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func isValidUrl(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
