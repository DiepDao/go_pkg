package array_test

import (
	"fmt"
	"testing"

	"github.com/DiepDao/go_pkg/array"
)

// Test structs
type Person struct {
	Name string
	Age  int
}

type Product struct {
	ID    int
	Name  string
	Price float64
}

// TestFilter_Integers tests filtering integers
func TestFilter_Integers(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Filter even numbers
	evens := array.Filter(numbers, func(n int) bool {
		return n%2 == 0
	})

	expected := []int{2, 4, 6, 8, 10}
	if len(evens) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(evens))
	}

	for i, v := range evens {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

// TestFilter_Strings tests filtering strings
func TestFilter_Strings(t *testing.T) {
	words := []string{"apple", "banana", "apricot", "cherry", "avocado"}

	// Filter words starting with 'a'
	aWords := array.Filter(words, func(s string) bool {
		return len(s) > 0 && s[0] == 'a'
	})

	expected := []string{"apple", "apricot", "avocado"}
	if len(aWords) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(aWords))
	}

	for i, v := range aWords {
		if v != expected[i] {
			t.Errorf("Expected %s at index %d, got %s", expected[i], i, v)
		}
	}
}

// TestFilter_Structs tests filtering structs
func TestFilter_Structs(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 17},
		{Name: "Charlie", Age: 30},
		{Name: "Diana", Age: 16},
		{Name: "Eve", Age: 22},
	}

	// Filter adults (age >= 18)
	adults := array.Filter(people, func(p Person) bool {
		return p.Age >= 18
	})

	if len(adults) != 3 {
		t.Errorf("Expected 3 adults, got %d", len(adults))
	}

	expectedNames := []string{"Alice", "Charlie", "Eve"}
	for i, person := range adults {
		if person.Name != expectedNames[i] {
			t.Errorf("Expected %s at index %d, got %s", expectedNames[i], i, person.Name)
		}
	}
}

// TestFilter_EmptySlice tests filtering an empty slice
func TestFilter_EmptySlice(t *testing.T) {
	empty := []int{}
	result := array.Filter(empty, func(n int) bool {
		return n > 0
	})

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(result))
	}
}

// TestFilter_NilPredicate tests filtering with nil predicate
func TestFilter_NilPredicate(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	result := array.Filter(numbers, nil)

	if len(result) != len(numbers) {
		t.Errorf("Expected length %d, got %d", len(numbers), len(result))
	}

	for i, v := range result {
		if v != numbers[i] {
			t.Errorf("Expected %d at index %d, got %d", numbers[i], i, v)
		}
	}

	// Verify it's a copy, not the same slice
	result[0] = 999
	if numbers[0] == 999 {
		t.Error("Filter with nil predicate should return a copy, not the same slice")
	}
}

// TestFilter_NoMatches tests filtering with no matches
func TestFilter_NoMatches(t *testing.T) {
	numbers := []int{1, 3, 5, 7, 9}
	evens := array.Filter(numbers, func(n int) bool {
		return n%2 == 0
	})

	if len(evens) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(evens))
	}
}

// TestFilter_AllMatches tests filtering where all elements match
func TestFilter_AllMatches(t *testing.T) {
	numbers := []int{2, 4, 6, 8, 10}
	evens := array.Filter(numbers, func(n int) bool {
		return n%2 == 0
	})

	if len(evens) != len(numbers) {
		t.Errorf("Expected length %d, got %d", len(numbers), len(evens))
	}
}

// TestFilterIndexed_WithIndex tests filtering using index
func TestFilterIndexed_WithIndex(t *testing.T) {
	numbers := []int{10, 20, 30, 40, 50}

	// Filter elements at even indices
	evenIndexElements := array.FilterIndexed(numbers, func(v int, index int, arr []int) bool {
		return index%2 == 0
	})

	expected := []int{10, 30, 50}
	if len(evenIndexElements) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(evenIndexElements))
	}

	for i, v := range evenIndexElements {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

// TestFilterIndexed_WithArray tests filtering using the array parameter
func TestFilterIndexed_WithArray(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	// Filter elements greater than the first element
	result := array.FilterIndexed(numbers, func(v int, index int, arr []int) bool {
		return v > arr[0]
	})

	expected := []int{2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

// TestFilterIndexed_ComplexPredicate tests complex predicate logic
func TestFilterIndexed_ComplexPredicate(t *testing.T) {
	products := []Product{
		{ID: 1, Name: "Laptop", Price: 999.99},
		{ID: 2, Name: "Mouse", Price: 25.50},
		{ID: 3, Name: "Keyboard", Price: 75.00},
		{ID: 4, Name: "Monitor", Price: 299.99},
		{ID: 5, Name: "USB Cable", Price: 10.00},
	}

	// Filter products: price > 50 AND index is odd
	result := array.FilterIndexed(products, func(p Product, index int, arr []Product) bool {
		return p.Price > 50 && index%2 == 1
	})

	// Index 1 (Mouse: 25.50) - fails price check
	// Index 3 (Monitor: 299.99) - passes both
	if len(result) != 1 {
		t.Errorf("Expected 1 product, got %d", len(result))
	}

	if len(result) > 0 && result[0].Name != "Monitor" {
		t.Errorf("Expected Monitor, got %s", result[0].Name)
	}
}

// TestFilterIndexed_NilPredicate tests FilterIndexed with nil predicate
func TestFilterIndexed_NilPredicate(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	result := array.FilterIndexed(numbers, nil)

	if len(result) != len(numbers) {
		t.Errorf("Expected length %d, got %d", len(numbers), len(result))
	}

	for i, v := range result {
		if v != numbers[i] {
			t.Errorf("Expected %d at index %d, got %d", numbers[i], i, v)
		}
	}

	// Verify it's a copy
	result[0] = 999
	if numbers[0] == 999 {
		t.Error("FilterIndexed with nil predicate should return a copy")
	}
}

// TestFilter_PreservesOrder tests that filtering preserves element order
func TestFilter_PreservesOrder(t *testing.T) {
	numbers := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}

	// Filter numbers > 5
	result := array.Filter(numbers, func(n int) bool {
		return n > 5
	})

	expected := []int{8, 9, 7, 6}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

// TestFilter_DoesNotModifyOriginal tests that original slice is not modified
func TestFilter_DoesNotModifyOriginal(t *testing.T) {
	original := []int{1, 2, 3, 4, 5}
	originalCopy := make([]int, len(original))
	copy(originalCopy, original)

	_ = array.Filter(original, func(n int) bool {
		return n%2 == 0
	})

	// Verify original is unchanged
	for i, v := range original {
		if v != originalCopy[i] {
			t.Errorf("Original slice was modified at index %d: expected %d, got %d", i, originalCopy[i], v)
		}
	}
}

// TestFilter_GenericAny tests filtering with any type
func TestFilter_GenericAny(t *testing.T) {
	mixed := []any{1, "hello", 2.5, "world", 42, true, "test"}

	// Filter only strings
	strings := array.Filter(mixed, func(v any) bool {
		_, ok := v.(string)
		return ok
	})

	if len(strings) != 3 {
		t.Errorf("Expected 3 strings, got %d", len(strings))
	}

	expectedStrings := []string{"hello", "world", "test"}
	for i, v := range strings {
		str, ok := v.(string)
		if !ok {
			t.Errorf("Expected string at index %d, got %T", i, v)
		}
		if str != expectedStrings[i] {
			t.Errorf("Expected %s at index %d, got %s", expectedStrings[i], i, str)
		}
	}
}

// BenchmarkFilter_Integers benchmarks filtering integers
func BenchmarkFilter_Integers(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = array.Filter(numbers, func(n int) bool {
			return n%2 == 0
		})
	}
}

// BenchmarkFilterIndexed_Integers benchmarks FilterIndexed with integers
func BenchmarkFilterIndexed_Integers(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = array.FilterIndexed(numbers, func(v int, index int, arr []int) bool {
			return v%2 == 0 && index < 500
		})
	}
}

// TestMap_IntToInt tests mapping integers to integers
func TestMap_IntToInt(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	// Double each number
	doubled := array.Map(numbers, func(n int) int {
		return n * 2
	})

	expected := []int{2, 4, 6, 8, 10}
	if len(doubled) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(doubled))
	}

	for i, v := range doubled {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

// TestMap_IntToString tests mapping integers to strings
func TestMap_IntToString(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	// Convert to strings
	strings := array.Map(numbers, func(n int) string {
		return fmt.Sprintf("Number: %d", n)
	})

	if len(strings) != len(numbers) {
		t.Errorf("Expected length %d, got %d", len(numbers), len(strings))
	}

	if strings[0] != "Number: 1" {
		t.Errorf("Expected 'Number: 1', got '%s'", strings[0])
	}
	if strings[4] != "Number: 5" {
		t.Errorf("Expected 'Number: 5', got '%s'", strings[4])
	}
}

// TestMap_StringToInt tests mapping strings to integers
func TestMap_StringToInt(t *testing.T) {
	words := []string{"hello", "world", "go", "test"}

	// Get length of each string
	lengths := array.Map(words, func(s string) int {
		return len(s)
	})

	expected := []int{5, 5, 2, 4}
	if len(lengths) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(lengths))
	}

	for i, v := range lengths {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

// TestMap_StructToStruct tests mapping structs to different structs
func TestMap_StructToStruct(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	type UserDTO struct {
		UserID   int
		FullName string
	}

	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}

	// Transform to DTOs
	dtos := array.Map(users, func(u User) UserDTO {
		return UserDTO{
			UserID:   u.ID,
			FullName: u.Name,
		}
	})

	if len(dtos) != len(users) {
		t.Errorf("Expected length %d, got %d", len(users), len(dtos))
	}

	if dtos[0].UserID != 1 || dtos[0].FullName != "Alice" {
		t.Errorf("Expected UserDTO{1, Alice}, got %+v", dtos[0])
	}
}

// TestMap_StructField tests extracting a field from structs
func TestMap_StructField(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
	}

	// Extract names
	names := array.Map(people, func(p Person) string {
		return p.Name
	})

	expected := []string{"Alice", "Bob", "Charlie"}
	if len(names) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(names))
	}

	for i, name := range names {
		if name != expected[i] {
			t.Errorf("Expected %s at index %d, got %s", expected[i], i, name)
		}
	}
}

// TestMap_EmptySlice tests mapping an empty slice
func TestMap_EmptySlice(t *testing.T) {
	empty := []int{}
	result := array.Map(empty, func(n int) int {
		return n * 2
	})

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(result))
	}
}

// TestMapIndexed_WithIndex tests mapping using index
func TestMapIndexed_WithIndex(t *testing.T) {
	numbers := []int{10, 20, 30, 40, 50}

	// Multiply each number by its index
	result := array.MapIndexed(numbers, func(v int, index int, arr []int) int {
		return v * index
	})

	expected := []int{0, 20, 60, 120, 200}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

// TestMapIndexed_WithArray tests mapping using the array parameter
func TestMapIndexed_WithArray(t *testing.T) {
	numbers := []int{5, 10, 15, 20, 25}

	// Calculate difference from first element
	result := array.MapIndexed(numbers, func(v int, index int, arr []int) int {
		return v - arr[0]
	})

	expected := []int{0, 5, 10, 15, 20}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

// TestMapIndexed_TypeConversion tests type conversion with index
func TestMapIndexed_TypeConversion(t *testing.T) {
	words := []string{"apple", "banana", "cherry"}

	// Create formatted strings with index
	result := array.MapIndexed(words, func(v string, index int, arr []string) string {
		return fmt.Sprintf("%d: %s", index+1, v)
	})

	expected := []string{"1: apple", "2: banana", "3: cherry"}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %s at index %d, got %s", expected[i], i, v)
		}
	}
}

// TestMap_ChainWithFilter tests chaining map and filter
func TestMap_ChainWithFilter(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Filter evens, then double them
	evens := array.Filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	doubled := array.Map(evens, func(n int) int {
		return n * 2
	})

	expected := []int{4, 8, 12, 16, 20}
	if len(doubled) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(doubled))
	}

	for i, v := range doubled {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

// TestMap_DoesNotModifyOriginal tests that original slice is not modified
func TestMap_DoesNotModifyOriginal(t *testing.T) {
	original := []int{1, 2, 3, 4, 5}
	originalCopy := make([]int, len(original))
	copy(originalCopy, original)

	_ = array.Map(original, func(n int) int {
		return n * 10
	})

	// Verify original is unchanged
	for i, v := range original {
		if v != originalCopy[i] {
			t.Errorf("Original slice was modified at index %d: expected %d, got %d", i, originalCopy[i], v)
		}
	}
}

// BenchmarkMap_Integers benchmarks mapping integers
func BenchmarkMap_Integers(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = array.Map(numbers, func(n int) int {
			return n * 2
		})
	}
}

// BenchmarkMapIndexed_Integers benchmarks MapIndexed with integers
func BenchmarkMapIndexed_Integers(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = array.MapIndexed(numbers, func(v int, index int, arr []int) int {
			return v * index
		})
	}
}

// TestIncludes_Integers tests Includes with integers
func TestIncludes_Integers(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	// Check if any number is greater than 3
	result := array.Includes(numbers, func(n int) bool {
		return n > 3
	})

	if !result {
		t.Error("Expected true, got false")
	}

	// Check if any number is greater than 10
	result = array.Includes(numbers, func(n int) bool {
		return n > 10
	})

	if result {
		t.Error("Expected false, got true")
	}
}

// TestIncludes_Strings tests Includes with strings
func TestIncludes_Strings(t *testing.T) {
	words := []string{"apple", "banana", "cherry", "date"}

	// Check if any word starts with 'b'
	result := array.Includes(words, func(s string) bool {
		return len(s) > 0 && s[0] == 'b'
	})

	if !result {
		t.Error("Expected true, got false")
	}

	// Check if any word starts with 'z'
	result = array.Includes(words, func(s string) bool {
		return len(s) > 0 && s[0] == 'z'
	})

	if result {
		t.Error("Expected false, got true")
	}
}

// TestIncludes_Structs tests Includes with structs
func TestIncludes_Structs(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 17},
		{Name: "Charlie", Age: 30},
	}

	// Check if any person is under 18
	result := array.Includes(people, func(p Person) bool {
		return p.Age < 18
	})

	if !result {
		t.Error("Expected true (Bob is 17), got false")
	}

	// Check if any person is over 40
	result = array.Includes(people, func(p Person) bool {
		return p.Age > 40
	})

	if result {
		t.Error("Expected false, got true")
	}
}

// TestIncludes_EmptySlice tests Includes with empty slice
func TestIncludes_EmptySlice(t *testing.T) {
	empty := []int{}
	result := array.Includes(empty, func(n int) bool {
		return n > 0
	})

	if result {
		t.Error("Expected false for empty slice, got true")
	}
}

// TestIncludes_NilPredicate tests Includes with nil predicate
func TestIncludes_NilPredicate(t *testing.T) {
	numbers := []int{1, 2, 3}
	result := array.Includes(numbers, nil)

	if result {
		t.Error("Expected false for nil predicate, got true")
	}
}

// TestIncludes_AllMatch tests Includes when all elements match
func TestIncludes_AllMatch(t *testing.T) {
	numbers := []int{2, 4, 6, 8, 10}
	result := array.Includes(numbers, func(n int) bool {
		return n%2 == 0
	})

	if !result {
		t.Error("Expected true (all are even), got false")
	}
}

// TestIncludes_FirstMatch tests Includes returns true on first match
func TestIncludes_FirstMatch(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	callCount := 0

	result := array.Includes(numbers, func(n int) bool {
		callCount++
		return n == 2
	})

	if !result {
		t.Error("Expected true, got false")
	}

	// Should stop at index 1 (second element)
	if callCount != 2 {
		t.Errorf("Expected 2 calls (stops at first match), got %d", callCount)
	}
}

// TestIncludesIndexed_WithIndex tests IncludesIndexed using index
func TestIncludesIndexed_WithIndex(t *testing.T) {
	numbers := []int{10, 20, 30, 40, 50}

	// Check if any element at even index is greater than 25
	result := array.IncludesIndexed(numbers, func(v int, index int, arr []int) bool {
		return index%2 == 0 && v > 25
	})

	if !result {
		t.Error("Expected true (30 at index 2, 50 at index 4), got false")
	}

	// Check if any element at odd index is greater than 100
	result = array.IncludesIndexed(numbers, func(v int, index int, arr []int) bool {
		return index%2 == 1 && v > 100
	})

	if result {
		t.Error("Expected false, got true")
	}
}

// TestIncludesIndexed_WithArray tests IncludesIndexed using array parameter
func TestIncludesIndexed_WithArray(t *testing.T) {
	numbers := []int{5, 10, 15, 20, 25}

	// Check if any element is greater than the first element
	result := array.IncludesIndexed(numbers, func(v int, index int, arr []int) bool {
		return v > arr[0]
	})

	if !result {
		t.Error("Expected true (10, 15, 20, 25 are all > 5), got false")
	}

	// Check if any element is less than the first element
	result = array.IncludesIndexed(numbers, func(v int, index int, arr []int) bool {
		return v < arr[0]
	})

	if result {
		t.Error("Expected false, got true")
	}
}

// TestIncludesIndexed_ComplexPredicate tests complex predicate logic
func TestIncludesIndexed_ComplexPredicate(t *testing.T) {
	products := []Product{
		{ID: 1, Name: "Laptop", Price: 999.99},
		{ID: 2, Name: "Mouse", Price: 25.50},
		{ID: 3, Name: "Keyboard", Price: 75.00},
		{ID: 4, Name: "Monitor", Price: 299.99},
	}

	// Check if any product at odd index has price > 50
	result := array.IncludesIndexed(products, func(p Product, index int, arr []Product) bool {
		return index%2 == 1 && p.Price > 50
	})

	if !result {
		t.Error("Expected true (Monitor at index 3 has price 299.99), got false")
	}
}

// TestIncludesIndexed_NilPredicate tests IncludesIndexed with nil predicate
func TestIncludesIndexed_NilPredicate(t *testing.T) {
	numbers := []int{1, 2, 3}
	result := array.IncludesIndexed(numbers, nil)

	if result {
		t.Error("Expected false for nil predicate, got true")
	}
}

// TestIncludesIndexed_EmptySlice tests IncludesIndexed with empty slice
func TestIncludesIndexed_EmptySlice(t *testing.T) {
	empty := []int{}
	result := array.IncludesIndexed(empty, func(v int, index int, arr []int) bool {
		return v > 0
	})

	if result {
		t.Error("Expected false for empty slice, got true")
	}
}

// BenchmarkIncludes_Integers benchmarks Includes with integers
func BenchmarkIncludes_Integers(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = array.Includes(numbers, func(n int) bool {
			return n > 500
		})
	}
}

// BenchmarkIncludesIndexed_Integers benchmarks IncludesIndexed with integers
func BenchmarkIncludesIndexed_Integers(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = array.IncludesIndexed(numbers, func(v int, index int, arr []int) bool {
			return v > 500 && index%2 == 0
		})
	}
}
