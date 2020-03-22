// addresscn can parse address string to standardize China address.
// 初始化阶段如果出错 会直接 panic
package addresscn

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/hyacinthus/x/xobj"
)

var (
	ErrorNotFound = errors.New("address code not found")
)

// Client 地址解析客户端
// 请调用方保证初始化完成后再使用
// 这些数据在初始化阶段完成写入，提供服务后不再写入，所以只有并行读取，是线程安全的。
type Client struct {
	provider string // 数据来源 github(default)/http/cos
	// url       string // http 模式时 文件的地址前缀 包含最后的斜线
	cos       xobj.Client
	province  map[string]string // code:name
	provinceR table             // name,code array
	city      map[string]City   // code:city
	cityR     table             // name,code array
	area      map[string]Area   // code:area
	areaR     map[string]table  // city: name,code array
}

// NewFromCOS 从腾讯云对象存储获取数据 用了我的其他库 内网速度快
func NewFromCOS(cos xobj.Client) *Client {
	var client = &Client{
		provider: "cos",
		cos:      cos,
		province: make(map[string]string),
		city:     make(map[string]City),
		area:     make(map[string]Area),
		areaR:    make(map[string]table),
		// 此处不用初始化几个反向映射，后面会新建并赋值
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
	var t = make(table, 0) // 反向映射
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

		// 加入反向映射
		t = t.Add(name, code)
		// 简称反向映射
		if strings.HasSuffix(name, "省") {
			t = t.Add(strings.TrimSuffix(name, "省"), code)
		}
		if strings.HasSuffix(name, "市") {
			t = t.Add(strings.TrimSuffix(name, "市"), code)
		}
	}
	// 特殊处理一些容易被叫错的省
	t = t.Add("广西", "45")
	t = t.Add("广西省", "45")
	t = t.Add("广西自治区", "45")
	t = t.Add("宁夏", "64")
	t = t.Add("宁夏省", "64")
	t = t.Add("宁夏自治区", "64")
	t = t.Add("新疆", "65")
	t = t.Add("新疆省", "65")
	t = t.Add("新疆自治区", "65")
	t = t.Add("内蒙古", "15")
	t = t.Add("内蒙古省", "15")
	t = t.Add("西藏", "54")
	t = t.Add("西藏省", "54")
	// 排序
	sort.Sort(t)
	c.provinceR = t
}

// LoadCity load the city data from a io reader.
func (c *Client) LoadCity(r io.Reader) {
	var t = make(table, 0) // 反向映射
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

		// 开始处理反向映射
		// 先跳过几种情况 直辖市的市辖区和县 直辖县 在市级都没名字
		if name == "市辖区" || name == "省直辖县级行政区划" || name == "县" {
			continue
		}
		// 剩下的全加进来
		t = t.Add(name, code)
		// 然后去掉通用后缀
		if strings.HasSuffix(name, "市") {
			t = t.Add(strings.TrimSuffix(name, "市"), code)
		}
		if strings.HasSuffix(name, "地区") {
			t = t.Add(strings.TrimSuffix(name, "地区"), code)
			t = t.Add(strings.TrimSuffix(name, "地区")+"市", code)
		}
		if strings.HasSuffix(name, "盟") {
			t = t.Add(strings.TrimSuffix(name, "盟")+"市", code)
		}
	}
	// 然后特殊处理不规范的自治州
	t = t.Add("延边自治州", "2224")
	t = t.Add("延边州", "2224")
	t = t.Add("延边市", "2224")
	t = t.Add("延边", "2224")
	t = t.Add("恩施自治州", "4228")
	t = t.Add("恩施州", "4228")
	t = t.Add("恩施市", "4228")
	t = t.Add("恩施", "4228")
	t = t.Add("湘西自治州", "4331")
	t = t.Add("湘西州", "4331")
	t = t.Add("湘西市", "4331")
	t = t.Add("湘西", "4331")
	t = t.Add("阿坝自治州", "5132")
	t = t.Add("阿坝州", "5132")
	t = t.Add("阿坝市", "5132")
	t = t.Add("阿坝", "5132")
	t = t.Add("甘孜自治州", "5133")
	t = t.Add("甘孜州", "5133")
	t = t.Add("甘孜市", "5133")
	t = t.Add("甘孜", "5133")
	t = t.Add("凉山自治州", "5134")
	t = t.Add("凉山州", "5134")
	t = t.Add("凉山市", "5134")
	t = t.Add("凉山", "5134")
	t = t.Add("黔西南自治州", "5223")
	t = t.Add("黔西南州", "5223")
	t = t.Add("黔西南市", "5223")
	t = t.Add("黔西南", "5223")
	t = t.Add("黔东南自治州", "5226")
	t = t.Add("黔东南州", "5226")
	t = t.Add("黔东南市", "5226")
	t = t.Add("黔东南", "5226")
	t = t.Add("黔南自治州", "5227")
	t = t.Add("黔南州", "5227")
	t = t.Add("黔南市", "5227")
	t = t.Add("黔南", "5227")
	t = t.Add("楚雄自治州", "5323")
	t = t.Add("楚雄州", "5323")
	t = t.Add("楚雄市", "5323")
	t = t.Add("楚雄", "5323")
	t = t.Add("红河自治州", "5325")
	t = t.Add("红河州", "5325")
	t = t.Add("红河市", "5325")
	t = t.Add("红河", "5325")
	t = t.Add("文山自治州", "5326")
	t = t.Add("文山州", "5326")
	t = t.Add("文山市", "5326")
	t = t.Add("文山", "5326")
	t = t.Add("西双版纳自治州", "5328")
	t = t.Add("西双版纳州", "5328")
	t = t.Add("西双版纳市", "5328")
	t = t.Add("西双版纳", "5328")
	t = t.Add("大理自治州", "5329")
	t = t.Add("大理州", "5329")
	t = t.Add("大理市", "5329")
	t = t.Add("大理", "5329")
	t = t.Add("德宏自治州", "5331")
	t = t.Add("德宏州", "5331")
	t = t.Add("德宏市", "5331")
	t = t.Add("德宏", "5331")
	t = t.Add("怒江自治州", "5333")
	t = t.Add("怒江州", "5333")
	t = t.Add("怒江市", "5333")
	t = t.Add("怒江", "5333")
	t = t.Add("迪庆自治州", "5334")
	t = t.Add("迪庆州", "5334")
	t = t.Add("迪庆市", "5334")
	t = t.Add("迪庆", "5334")
	t = t.Add("临夏自治州", "6229")
	t = t.Add("临夏州", "6229")
	t = t.Add("临夏市", "6229")
	t = t.Add("临夏", "6229")
	t = t.Add("甘南自治州", "6230")
	t = t.Add("甘南州", "6230")
	t = t.Add("甘南市", "6230")
	t = t.Add("甘南", "6230")
	t = t.Add("海北自治州", "6322")
	t = t.Add("海北州", "6322")
	t = t.Add("海北市", "6322")
	t = t.Add("海北", "6322")
	t = t.Add("黄南自治州", "6323")
	t = t.Add("黄南州", "6323")
	t = t.Add("黄南市", "6323")
	t = t.Add("黄南", "6323")
	t = t.Add("海南自治州", "6325")
	t = t.Add("海南州", "6325")
	t = t.Add("海南市", "6325")
	t = t.Add("果洛自治州", "6326")
	t = t.Add("果洛州", "6326")
	t = t.Add("果洛市", "6326")
	t = t.Add("果洛", "6326")
	t = t.Add("玉树自治州", "6327")
	t = t.Add("玉树州", "6327")
	t = t.Add("玉树市", "6327")
	t = t.Add("玉树", "6327")
	t = t.Add("海西自治州", "6328")
	t = t.Add("海西州", "6328")
	t = t.Add("海西市", "6328")
	t = t.Add("海西", "6328")
	t = t.Add("昌吉自治州", "6523")
	t = t.Add("昌吉州", "6523")
	t = t.Add("昌吉市", "6523")
	t = t.Add("昌吉", "6523")
	t = t.Add("博尔塔拉自治州", "6527")
	t = t.Add("博尔塔拉州", "6527")
	t = t.Add("博尔塔拉市", "6527")
	t = t.Add("博尔塔拉", "6527")
	t = t.Add("巴音郭楞自治州", "6528")
	t = t.Add("巴音郭楞州", "6528")
	t = t.Add("巴音郭楞市", "6528")
	t = t.Add("巴音郭楞", "6528")
	// t = t.Add("克孜勒苏柯尔克孜自治州", "6530")
	t = t.Add("克州", "6530")
	t = t.Add("伊犁自治州", "6540")
	t = t.Add("伊犁州", "6540")
	t = t.Add("伊犁市", "6540")
	t = t.Add("伊犁", "6540")

	// 最后排序保存
	sort.Sort(t)
	c.cityR = t
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
		city := line[2]
		area := Area{
			Code:         code,
			Name:         name,
			CityCode:     city,
			ProvinceCode: line[3],
		}
		x, ok := c.area[code]
		if ok {
			panic(fmt.Sprintf("duplicate area code %s with name %s and %s", code, name, x))
		}
		c.area[code] = area

		// 开始保存反向映射
		if c.areaR[city] == nil {
			c.areaR[city] = make(table, 0)
		}
		c.areaR[city] = c.areaR[city].Add(name, code)
		if strings.HasSuffix(name, "经济开发区") || strings.HasSuffix(name, "经济技术开发区") {
			c.areaR[city] = c.areaR[city].Add("经济开发区", code)
			c.areaR[city] = c.areaR[city].Add("经开区", code)
		}
		if strings.HasSuffix(name, "高新技术产业开发区") {
			c.areaR[city] = c.areaR[city].Add("高新区", code)
		}
	}
	// 排序每个城市的反向映射
	for i := range c.areaR {
		sort.Sort(c.areaR[i])
	}
}

// GetPCA load the province city area stream from cos.
func (c *Client) GetPCA() (io.ReadCloser, error) {
	return c.cos.Reader("/division/pca-code.json")
}
