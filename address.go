package addresscn

// Address 通信地址 可被 gorm 公用
type Address struct {
	ProvinceCode string `json:"province_code" gorm:"type:varchar(2)"` // 省 必填
	CityCode     string `json:"city_code" gorm:"type:varchar(4)"`     // 市 必填
	AreaCode     string `json:"area_code" gorm:"type:varchar(6)"`     // 区 必填
	Detail       string `json:"detail" gorm:"type:varchar(255)"`      // 街道一下具体地址 必填
	ZipCode      string `json:"zip_code" gorm:"type:varchar(6)"`      // 邮编 非必填
}
