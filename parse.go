package addresscn

import "strings"

// ParseProvince 从名字解析省份代码
func (c *Client) ParseProvince(name string) (string, error) {
	p, ok := c.provinceR[name]
	if !ok {
		return "", ErrorNotFound
	}
	return p, nil
}

// ParseCity 从名字解析市代码
func (c *Client) ParseCity(name string) (City, error) {
	p, ok := c.cityR[name]
	if !ok {
		return p, ErrorNotFound
	}
	return p, nil
}

// FindProvinces 获取省份
func (c *Client) FindProvinces() map[string]string {
	return c.province
}

// FindCitys 获取城市
func (c *Client) FindCitys(code string) (map[string]string, error) {
	rows := make(map[string]string)
	if citys, ok := c.cityP[code]; ok {
		for _, v := range citys {
			rows[v.Code] = v.Name
		}
		return rows, nil
	}
	return nil, ErrorNotFound
}

// FindAreas 获取县区
func (c *Client) FindAreas(code string) (map[string]string, error) {
	rows := make(map[string]string)
	if areas, ok := c.areaP[code]; ok {
		for _, v := range areas {
			rows[v.Code] = v.Name
		}
		return rows, nil
	}
	return nil, ErrorNotFound
}

// GetProvinceName 获取省份名称
func (c *Client) GetProvinceName(code string) (string, error) {
	if name, ok := c.province[code]; ok {
		return name, nil
	}
	return "", ErrorNotFound
}

// GetCityName 获取市名称
func (c *Client) GetCityName(code string) (string, error) {
	if city, ok := c.city[code]; ok {
		return city.Name, nil
	}
	return "", ErrorNotFound
}

// GetAreaName 获取县区名称
func (c *Client) GetAreaName(code string) (string, error) {
	if area, ok := c.area[code]; ok {
		return area.Name, nil
	}
	return "", ErrorNotFound
}

// ParseAddress 解析地址
func (c *Client) ParseAddress(addr string) (result Address, err error) {
	cAddr := addr
	for k, v := range c.provinceR {
		if strings.HasPrefix(addr, k) {
			result.ProvinceCode = v
			switch k {
			case "北京":
			case "北京市":
			case "天津":
			case "天津市":
			case "重庆":
			case "重庆市":
			case "上海":
			case "上海市":
			default:
				cAddr = strings.TrimPrefix(addr, k)
				cAddr = strings.TrimPrefix(cAddr, "省")
				cAddr = strings.TrimPrefix(cAddr, "市")
				cAddr = strings.TrimPrefix(cAddr, "自治区")
			}
			break
		}
	}
	aAddr := cAddr
	for k, v := range c.cityR {
		if strings.HasPrefix(cAddr, k) {
			if len(result.ProvinceCode) == 0 {
				result.ProvinceCode = v.ProvinceCode
			} else if result.ProvinceCode != v.ProvinceCode {
				// 省ID对不上则直接返回，只解析到省
				err = ErrorInvalidAddr
				return
			}
			result.CityCode = v.Code
			aAddr = strings.TrimPrefix(cAddr, k)
			aAddr = strings.TrimPrefix(aAddr, "市")
			aAddr = strings.TrimPrefix(aAddr, "自治州")
			aAddr = strings.TrimPrefix(aAddr, "地区")
		}
	}
	// 若省市找不到，直接返回，必须包含省市
	if len(result.ProvinceCode) == 0 || len(result.CityCode) == 0 {
		err = ErrorInvalidAddr
		return
	}
	for k, v := range c.areaR {
		area := strings.Split(k, "-")
		if strings.HasPrefix(aAddr, area[1]) {
			if result.ProvinceCode == v.ProvinceCode && result.CityCode == v.CityCode {
				result.AreaCode = v.Code
				result.Detail = strings.TrimPrefix(aAddr, area[1])
				return
			}
		}
	}
	if len(result.ProvinceCode) == 0 || len(result.CityCode) == 0 || len(result.AreaCode) == 0 {
		err = ErrorInvalidAddr
	}
	return
}
