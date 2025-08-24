package eval

import (
	"math/big"
	"testing"

	"src.elv.sh/pkg/eval/vals"
)

// TestOptimizedIntegerArithmetic tests the fast path for pure integer operations
func TestOptimizedIntegerArithmetic(t *testing.T) {
	tests := []struct {
		name     string
		op       func(...vals.Num) vals.Num
		args     []vals.Num
		expected vals.Num
	}{
		// Addition tests
		{"add_pure_ints", add, []vals.Num{1, 2, 3, 4}, 10},
		{"add_single_int", add, []vals.Num{42}, 42},
		{"add_empty", add, []vals.Num{}, 0},
		{"add_zero", add, []vals.Num{0, 5, 0}, 5},
		{"add_negative", add, []vals.Num{-3, 7, -2}, 2},
		
		// Multiplication tests  
		{"mul_pure_ints", mul, []vals.Num{2, 3, 4}, 24},
		{"mul_single_int", mul, []vals.Num{7}, 7},
		{"mul_empty", mul, []vals.Num{}, 1},
		{"mul_with_zero", mul, []vals.Num{5, 0, 3}, 0},
		{"mul_negative", mul, []vals.Num{-2, 3, -4}, 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.op(tt.args...)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
			// Verify the result is still an int (not promoted to big.Int)
			if _, ok := result.(int); !ok && tt.expected != 0 {
				t.Errorf("Expected int type, got %T", result)
			}
		})
	}
}

// TestMixedTypeHandling tests that mixed types still work correctly 
func TestMixedTypeHandling(t *testing.T) {
	bigInt := big.NewInt(1000000000000)
	bigRat := big.NewRat(3, 2) // 1.5
	
	tests := []struct {
		name string
		op   func(...vals.Num) vals.Num
		args []vals.Num
		check func(vals.Num) bool
	}{
		{
			"add_mixed_int_bigint", 
			add, 
			[]vals.Num{5, bigInt}, 
			func(result vals.Num) bool {
				// vals.NormalizeBigInt converts back to int if it fits
				switch result := result.(type) {
				case int:
					return result == 1000000000005
				case *big.Int:
					return result.Cmp(big.NewInt(1000000000005)) == 0
				default:
					return false
				}
			},
		},
		{
			"add_mixed_int_rat", 
			add,
			[]vals.Num{1, bigRat},
			func(result vals.Num) bool {
				br, ok := result.(*big.Rat)
				expected := big.NewRat(5, 2) // 1 + 3/2 = 5/2
				return ok && br.Cmp(expected) == 0
			},
		},
		{
			"add_mixed_int_float", 
			add,
			[]vals.Num{3, 2.5},
			func(result vals.Num) bool {
				f, ok := result.(float64)
				return ok && f == 5.5
			},
		},
		{
			"mul_mixed_int_float",
			mul,
			[]vals.Num{4, 2.5},
			func(result vals.Num) bool {
				f, ok := result.(float64)
				return ok && f == 10.0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.op(tt.args...)
			if !tt.check(result) {
				t.Errorf("Check failed for result: %v (type %T)", result, result)
			}
		})
	}
}

// TestOverflowHandling tests that integer overflow is properly detected
func TestOverflowHandling(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	
	tests := []struct {
		name     string
		op       func(...vals.Num) vals.Num  
		args     []vals.Num
		checkType string // "bigint" or "int"
	}{
		{
			"add_overflow_positive",
			add,
			[]vals.Num{maxInt, 1},
			"bigint", // should overflow to big.Int
		},
		{
			"add_overflow_negative", 
			add,
			[]vals.Num{minInt, -1},
			"bigint", // should overflow to big.Int
		},
		{
			"mul_overflow",
			mul, 
			[]vals.Num{maxInt, 2},
			"bigint", // should overflow to big.Int
		},
		{
			"add_no_overflow",
			add,
			[]vals.Num{100, 200},
			"int", // should stay int
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.op(tt.args...)
			
			switch tt.checkType {
			case "int":
				if _, ok := result.(int); !ok {
					t.Errorf("Expected int, got %T", result)
				}
			case "bigint":
				if _, ok := result.(*big.Int); !ok {
					t.Errorf("Expected *big.Int, got %T", result)
				}
			}
		})
	}
}

// TestSafeIntArithmetic tests the helper functions
func TestSafeIntArithmetic(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	
	// Test safeIntAdd
	t.Run("safeIntAdd", func(t *testing.T) {
		// Normal case
		if result, ok := safeIntAdd(5, 3); !ok || result != 8 {
			t.Errorf("Expected (8, true), got (%d, %v)", result, ok)
		}
		
		// Overflow case
		if _, ok := safeIntAdd(maxInt, 1); ok {
			t.Error("Expected overflow detection to return false")
		}
		
		// Underflow case
		if _, ok := safeIntAdd(minInt, -1); ok {
			t.Error("Expected underflow detection to return false")
		}
	})
	
	// Test safeIntMul
	t.Run("safeIntMul", func(t *testing.T) {
		// Normal case
		if result, ok := safeIntMul(6, 7); !ok || result != 42 {
			t.Errorf("Expected (42, true), got (%d, %v)", result, ok)
		}
		
		// Zero case
		if result, ok := safeIntMul(0, 12345); !ok || result != 0 {
			t.Errorf("Expected (0, true), got (%d, %v)", result, ok)
		}
		
		// Overflow case
		if _, ok := safeIntMul(maxInt, 2); ok {
			t.Error("Expected overflow detection to return false")
		}
	})
}

// BenchmarkArithmeticOperations provides performance comparisons
func BenchmarkArithmeticOperations(b *testing.B) {
	intArgs := []vals.Num{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	mixedArgs := []vals.Num{1, 2, 3, 4, big.NewInt(5), 6, 7, 8, 9, 10}
	
	b.Run("add_pure_ints", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			add(intArgs...)
		}
	})
	
	b.Run("add_mixed_types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			add(mixedArgs...)
		}
	})
	
	b.Run("mul_pure_ints", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mul(intArgs...)
		}
	})
	
	b.Run("mul_mixed_types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mul(mixedArgs...)
		}
	})
}