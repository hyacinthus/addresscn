package addresscn

// Province 省/直辖市
type Province struct {
	Code string
	Name string
}

// City 市
type City struct {
	Code         string
	Name         string
	ProvinceCode string
}

// Area 区/县
type Area struct {
	Code         string
	Name         string
	ProvinceCode string
	CityCode     string
}

// item is 映射表元素
type item struct {
	name string
	code string
}

// 映射表 其实是个可排序列表
type table []item

// Add is a shortcut for append
func (t table) Add(name, code string) table {
	return append(t, item{name, code})
}

// Len for sort interface
func (t table) Len() int {
	return len(t)
}

// Swap for sort interface
func (t table) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Less for sort interface
func (t table) Less(i, j int) bool {
	return len(t[i].name) > len(t[j].name)
}
