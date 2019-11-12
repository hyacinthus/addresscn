// addresscn can parse address string to standardize China address.
// 初始化阶段如果出错 会直接 panic
package addresscn

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/hyacinthus/x/xobj"
)

var (
	ErrorNotFound = errors.New("the key not found")
)

// Client 地址解析客户端
type Client struct {
	provider string // 数据来源 github(default)/http/cos
	//url       string // http 模式时 文件的地址前缀 包含最后的斜线
	cos       xobj.Client
	province  map[string]string // code-name
	provinceR map[string]string // name-code
	city      map[string]City   // code-city
	cityR     map[string]City   // name-city
	area      map[string]Area   // code-area
}

// NewFromCOS 从腾讯云对象存储获取数据 用了我的其他库 内网速度快
func NewFromCOS(cos xobj.Client) *Client {
	var client = &Client{
		provider: "cos",
		cos:      cos,
	}
	p, err := cos.Reader("/division/provinces.csv")
	if err != nil {
		panic(err)
	}
	client.LoadProvince(p)
	err = p.Close()
	if err != nil {
		panic(err)
	}
	c, err := cos.Reader("/division/cities.csv")
	if err != nil {
		panic(err)
	}
	client.LoadCity(c)
	err = c.Close()
	if err != nil {
		panic(err)
	}
	a, err := cos.Reader("/division/areas.csv")
	if err != nil {
		panic(err)
	}
	client.LoadArea(a)
	err = a.Close()
	if err != nil {
		panic(err)
	}
	return client
}

// LoadProvince load the province data from a io reader.
func (c *Client) LoadProvince(r io.Reader) {
	pr := csv.NewReader(r)
	// skip the title line
	_, err := pr.Read()
	if err != nil {
		panic(err)
	}
	// parse every line
	for {
		line, err := pr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		code := line[0]
		name := line[1]
		x, ok := c.province[code]
		if ok {
			panic(fmt.Sprintf("duplicate province code %s with name %s and %s", code, name, x))
		}
		c.province[code] = name

		y, ok := c.provinceR[name]
		if ok {
			panic(fmt.Sprintf("duplicate province name %s with code %s and %s", name, code, y))
		}
		c.provinceR[name] = code
	}
}

// LoadCity load the city data from a io reader.
func (c *Client) LoadCity(r io.Reader) {
	cr := csv.NewReader(r)
	// skip the title line
	_, err := cr.Read()
	if err != nil {
		panic(err)
	}
	// parse every line
	for {
		line, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		code := line[0]
		name := line[1]
		city := City{
			Code:         code,
			Name:         name,
			ProvinceCode: line[2],
		}
		x, ok := c.city[code]
		if ok {
			panic(fmt.Sprintf("duplicate city code %s with name %s and %s", code, name, x))
		}
		c.city[code] = city
		y, ok := c.cityR[name]
		if ok {
			panic(fmt.Sprintf("duplicate city name %s with code %s and %s", name, code, y))
		}
		c.cityR[name] = city
	}
}

// LoadArea load the area data from a io reader.
func (c *Client) LoadArea(r io.Reader) {
	ar := csv.NewReader(r)
	// skip the title line
	_, err := ar.Read()
	if err != nil {
		panic(err)
	}
	// parse every line
	for {
		line, err := ar.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		code := line[0]
		name := line[1]
		area := Area{
			Code:         code,
			Name:         name,
			CityCode:     line[2],
			ProvinceCode: line[3],
		}
		x, ok := c.area[code]
		if ok {
			panic(fmt.Sprintf("duplicate area code %s with name %s and %s", code, name, x))
		}
		c.area[code] = area
	}
}
