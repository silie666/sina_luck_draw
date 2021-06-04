package function

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"sina/config"
	"sina/logger"
	"sina/model"
	"sina/respdata"
	"strconv"
	"strings"
	"time"
)

func GetSinaLucking() {
	var config = config.GetSinaUrl()
	c := colly.NewCollector()

	c.OnResponse(func(response *colly.Response) {
		compile := regexp.MustCompile(`<script>FM\.view\((\{"ns":"pl\.content\.homeFeed\.index".*\})\)<\/script>`)
		submatch := compile.FindAllSubmatch(response.Body,-1)
		var SinaLuck respdata.SinaLuckData
		json.Unmarshal(submatch[0][1],&SinaLuck)

		//ioutil.WriteFile("c.html",[]byte(SinaLuck.Html),0777)


		str1complie := regexp.MustCompile(`(<div\s+tbinfo="ouid=\d*" class="WB_cardwrap WB_feed_type S_bg2 WB_feed_vipcover WB_feed_like"\s*mid="\d*"  action-type="feed_list_item" diss-data="filter_actionlog=">(?s:.*?)<div node-type="feed_list_repeat" class="WB_feed_repeat S_bg1" style="display:none;"><\/div>\s*<\/div>)`)
		str1submatch := str1complie.FindAllSubmatch([]byte(SinaLuck.Html),-1)


		for i:=0;i<len(str1submatch);i++ {
			info_complie := regexp.MustCompile(`<div\s+tbinfo="ouid=(\d*)" class="WB_cardwrap WB_feed_type S_bg2 WB_feed_vipcover WB_feed_like" mid="(\d*)"  action-type="feed_list_item" diss-data="filter_actionlog=">`)
			info_submatch := info_complie.FindAllSubmatch(str1submatch[i][1],-1)
			/*判断是否在抽奖详情*/
			is_luckling_url := config["LUCKING_STATUS"].(string)+string(info_submatch[0][2])
			req, _ := http.NewRequest("GET", is_luckling_url, nil)
			req.Header.Set("cookie", "SUB="+config["SUB"].(string))
			c := http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
				Timeout: 30 * time.Second,
			}
			r, err := c.Do(req)
			if err != nil {
				logger.LoggerToFile("错误："+err.Error())
				return
			}
			resp,err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			lucking_complie := regexp.MustCompile(`window\.__DATA__\s=\s({(?s:.*?)})\s\|\|`)
			lucking_info := lucking_complie.FindAllSubmatch(resp,-1)
			data := struct {
				Total int `json:"total"`
				List []struct{
					Time string `json:"time"`
				}`json:"list"`
				Weibo struct{
					User struct{
						Name string `json:"name"`
					} `json:"user"`
				}	`json:"weibo"`
			}{}
			if len(lucking_info) == 0 {
				lucking_info[0][1] = []byte(strings.Replace(string(lucking_info[0][1]), " ", "", -1))
				lucking_info[0][1] = []byte(strings.Replace(string(lucking_info[0][1]), "\n", "", -1))
				json.Unmarshal(lucking_info[0][1],&data)
			} else {
				logger.LoggerToFile("解析抽奖页面错误:"+is_luckling_url)
				return
			}
			/*判断是否在抽奖详情*/

			if data.Total != 0 {
				luck_time := Timetoymd(data.List[0].Time)
				timeNow := time.Now().Unix()
				timess, _ := time.Parse("2006-01-02 15:04:05", luck_time)
				timeUnix := timess.Unix()
				if timeNow <= timeUnix {
					logger.LoggerToFile("抽奖已结束:"+is_luckling_url)
					return
				}

				var detail model.SinaDetail
				var follow model.SinaFollow
				str_uid := string(info_submatch[0][1])
				to_uid,_ := strconv.Atoi(str_uid)

				detail.SinaDetailAdd(model.SinaDetail{
					HtmlStr: string(str1submatch[i][1]),
					Mid:string(info_submatch[0][2]),
					ToUid: to_uid,
					Uid: config["UID"].(string),
					LuckTime:luck_time,
				})
				follow.SinaFollowAdd(model.SinaFollow{
					Url: "https://weibo.com/u/"+string(info_submatch[0][1]),
					Uid: config["UID"].(string),
					Mid: string(info_submatch[0][2]),
					ToUid: to_uid,
					Nick: data.Weibo.User.Name,
				})
				at_complie := regexp.MustCompile(`<a target="_blank" render="ext" extra-data="type=atname" href="([^"]*)" usercard="name=[^"]*">@([^<]*)<\/a>`)
				at_submatch := at_complie.FindAllSubmatch(str1submatch[i][1],-1)

				//strr := fmt.Sprintf("%s",at_submatch)
				//ioutil.WriteFile("error1.txt",[]byte(strr),0777)

				for j:=0;j<len(at_submatch);j++{
					to_nick_url := `https:`+string(at_submatch[j][1])
					to_uid_url := GetLocation(to_nick_url,config["SUB"].(string))

					lucking_to_uid_complie := regexp.MustCompile(`https:\/\/weibo\.com\/?u?\/(\d+)\?from=feed`)
					lucking_to_uid_submatch := lucking_to_uid_complie.FindAllSubmatch([]byte(to_uid_url),-1)
					if len(lucking_to_uid_submatch) == 0 {
						/*获取uid*/
						req, _ := http.NewRequest("GET", to_uid_url, nil)
						req.Header.Set("cookie", "SUB="+config["SUB"].(string))
						c := http.Client{
							CheckRedirect: func(req *http.Request, via []*http.Request) error {
								return http.ErrUseLastResponse
							},
							Timeout: 30 * time.Second,
						}
						res, err := c.Do(req)
						if err != nil {
							logger.LoggerToFile("错误："+err.Error())
							return
						}
						resp,err := ioutil.ReadAll(res.Body)
						defer res.Body.Close()

						lucking_to_uid_complie = regexp.MustCompile(`\$CONFIG\['oid'\]='(\d+)'`)
						lucking_to_uid_submatch = lucking_to_uid_complie.FindAllSubmatch(resp,-1)

					}

					if len(lucking_to_uid_submatch) != 0 {
						str_uid := string(lucking_to_uid_submatch[0][1])
						to_uid,_ := strconv.Atoi(str_uid)
						follow.SinaFollowAdd(model.SinaFollow{
							Url: to_uid_url,
							Uid: config["UID"].(string),
							Mid: string(info_submatch[0][2]),
							ToUid: to_uid,
							Nick: string(at_submatch[j][2]),
						})

					} else {
						str := "总执行:"+strconv.Itoa(len(str1submatch)) +",当前执行到："+strconv.Itoa(i)  +",长度："+strconv.Itoa(len(lucking_to_uid_submatch))+",昵称地址："+to_nick_url+",uid地址:"+to_uid_url
						logger.LoggerToFile(str)
					}

				}
			}
		}
	})

	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("cookie", "SUB="+config["SUB"].(string))
		request.Headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
		fmt.Println("Visiting", request.URL)
	})
	c.OnError(func(response *colly.Response, err error) {
		logger.LoggerToFile(err.Error())
	})
	c.Visit(config["LUCKING"].(string))
}

func GetSinaLuckingApi() {
	var config = config.GetSinaUrl()
	c := colly.NewCollector()
	c.OnResponse(func(response *colly.Response) {
		var sina_luck_data respdata.SinaLuckDataApi
		json.Unmarshal(response.Body,&sina_luck_data)
		if len(sina_luck_data.Data.Cards) != 5{
			return
		}
		data := sina_luck_data.Data.Cards[4].CardGroup
		for i:=0;i<len(data);i++{
			mid := data[i].Mblog.Mid
			html_str := data[i].Mblog.Text
			to_uid := data[i].Mblog.User.Id
			nick := data[i].Mblog.User.ScreenName


			/*判断是否在抽奖详情*/
			is_luckling_url := config["LUCKING_STATUS"].(string)+mid

			req, _ := http.NewRequest("GET", is_luckling_url, nil)
			req.Header.Set("cookie", "SUB="+config["SUB"].(string))
			c := http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
				Timeout: 30 * time.Second,
			}
			r, err := c.Do(req)

			if err != nil {
				logger.LoggerToFile("错误："+err.Error())
				return
			}
			resp,err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			lucking_complie := regexp.MustCompile(`window\.__DATA__\s=\s({(?s:.*?)})\s\|\|`)
			lucking_info := lucking_complie.FindAllSubmatch(resp,-1)
			info_data := struct {
				Total int `json:"total"`
				List []struct{
					Time string `json:"time"`
				}`json:"list"`
				Weibo struct{
					User struct{
						Name string `json:"name"`
					} `json:"user"`
				}	`json:"weibo"`
			}{}

			if len(lucking_info) == 0 {
				lucking_info[0][1] = []byte(strings.Replace(string(lucking_info[0][1]), " ", "", -1))
				lucking_info[0][1] = []byte(strings.Replace(string(lucking_info[0][1]), "\n", "", -1))
				json.Unmarshal(lucking_info[0][1],&info_data)
			} else {
				logger.LoggerToFile("解析抽奖页面错误:"+is_luckling_url)
				return
			}

			if info_data.Total != 0 {

				luck_time := Timetoymd(info_data.List[0].Time)
				timeNow := time.Now().Unix()
				timess, _ := time.Parse("2006-01-02 15:04:05", luck_time)
				timeUnix := timess.Unix()
				if timeNow <= timeUnix {
					logger.LoggerToFile("抽奖已结束:"+is_luckling_url)
					return
				}


				var detail model.SinaDetail
				var follow model.SinaFollow
				detail.SinaDetailAdd(model.SinaDetail{
					HtmlStr: html_str,
					Mid:mid,
					ToUid: to_uid,
					Uid: config["UID"].(string),
					LuckTime: luck_time,
				})
				follow.SinaFollowAdd(model.SinaFollow{
					Url: "https://weibo.com/u/"+strconv.Itoa(to_uid),
					Uid: config["UID"].(string),
					Mid: mid,
					ToUid:to_uid,
					Nick: nick,
				})


				at_complie := regexp.MustCompile(`<a href='([^"]*)'>@([^<]*)<\/a>`)
				at_submatch := at_complie.FindAllSubmatch([]byte(html_str),-1)

				for j:=0;j<len(at_submatch);j++{
					to_nick_url := `https://weibo.com`+string(at_submatch[j][1])
					to_uid_url := GetLocation(to_nick_url,config["SUB"].(string))

					lucking_to_uid_complie := regexp.MustCompile(`https:\/\/weibo\.com\/?u?\/(\d+)\?from=feed`)
					lucking_to_uid_submatch := lucking_to_uid_complie.FindAllSubmatch([]byte(to_uid_url),-1)



					//fmt.Println("总执行:"+strconv.Itoa(len(at_submatch)))
					//fmt.Println("当前执行到："+strconv.Itoa(j))
					//fmt.Println("长度："+strconv.Itoa(len(lucking_to_uid_submatch)))
					//str := fmt.Sprintf("%s",to_nick_url)
					//fmt.Println(str)

					if len(lucking_to_uid_submatch) == 0 {
						/*获取uid*/
						req, _ := http.NewRequest("GET", to_uid_url, nil)
						req.Header.Set("cookie", "SUB="+config["SUB"].(string))
						c := http.Client{
							CheckRedirect: func(req *http.Request, via []*http.Request) error {
								return http.ErrUseLastResponse
							},
							Timeout: 30 * time.Second,
						}
						res, err := c.Do(req)
						if err != nil {
							logger.LoggerToFile("错误："+err.Error())
							return
						}
						resp,err := ioutil.ReadAll(res.Body)
						defer res.Body.Close()
						lucking_to_uid_complie = regexp.MustCompile(`\$CONFIG\['oid'\]='(\d+)'`)
						lucking_to_uid_submatch = lucking_to_uid_complie.FindAllSubmatch(resp,-1)
					}

					if len(lucking_to_uid_submatch) != 0 {
						toto_uid := string(lucking_to_uid_submatch[0][1])
						ttoo_uid,_ := strconv.Atoi(toto_uid)
						follow.SinaFollowAdd(model.SinaFollow{
							Url: to_uid_url,
							Uid: config["UID"].(string),
							Mid: mid,
							ToUid: ttoo_uid,
							Nick: string(at_submatch[j][2]),
						})
					} else {
						str := "总执行:"+strconv.Itoa(len(data)) +",当前执行到："+strconv.Itoa(i)  +",长度："+strconv.Itoa(len(lucking_to_uid_submatch))+",昵称地址："+to_nick_url+",uid地址:"+to_uid_url
						logger.LoggerToFile(str)
					}
				}
			}
		}
	})

	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("cookie", "SUB="+config["SUB"].(string))
		request.Headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
		fmt.Println("Visiting", request.URL)
	})
	c.OnError(func(response *colly.Response, err error) {
		logger.LoggerToFile(err.Error())
	})
	c.Visit(config["LUCKING_TIME"].(string))
}

func GetLuckSearchApi() {
	var config = config.GetSinaUrl()
	c := colly.NewCollector()
	c.OnResponse(func(response *colly.Response) {
		var sina_luck_search respdata.SinaLuckSearchApi
		json.Unmarshal(response.Body,&sina_luck_search)
		if len(sina_luck_search.Data.Cards) == 0 {
			return
		}
		data := sina_luck_search.Data.Cards
		//str:= fmt.Sprintf("%s",data[0].Mblog)
		//ioutil.WriteFile("a.txt",[]byte(str),0777)

		for i:=0;i<len(data);i++{
			if data[i].CardType == 9 {
				mid := data[i].Mblog.Mid
				html_str := data[i].Mblog.Text
				to_uid := data[i].Mblog.User.Id
				nick := data[i].Mblog.User.ScreenName

				/*判断是否在抽奖详情*/
				is_luckling_url := config["LUCKING_STATUS"].(string)+mid
				req, _ := http.NewRequest("GET", is_luckling_url, nil)
				req.Header.Set("cookie", "SUB="+config["SUB"].(string))
				c := http.Client{
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
					Timeout: 30 * time.Second,
				}
				r, err := c.Do(req)
				if err != nil {
					logger.LoggerToFile("错误："+err.Error())
					return
				}
				resp,err := ioutil.ReadAll(r.Body)
				defer r.Body.Close()
				lucking_complie := regexp.MustCompile(`window\.__DATA__\s+=\s+({(?s:.*?)})\s+\|\|`)
				lucking_info := lucking_complie.FindAllSubmatch(resp,-1)
				info_data := struct {
					Total int `json:"total"`
					List []struct{
						Time string `json:"time"`
					}`json:"list"`
					Weibo struct{
						User struct{
							Name string `json:"name"`
						} `json:"user"`
					}	`json:"weibo"`
				}{}

				if len(lucking_info) != 0 {
					lucking_info[0][1] = []byte(strings.Replace(string(lucking_info[0][1]), " ", "", -1))
					lucking_info[0][1] = []byte(strings.Replace(string(lucking_info[0][1]), "\n", "", -1))
					json.Unmarshal(lucking_info[0][1],&info_data)
				} else {
					logger.LoggerToFile("解析抽奖页面错误:"+is_luckling_url)
					return
				}

				if info_data.Total != 0 {

					luck_time := Timetoymd(info_data.List[0].Time)
					timeNow := time.Now().Unix()
					timess, _ := time.Parse("2006-01-02 15:04:05", luck_time)
					timeUnix := timess.Unix()
					if timeNow >= timeUnix {
						logger.LoggerToFile("抽奖已结束:"+is_luckling_url)
						return
					}

					var detail model.SinaDetail
					var follow model.SinaFollow
					detail.SinaDetailAdd(model.SinaDetail{
						HtmlStr: html_str,
						Mid:mid,
						ToUid: to_uid,
						Uid: config["UID"].(string),
						LuckTime: luck_time,
					})
					follow.SinaFollowAdd(model.SinaFollow{
						Url: "https://weibo.com/u/"+strconv.Itoa(to_uid),
						Uid: config["UID"].(string),
						Mid: mid,
						ToUid:to_uid,
						Nick: nick,
					})
					at_complie := regexp.MustCompile(`<a href='([^"']*)'>@([^<]*)<\/a>`)
					at_submatch := at_complie.FindAllSubmatch([]byte(html_str),-1)



					for j:=0;j<len(at_submatch);j++{
						to_nick_url := `https://weibo.com`+string(at_submatch[j][1])
						to_uid_url := GetLocation(to_nick_url,config["SUB"].(string))
						lucking_to_uid_complie := regexp.MustCompile(`https:\/\/weibo\.com\/?u?\/(\d+)\?from=feed`)
						lucking_to_uid_submatch := lucking_to_uid_complie.FindAllSubmatch([]byte(to_uid_url),-1)


						//fmt.Println("总执行:"+strconv.Itoa(len(data)))
						//fmt.Println("当前执行到："+strconv.Itoa(i))
						//fmt.Println("长度："+strconv.Itoa(len(lucking_to_uid_submatch)))
						//str := fmt.Sprintf("%s",to_nick_url)
						//fmt.Println(str)

						if len(lucking_to_uid_submatch) == 0 {
							/*获取uid*/
							req, _ := http.NewRequest("GET", to_uid_url, nil)
							req.Header.Set("cookie", "SUB="+config["SUB"].(string))
							c := http.Client{
								CheckRedirect: func(req *http.Request, via []*http.Request) error {
									return http.ErrUseLastResponse
								},
								Timeout: 30 * time.Second,
							}
							res, err := c.Do(req)
							if err != nil {
								logger.LoggerToFile("错误："+err.Error())
								return
							}
							resp,err := ioutil.ReadAll(res.Body)
							defer res.Body.Close()
							lucking_to_uid_complie = regexp.MustCompile(`\$CONFIG\['oid'\]='(\d+)'`)
							lucking_to_uid_submatch = lucking_to_uid_complie.FindAllSubmatch(resp,-1)
						}

						if len(lucking_to_uid_submatch) != 0 {
							toto_uid := string(lucking_to_uid_submatch[0][1])
							ttoo_uid,_ := strconv.Atoi(toto_uid)
							follow.SinaFollowAdd(model.SinaFollow{
								Url: to_uid_url,
								Uid: config["UID"].(string),
								Mid: mid,
								ToUid: ttoo_uid,
								Nick: string(at_submatch[j][2]),
							})
						} else {
							str := "总执行:"+strconv.Itoa(len(data)) +",当前执行到："+strconv.Itoa(i)  +",长度："+strconv.Itoa(len(lucking_to_uid_submatch))+",昵称地址："+to_nick_url+",uid地址:"+to_uid_url
							logger.LoggerToFile(str)
						}
					}
				}

			}
		}
	})

	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("cookie", "SUB="+config["SUB"].(string))
		request.Headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
		fmt.Println("Visiting", request.URL)
	})
	c.OnError(func(response *colly.Response, err error) {
		logger.LoggerToFile(err.Error())
	})
	c.Visit(config["LUCKING_SEARCH"].(string))
	c.Visit(config["LUCKING_SEARCH_ZHUANFA"].(string))
	c.Visit(config["LUCKING_SEARCH_XIANGQING"].(string))
}



func GetLocation(url,sub string)string{
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("cookie", "SUB="+sub)
	c := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 30 * time.Second,
	}
	resp, err := c.Do(req)
	if err != nil {
		return ""
	}
	return resp.Header.Get("location")
}


//关注
func FollowSet() {
	var config = config.GetSinaUrl()
	c := colly.NewCollector()
	var sina_follow model.SinaFollow
	data := sina_follow.SinaFollowList("is_modify = 0 and uid = "+config["UID"].(string))

	for _,v := range data{
		c.OnResponse(func(response *colly.Response) {
			var sina_code respdata.SinaCode
			json.Unmarshal(response.Body, &sina_code)
			if sina_code.Code != "100000" {
				logger.LoggerToFile("错误：" + sina_code.Msg+"，错误码:"+sina_code.Code)
				return
			}
			var sina_follow model.SinaFollow
			sina_follow.SinaFollowSave(model.SinaFollow{
				Id: v.Id,
				IsModify: 1,
			})
			fmt.Println("关注成功")
		})
		c.OnRequest(func(request *colly.Request) {
			request.Headers.Set("cookie", "SUB="+config["SUB"].(string))
			request.Headers.Set("referer", "https://weibo.com/u/"+strconv.Itoa(v.ToUid))
			request.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
			request.Headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
		})
		c.OnError(func(response *colly.Response, err error) {
			logger.LoggerToFile(err.Error())
		})
		to_uid := strconv.Itoa(v.ToUid)
		c.Post(config["FOLLOW_URL"].(string), map[string]string{
			"uid":    to_uid,
			"location":    config["PAGE_ID"].(string),
			"oid":    to_uid,
		})
		time.Sleep(10*time.Second)

	}
}

//点赞
func LikeSet() {
	var config = config.GetSinaUrl()
	c := colly.NewCollector()
	var sina_detail model.SinaDetail
	data := sina_detail.SinaDetailList("is_like = 0 and uid = "+config["UID"].(string))


	for _,v := range data{
		c.OnResponse(func(response *colly.Response) {
			var sina_code respdata.SinaCode
			json.Unmarshal(response.Body, &sina_code)
			if sina_code.Code != "100000" {
				logger.LoggerToFile("错误：" + sina_code.Msg+"，错误码:"+sina_code.Code)
				return
			}
			var sina_detail model.SinaDetail
			sina_detail.SinaDetailSave(model.SinaDetail{
				Id: v.Id,
				IsLike: 1,
			})

			fmt.Println("点赞成功")
		})
		c.OnRequest(func(request *colly.Request) {
			request.Headers.Set("cookie", "SUB="+config["SUB"].(string))
			request.Headers.Set("referer", "https://weibo.com/u/"+strconv.Itoa(v.ToUid))
			request.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
			request.Headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
		})
		c.OnError(func(response *colly.Response, err error) {
			logger.LoggerToFile(err.Error())
		})


		c.Post(config["LIKE_URL"].(string), map[string]string{
			"mid":    v.Mid,
		})
		time.Sleep(10*time.Second)
	}


}

//转发抽奖话题 转发+评论
func HuaTiZhuanFa() {
	var config = config.GetSinaUrl()
	c := colly.NewCollector()

	var sina_detail model.SinaDetail
	data := sina_detail.SinaDetailList("is_repost = 0 and uid = "+config["UID"].(string))

	for _,v := range data{

		strArr := [5]string{"不错","来了来了","冲","真不错","[抱一抱]"}
		rand.Seed(time.Now().UnixNano())
		str := strArr[rand.Intn(len(strArr)-1)]
		var strr string
		if strings.Index(v.HtmlStr,"好友") != -1 || strings.Index(v.HtmlStr,"1好友") != -1 || strings.Index(v.HtmlStr,"1个好友") != -1 || strings.Index(v.HtmlStr,"一好友") != -1 || strings.Index(v.HtmlStr,"一个好友") != -1 || strings.Index(v.HtmlStr,"1位好友") != -1 || strings.Index(v.HtmlStr,"一位好友") != -1{
			strr = " @用户6287329627 "
		}
		if strings.Index(v.HtmlStr,"2好友") != -1 || strings.Index(v.HtmlStr,"2个好友") != -1 || strings.Index(v.HtmlStr,"二好友") != -1 || strings.Index(v.HtmlStr,"二个好友") != -1 || strings.Index(v.HtmlStr,"两好友") != -1 || strings.Index(v.HtmlStr,"两位好友") != -1 || strings.Index(v.HtmlStr,"二位好友") != -1 || strings.Index(v.HtmlStr,"两个好友") != -1{
			strr = " @用户6287329627 @用户6930417324 "
		}

		if strings.Index(v.HtmlStr,"3好友") != -1 || strings.Index(v.HtmlStr,"3个好友") != -1 || strings.Index(v.HtmlStr,"三好友") != -1 || strings.Index(v.HtmlStr,"三个好友") != -1 || strings.Index(v.HtmlStr,"三好友") != -1 || strings.Index(v.HtmlStr,"三位好友") != -1 || strings.Index(v.HtmlStr,"三位好友") != -1 || strings.Index(v.HtmlStr,"三个好友") != -1{
			strr = " @用户6287329627 @用户6930417324 @用户6975713923 "
		}
		str += strr

		c.OnResponse(func(response *colly.Response) {
			var sina_code respdata.SinaCode
			json.Unmarshal(response.Body, &sina_code)
			if sina_code.Code == "100001" {
				//只转发
				var r http.Request
				r.ParseForm()
				r.Form.Add("mid",v.Mid)
				r.Form.Add("style_type","1")
				r.Form.Add("reason",str)
				r.Form.Add("location",config["PAGES_ID"].(string))
				r.Form.Add("pdetail","100606"+strconv.Itoa(v.ToUid))
				body_str := strings.TrimSpace(r.Form.Encode())
				reqs,err := http.NewRequest("POST",config["ZHUANFA_URL"].(string),strings.NewReader(body_str))
				if err != nil {
					logger.LoggerToFile("错误："+err.Error())
					return
				}
				reqs.Header.Set("cookie", "SUB="+config["SUB"].(string))
				reqs.Header.Set("referer", "https://weibo.com/u/"+strconv.Itoa(v.ToUid))
				reqs.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				reqs.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
				var resp *http.Response
				resp, err = http.DefaultClient.Do(reqs)
				if err != nil {
					logger.LoggerToFile("错误："+err.Error())
					return
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				json.Unmarshal(body, &sina_code)
				if sina_code.Code != "100000" {
					logger.LoggerToFile("错误：" + sina_code.Msg+"，错误码："+sina_code.Code)
					return
				}
				var sina_detail model.SinaDetail
				sina_detail.SinaDetailSave(model.SinaDetail{
					Id: v.Id,
					IsRepost: 1,
				})
				fmt.Println("转发成功")
				//只转发
			}else {
				if sina_code.Code != "100000" {
					logger.LoggerToFile("错误：" + sina_code.Msg+"，错误码："+sina_code.Code)
					return
				}
				var sina_detail model.SinaDetail
				sina_detail.SinaDetailSave(model.SinaDetail{
					Id: v.Id,
					IsRepost: 1,
				})
				fmt.Println("转发评论成功")
			}

		})
		c.OnRequest(func(request *colly.Request) {
			request.Headers.Set("cookie", "SUB="+config["SUB"].(string))
			request.Headers.Set("referer", "https://weibo.com/u/"+strconv.Itoa(v.ToUid))
			request.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
			request.Headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
		})
		c.OnError(func(response *colly.Response, err error) {
			logger.LoggerToFile(err.Error())
		})

		c.Post(config["COMMENT_URL"].(string), map[string]string{
			"mid":    v.Mid,
			"uid":   v.Uid,
			"forward": "1",
			"content":   str,
			"location":  config["PAGE_ID"].(string),
			"pdetail":  "100505"+strconv.Itoa(v.ToUid),
		})
		time.Sleep(10*time.Second)

	}
}

func Timetoymd(str string) string {
	ymd := strings.Replace(str,"年","-",-1)
	ymd = strings.Replace(ymd,"月","-",-1)
	ymd = strings.Replace(ymd,"日"," ",-1)
	return ymd+":00"
}