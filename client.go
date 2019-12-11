// addresscn can parse address string to standardize China address.
// 初始化阶段如果出错 会直接 panic
package addresscn

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

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
	cityP     map[string][]City // provinceCode-citys
	area      map[string]Area   // code-area
	areaP     map[string][]Area // cityCode-areas
}

// NewFromCOS 从腾讯云对象存储获取数据 用了我的其他库 内网速度快
func NewFromCOS(cos xobj.Client) *Client {
	var client = &Client{
		provider:  "cos",
		cos:       cos,
		province:  make(map[string]string),
		provinceR: make(map[string]string),
		city:      make(map[string]City),
		cityR:     make(map[string]City),
		cityP:     make(map[string][]City),
		area:      make(map[string]Area),
		areaP:     make(map[string][]Area),
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
		// 去除省字再来一遍
		c.provinceR[strings.TrimSuffix(name, "省")] = code
	}
	// 特殊处理一些容易被叫错的省 反正能查出来就行
	c.provinceR["广西"] = "45"
	c.provinceR["广西省"] = "45"
	c.provinceR["宁夏"] = "64"
	c.provinceR["宁夏省"] = "64"
	c.provinceR["新疆"] = "65"
	c.provinceR["新疆省"] = "65"
	c.provinceR["内蒙古"] = "15"
	c.provinceR["内蒙古省"] = "15"
	c.provinceR["西藏"] = "54"
	c.provinceR["西藏省"] = "54"
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
		// 代码字典
		x, ok := c.city[code]
		if ok {
			panic(fmt.Sprintf("duplicate city code %s with name %s and %s", code, name, x))
		}
		c.city[code] = city
		// 名称字典特殊处理直辖市 把 key 由"市辖区"改成直辖市名
		switch code {
		case "1101":
			name = "北京市"
		case "1201":
			name = "天津市"
		case "3101":
			name = "上海市"
		case "5001":
			name = "重庆市"
		}
		// 跳过直辖县
		if strings.Contains(name, "直辖县") {
			continue
		}
		// 名称字典
		y, ok := c.cityR[name]
		if ok {
			panic(fmt.Sprintf("duplicate city name %s with code %s and %s", name, code, y))
		}
		c.cityR[name] = city

		// 存储省-市关系
		if _, ok := c.cityP[line[2]]; !ok {
			c.cityP[line[2]] = make([]City, 0, 99)
		}
		c.cityP[line[2]] = append(c.cityP[line[2]], city)
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

		// 存储市-区关系
		if _, ok := c.areaP[line[2]]; !ok {
			c.areaP[line[2]] = make([]Area, 0, 99)
		}
		c.areaP[line[2]] = append(c.areaP[line[2]], area)
	}
}

// GetPCA load the province city area stream from cos.
func (c *Client) GetPCA() (io.Reader, error) {
	return c.cos.Reader("/division/pca-code.json")
}
