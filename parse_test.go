package addresscn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type data struct {
	src string
	dst Address
}

func TestParseAddress(t *testing.T) {
	var c = NewFromGithub()
	var collection = []data{
		data{
			src: "安徽省合肥市高新区柏堰科技园石楠路13号",
			dst: Address{
				ProvinceCode: "34",
				CityCode:     "3401",
				AreaCode:     "340176",
				Detail:       "柏堰科技园石楠路13号",
			},
		},
	}

	for _, test := range collection {
		v, err := c.ParseAddress(test.src)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, test.dst, v)
	}
}
