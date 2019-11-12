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
