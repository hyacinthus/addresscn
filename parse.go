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

// ParseCity 从名字解析市
func (c *Client) ParseCity(name string) (City, error) {
	city, ok := c.cityR[name]
	if !ok {
		return city, ErrorNotFound
	}
	return city, nil
}

// ParseArea 从名字解析区县
func (c *Client) ParseArea(name string) (Area, error) {
	area, ok := c.areaR[name]
	if !ok {
		return area, ErrorNotFound
	}
	return area, nil
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

// ParseAddress 解析地址 TODO: 这个问题不少，还需要完善
func (c *Client) ParseAddress(src string) (result Address, err error) {
	cAddr := src
	for k, v := range c.provinceR {
		if strings.HasPrefix(src, k) {
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
				standardName, err := c.ProvinceName(v)
				if err != nil {
					return Address{}, err
				}
				cAddr = strings.TrimPrefix(src, standardName)
				cAddr = strings.TrimPrefix(cAddr, k)
				cAddr = strings.TrimPrefix(cAddr, "自治区")
				cAddr = strings.TrimPrefix(cAddr, "市")
				cAddr = strings.TrimPrefix(cAddr, "省")
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
			standardName, err := c.CityName(v.Code)
			if err != nil {
				return Address{}, err
			}
			aAddr = strings.TrimPrefix(cAddr, standardName)
			aAddr = strings.TrimPrefix(aAddr, k)
			aAddr = strings.TrimPrefix(aAddr, "地区")
			aAddr = strings.TrimPrefix(aAddr, "自治州")
			aAddr = strings.TrimPrefix(aAddr, "市")
			if len(result.CityCode) > 0 {
				break
			}
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
				standardName, err := c.AreaName(v.Code)
				if err != nil {
					return Address{}, err
				}
				src = strings.TrimPrefix(aAddr, standardName)
				src = strings.TrimPrefix(src, area[1])
				result.Detail = strings.TrimPrefix(src, "旗")
				result.Detail = strings.TrimPrefix(result.Detail, "县")
				result.Detail = strings.TrimPrefix(result.Detail, "区")
				break
			}
		}
	}
	if len(result.ProvinceCode) == 0 || len(result.CityCode) == 0 || len(result.AreaCode) == 0 {
		err = ErrorInvalidAddr
	}
	return
}
