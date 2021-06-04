package main

import (
	"sina/function"
	"time"
)

func main() {
	//先获取信息
	//function.GetSinaLuckingApi()
	//function.GetSinaLucking()
	//ticker1 := time.NewTicker(time.Minute * 1)
	//ticker2 := time.NewTicker(time.Minute * 2)
	//ticker3 := time.NewTicker(time.Hour * 24)
	//ticker4 := time.NewTicker(time.Minute * 30)
	//for {
	//	select {
	//	case <-ticker1.C:
	//		function.GetLuckSearchApi()
	//	case <-ticker2.C:
	//		function.GetSinaLuckingApi()
	//	case <-ticker3.C:
	//		function.GetSinaLucking()
	//	case <-ticker4.C:
	//		function.FollowSet()
	//		function.HuaTiZhuanFa()
	//		function.LikeSet()
	//	}
	//}
	for {
		time.Sleep(5*time.Minute)
		function.GetLuckSearchApi()
		function.FollowSet()
		function.HuaTiZhuanFa()
		function.LikeSet()
	}

}
