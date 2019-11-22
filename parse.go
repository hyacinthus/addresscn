package addresscn

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

// GetCityName 获取县区名称
func (c *Client) GetAreaName(code string) (string, error) {
	if area, ok := c.area[code]; ok {
		return area.Name, nil
	}
	return "", ErrorNotFound
}
