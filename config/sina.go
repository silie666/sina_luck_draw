package config


func GetSinaUrl() map[string]interface{} {
	// 初始化数据库配置map
	urlConfig := make(map[string]interface{})


	urlConfig["LUCKING"] = "https://weibo.com/p/100808557b69009a8ef6588f9124fe9c30d36c/super_index"
	urlConfig["LUCKING_TIME"] = "https://m.weibo.cn/api/container/getIndex?jumpfrom=weibocom&containerid=100808557b69009a8ef6588f9124fe9c30d36c_-_sort_time"
	urlConfig["LUCKING_SEARCH"] = "https://m.weibo.cn/api/container/getIndex?containerid=100103type%3D1%26q%3D@%E5%BE%AE%E5%8D%9A%E6%8A%BD%E5%A5%96%E5%B9%B3%E5%8F%B0&page_type=searchall"  //关键词1
	urlConfig["LUCKING_SEARCH_ZHUANFA"] = "https://m.weibo.cn/api/container/getIndex?containerid=100103type%3D1%26q%3D%25E8%25BD%25AC%25E5%258F%2591%25E6%258A%25BD%25E5%25A5%2596&page_type=searchall" //关键词2
	urlConfig["LUCKING_SEARCH_XIANGQING"] = "https://m.weibo.cn/api/container/getIndex?containerid=100103type%3D1%26q%3D%E6%8A%BD%E5%A5%96%E8%AF%A6%E6%83%85&page_type=searchall" //关键词3
	urlConfig["LUCKING_STATUS"] = "https://lottery.media.weibo.com/lottery/h5/history/list?mid="   //查看是否存在页面


	urlConfig["REFERER"] = "https://m.weibo.cn/p/100808557b69009a8ef6588f9124fe9c30d36c/super_index?jumpfrom=weibocom"
	urlConfig["PDETAIL"] = "100808557b69009a8ef6588f9124fe9c30d36c"
	//urlConfig["PAGE_ID"] = "page_100808_super_index"
	urlConfig["LOCATION"] = "100808557b69009a8ef6588f9124fe9c30d36c"

	urlConfig["SUB"] = "_2A25NhUDTDeRhGeNL7VYZ8C3EyzuIHXVuhmCbrDV8PUJbkNANLWn8kW1NSOd6D2se5X_gGoyofenFUhb13_u6Casd"  //微博的sub，再cookie查找
	urlConfig["UID"] = "5564801254"  //自己的uid


	urlConfig["COMMENT_URL"] = "https://weibo.com/aj/v6/comment/add"
	urlConfig["FOLLOW_URL"] = "https://weibo.com/aj/f/followed"
	urlConfig["LIKE_URL"] = "https://weibo.com/aj/v6/like/add"
	urlConfig["PAGE_ID"] = "page_100505_home"
	urlConfig["ZHUANFA_URL"] = "https://weibo.com/aj/v6/mblog/forward"
	urlConfig["PAGES_ID"] = "page_100606_home"



	return urlConfig
}
