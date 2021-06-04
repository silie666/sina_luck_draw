package respdata

type SinaLuckData struct {
	Html string `json:"html"`
}

type SinaLuckDataApi struct {
	Data struct{
		Cards []struct{
			CardGroup []struct{
				Mblog struct{
					Mid string `json:"mid"`
					Text string `json:"text"`
					User struct{
						Id int `json:"id"`
						ScreenName string `json:"screen_name"`
					}
				}`json:"mblog"`
			} `json:"card_group"`
		}`json:"cards"`
	}`json:"data"`

}


type SinaLuckSearchApi struct {
	Data struct{
		Cards []struct{
			CardType int `json:"card_type"`
			Mblog struct{
				Mid string `json:"mid"`
				Text string `json:"text"`
				User struct{
					Id int `json:"id"`
					ScreenName string `json:"screen_name"`
				}
			}`json:"mblog"`
		}`json:"cards"`
	}`json:"data"`

}


type SinaCode struct {
	Code string `json:"code"`
	Msg string `json:"msg"`
}
