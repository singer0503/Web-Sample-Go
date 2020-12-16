//在套件中，測試程式碼必須是 _test.go 結尾
package cal

import (
	"testing"
)

var (
	input = [][]int{
		{1, 2},
		{2, 3},
		{3, 4},
	}
	output = []int{3, 5, 7}
)

func TestAdd(t *testing.T) {
	for i, in := range input {
		result := Add(in[0], in[1])
		t.Log("in =", in, "output should be", output[i], " actual output equal", result)
		if output[i] == result {
			t.Log("Pass")
		} else {
			t.Error("Case", i, "is failed")
		}
	}
}
