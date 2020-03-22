package addresscn

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSort 测试按照名称长短排序
func TestSort(t *testing.T) {
	var src = table{
		item{"黑龙江省", "2"},
		item{"陕西省", "1"},
		item{"宁夏回族自治区", "3"},
	}
	var dst = table{
		item{"宁夏回族自治区", "3"},
		item{"黑龙江省", "2"},
		item{"陕西省", "1"},
	}
	sort.Sort(src)
	assert.Equal(t, dst, src)
}
