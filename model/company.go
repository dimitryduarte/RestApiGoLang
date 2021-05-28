package model

type Company struct {
	IdCompany   int64  `json:"idcompany" gorm:"column:id;primaryKey;autoIncrement"`
	CompanyName string `json:"companyname" gorm:"column:companyname"`
	Name        string `json:"name" gorm:"column:name"`
	Cnpj        string `json:"cnpj" gorm:"column:cnpj"`
	TaxesId     int64  `json:"taxesid" gorm:"column:TaxesId"`
	Taxes       Taxes  `json:"taxes" gorm:"foreignKey:TaxesId;References:id"`
}

type Taxes struct {
	TaxesId int64   `json:"taxesid" gorm:"column:id;primaryKey;autoIncrement"`
	Taxa1   float64 `json:"taxa1" gorm:"column:taxa1"`
	Taxa2   float64 `json:"taxa2" gorm:"column:taxa2"`
	Taxa3   float64 `json:"taxa3" gorm:"column:taxa3"`
}
