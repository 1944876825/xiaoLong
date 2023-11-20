package could

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"xiaoLong/xc"
)

var XiaoLong = &xiaoLong{
	ForMatText: "{FileName}\t{FileCid}\t{FileDir}",
}

type xiaoLong struct {
	AccessToken string
	ForMatText  string
	SaveType    int
	Count       int
}

type fileSuperListParams struct {
	FileType   []interface{} `json:"fileType"`
	Keywords   string        `json:"keywords"`
	PageNum    int           `json:"pageNum"`
	PageSize   int           `json:"pageSize"`
	ParentId   string        `json:"parentId"`
	SortMethod string        `json:"sortMethod"`
	SortType   string        `json:"sortType"`
}
type fileSuperListData struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Submessage string      `json:"submessage"`
	Data       []fileModel `json:"data"`
	Count      int         `json:"count"`
	Stime      int         `json:"stime"`
}
type fileModel struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	ParentId    string `json:"parentId"`
	FileName    string `json:"fileName"`
	FileCid     string `json:"fileCid"`
	FileSize    int    `json:"fileSize"`
	FileType    int    `json:"fileType"`
	IsFolder    int    `json:"isFolder"`
	IsDel       int    `json:"isDel"`
	Thumbnail   string `json:"thumbnail"`
	Width       string `json:"width"`
	Height      string `json:"height"`
	Duration    int    `json:"duration"`
	MigrateHash string `json:"migrateHash"`
	Collect     int    `json:"collect"`
	Suffix      string `json:"suffix"`
	IsTrash     int    `json:"isTrash"`
	Cover       string `json:"cover"`
	Storage     string `json:"storage"`
	Ptime       int    `json:"ptime"`
	IsCheck     int    `json:"is_check"`
	FileDir     string `json:"fileDir"`
}
type loginData struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Submessage string `json:"submessage"`
	Data       struct {
		Token        string `json:"token"`
		Id           string `json:"id"`
		PeerId       string `json:"peerId"`
		Name         string `json:"name"`
		Phone        string `json:"phone"`
		Sex          int    `json:"sex"`
		Ptime        int    `json:"ptime"`
		Utime        int    `json:"utime"`
		Nickname     string `json:"nickname"`
		Img          string `json:"img"`
		LikeNum      int    `json:"likeNum"`
		AttentionNum int    `json:"attentionNum"`
		FansNum      int    `json:"fansNum"`
		ArticleNum   int    `json:"articleNum"`
		IsAttention  bool   `json:"isAttention"`
		Showhl       bool   `json:"showhl"`
		Profile      string `json:"profile"`
		HadPayPasswd bool   `json:"hadPayPasswd"`
		IdCard       string `json:"idCard"`
		Role         string `json:"role"`
		InWhiteList  int    `json:"inWhiteList"`
		IsVip        int    `json:"isVip"`
		VipDeadline  int    `json:"vipDeadline"`
	} `json:"data"`
	Count int `json:"count"`
	Stime int `json:"stime"`
}

func (x *xiaoLong) GetFileSuperList(file *os.File, parentId string, parentName string, pageNum int) {
	var list []string
	data := `{"fileType":[],"keywords":"","pageNum":%d,"pageSize":100,"parentId":"%s","sortMethod":"desc","sortType":"fileName"}`
	data = fmt.Sprintf(data, pageNum, parentId)
	api := "https://productapi.stariverpan.com/cmsprovider/v1.2/cloud/fileSuperList"
	req, err := http.NewRequest("POST", api, strings.NewReader(data))
	if err != nil {
		fmt.Println("GetFileSuperList req err", err)
		xc.Wg.Done()
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", x.AccessToken)
	req.Header.Set("Referer", "https://share.stariverpan.com/")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("GetFileSuperList resp err", err)
		xc.Wg.Done()
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("GetFileSuperList io read err", err)
		xc.Wg.Done()
		return
	}
	if resp.StatusCode == http.StatusOK {
		var s fileSuperListData
		err := json.Unmarshal(body, &s)
		if err == nil {
			if s.Code == 10010 {
				fmt.Println("GetFileSuperList 登录失败", string(body))
				xc.Wg.Done()
				return
			}
			if s.Code == 200 {
				for _, d := range s.Data {
					// 目录处理
					if d.IsFolder == 1 {
						if x.SaveType == 1 {
							xc.Wg.Add(1)
							go x.GetFileSuperList(file, d.Id, parentName+"/"+d.FileName, 0)
						}
						if x.SaveType == 2 {
							var fn string
							if d.FileName != "" {
								fn = strings.Replace(d.FileName, "\t", "-", -1)
							} else {
								fn = "没有名称-" + d.FileCid
							}
							fmt.Println(fn)
							nFile, err := os.Create(filepath.FromSlash("cid/" + fn + ".txt"))
							xc.FileXc = append(xc.FileXc, nFile)
							if err == nil {
								xc.Wg.Add(1)
								go x.GetFileSuperList(nFile, d.Id, parentName+"/"+d.FileName, 0)
							} else {
								fmt.Println("文件创建失败", fn, err)
							}
						}
					} else {
						fmt.Println(d.FileName)
						d.FileDir = parentName
						list = append(list, x.myFormat(d, x.ForMatText)+"\n")
						xc.Mutex.Lock()
						x.Count += 1
						xc.Mutex.Unlock()
					}
				}
			}
		} else {
			fmt.Println("GetFileSuperList json un err", err)
		}
		if len(list) > 0 {
			for _, l := range list {
				xc.Mutex.Lock()
				_, _ = file.WriteString(l)
				xc.Mutex.Unlock()
			}
		}
		if len(s.Data) >= 100 {
			pageNum++
			xc.Wg.Add(1)
			go x.GetFileSuperList(file, parentId, parentName, pageNum)
		}
	}
	xc.Wg.Done()
}

func (x *xiaoLong) Login(token string) bool {
	api := "https://productapi.stariverpan.com/cmsprovider/v2.5/user/login"
	data := fmt.Sprintf(`{"token":"%s"}`, token)
	req, err := http.NewRequest("POST", api, strings.NewReader(data))
	if err != nil {
		fmt.Println("login req err", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://share.stariverpan.com/")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("login resp err", err)
		return false
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("login io read err", err)
		return false
	}
	if resp.StatusCode == http.StatusOK {
		var s loginData
		err := json.Unmarshal(body, &s)
		if err == nil {
			if s.Code == 200 {
				fmt.Println("欢迎您！", s.Data.Name)
				x.AccessToken = "Bearer " + token
				return true
			}
		} else {
			fmt.Println("login json un err", err)
			return false
		}
	}
	fmt.Println("login err", string(body))
	return false
}

type FormatItem struct {
	Key   string
	Title string
}

func (*xiaoLong) myFormat(json fileModel, text string) string {
	if text == "" {
		return ""
	}
	var arr = []FormatItem{
		{
			Key:   "FileName",
			Title: "文件名称",
		},
		{
			Key:   "FileCid",
			Title: "文件cid",
		},
		{
			Key:   "FileDir",
			Title: "文件路径",
		},
		{
			Key:   "FileSize",
			Title: "文件大小",
		},
		{
			Key:   "Cover",
			Title: "文件封面cid",
		},
		{
			Key:   "Suffix",
			Title: "扩展名",
		},
	}
	v := reflect.ValueOf(json)
	text = strings.Replace(text, "\\t", "\t", -1)
	for _, item := range arr {
		if strings.Contains(text, "{"+item.Key+"}") {
			field := v.FieldByName(item.Key)
			if field.IsValid() {
				value := field.Interface()
				if str, ok := value.(string); ok {
					text = strings.Replace(text, "{"+item.Key+"}", str, -1)
				} else if num, ok := value.(int); ok {
					text = strings.Replace(text, "{"+item.Key+"}", strconv.Itoa(num), -1)
				}
			}
		}
	}
	return text
}
