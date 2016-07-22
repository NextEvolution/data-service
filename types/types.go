package types


//type Seller struct {
//	FbUser FbUser
//	PictureUrl string
//	Name string
//	FbToken string
//	//settings. Albums to ignore, groups to watch, keyword, etc.
//	//facebook token
//
//}

type SellerAlbumScan struct {
	Date int `json:"date"`
	Products []Product `json:"products"`
}

type Product struct {
	Album       string `json:"album"`
	Description string `json:"description"`
	Metadata    FbPicture `json:"metadata"`
	SaleEvents  []SaleEvent `json:"sales_events"`
}

type SaleEvent struct {
	Metadata FbComment `json:"metadata"`
	Customer Customer `json:"customer"`
	Date int `json:"date"`
}

type Customer struct {
	Name     string `json:"name"`
	Metadata FbUser `json:"metadata"`
	//ContactInfo
}

type FbUser struct {
	FbId string `json:"fb_id"`
	Name string `json:"name"`
}

type FbPicture struct {
	Height int `json:"height"`
	Width int `json:"width"`
	ImageUrl string `json:"image_url"`
	FbId string `json:"fb_id"`
	CreatedTime int `json:"created_time"`
}

type FbComment struct {
	Text string `json:"text"`
	FbId string `json:"fb_id"`
}