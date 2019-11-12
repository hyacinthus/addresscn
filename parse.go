package addresscn

// ParseProvince 从名字解析省份代码
func (c *Client) ParseProvince(name string) (string, error) {
	p, ok := c.provinceR[name]
	if !ok {
		return "", ErrorNotFound
	}
	return p, nil
}
