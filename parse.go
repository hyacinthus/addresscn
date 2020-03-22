package addresscn

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// ParseProvince 从名字解析省份代码
func (c *Client) ParseProvince(name string) (string, error) {
	for _, item := range c.provinceR {
		if name == item.name {
			return item.code, nil
		}
	}
	return "", ErrorNotFound
}

// ParseCity 从名字解析市
func (c *Client) ParseCity(name string) (City, error) {
	for _, item := range c.cityR {
		if name == item.name {
			return c.city[item.code], nil
		}
	}
	return City{}, ErrorNotFound
}

// ProvinceName 获取省份名称
func (c *Client) ProvinceName(code string) (string, error) {
	if name, ok := c.province[code]; ok {
		return name, nil
	}
	return "", ErrorNotFound
}

// CityName 获取市名称
func (c *Client) CityName(code string) (string, error) {
	if city, ok := c.city[code]; ok {
		return city.Name, nil
	}
	return "", ErrorNotFound
}

// AreaName 获取县区名称
func (c *Client) AreaName(code string) (string, error) {
	if area, ok := c.area[code]; ok {
		return area.Name, nil
	}
	return "", ErrorNotFound
}

// MustParseAddress 解析地址 报告任何错误
func (c *Client) MustParseAddress(src string) (Address, error) {
	addr, err := c.ParseAddress(src)
	if err != nil {
		return addr, err
	}
	if addr.AreaCode == "" || addr.Detail == "" {
		return addr, errors.New("解析地址失败：" + src)
	}
	return addr, nil
}

// ParseAddress 解析地址 省市解析成功就算成功 忽略区解析失败的错误
func (c *Client) ParseAddress(src string) (Address, error) {
	var addr Address
	var err = errors.New("解析地址失败：" + src)
	if utf8.RuneCountInString(src) < 6 {
		return addr, err
	}
	// 去掉空格和标点
	cur := strings.ReplaceAll(src, " ", "")
	cur = strings.ReplaceAll(cur, ",", "")
	cur = strings.ReplaceAll(cur, "，", "")
	cur = strings.ReplaceAll(cur, "。", "")
	cur = strings.TrimPrefix(cur, "中国")
	// 开始解析省份
	for _, pr := range c.provinceR {
		if strings.HasPrefix(cur, pr.name) {
			addr.ProvinceCode = pr.code
			cur = strings.TrimPrefix(cur, pr.name)
			// 如果是直辖市，直接确定市级code，后续不再解析市，另外兼容写两遍市和写市辖区的情况
			if pr.code == "11" {
				addr.CityCode = "1101"
				cur = strings.TrimPrefix(cur, pr.name)
				cur = strings.TrimPrefix(cur, "市辖区")
			}
			if pr.code == "12" {
				addr.CityCode = "1201"
				cur = strings.TrimPrefix(cur, pr.name)
				cur = strings.TrimPrefix(cur, "市辖区")
			}
			if pr.code == "31" {
				addr.CityCode = "3101"
				cur = strings.TrimPrefix(cur, pr.name)
				cur = strings.TrimPrefix(cur, "市辖区")
			}
			if pr.code == "50" {
				addr.CityCode = "5001"
				cur = strings.TrimPrefix(cur, pr.name)
				cur = strings.TrimPrefix(cur, "市辖区")
			}
			// 解析到了就跳出省份解析
			break
		}
	}
	// 开始解析市
	if addr.CityCode == "" {
		for _, ci := range c.cityR {
			if strings.HasPrefix(cur, ci.name) {
				addr.CityCode = ci.code
				cur = strings.TrimPrefix(cur, ci.name)
			}
		}
	}
	// 开始解析区 因为要用到城市 所以城市如果没有解析成功则失败退出
	if addr.CityCode == "" {
		addr.Detail = cur
		return addr, err
	}
	for _, ar := range c.areaR[addr.CityCode] {
		if strings.HasPrefix(cur, ar.name) {
			addr.AreaCode = ar.code
			cur = strings.TrimPrefix(cur, ar.name)
		}
	}
	// 剩下的是详情
	addr.Detail = cur

	// 能到这里省市肯定解析成功了，区没有解析成功不报错
	return addr, nil
}
