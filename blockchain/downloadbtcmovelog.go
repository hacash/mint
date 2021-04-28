package blockchain

import (
	"fmt"
	"github.com/hacash/core/stores"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 下载比特币单向转移记录
func (bc *BlockChain) downLoadBTCMoveLog() {

	sleepTimeErr := time.Minute * 45
	sleepTimeEmpty := time.Hour * 8

	// test start
	//sleepTimeErr = time.Second * 3
	//sleepTimeEmpty = sleepTimeErr
	// test end

	store := bc.chainstate.BlockStore()
	realpage, e := store.GetBTCMoveLogTotalPage()
	if e != nil {
		fmt.Println(e)
		return
	}
	downloadUrl := bc.config.DownloadBTCMoveLogUrl
	reqpage := realpage
	// 分页读取
	limit := stores.SatoshiGenesisLogStorePageLimit
	lastpageData := []string{}
	if reqpage == 0 {
		reqpage = 1 // 首次获取第一页
		lastpageData = []string{}
	} else {
		// 读取最后一页数据
		list, e := store.GetBTCMoveLogPageData(reqpage)
		if e != nil {
			fmt.Println(e)
			return
		}
		lastpageData = list
	}
	// 等待读取新增
	lastpage := reqpage
	// 循环读取
	for {
		lastdatasize := len(lastpageData)
		addgetstart := lastdatasize + ((lastpage - 1) * limit) + 1
		addgetlimit := limit - lastdatasize
		// 读取
		pagedata, err := readSatoshiGenesisByUrl(downloadUrl, addgetstart, addgetlimit)
		if err != nil || len(pagedata) == 0 {
			time.Sleep(sleepTimeErr)
			continue // 下载错误，休眠后重试
		}
		if len(pagedata) == 0 {
			// 下载为空
			time.Sleep(sleepTimeEmpty)
			continue // 数据为空，休眠 8 小时 后重试
		}
		lastpageData = append(lastpageData, pagedata...)
		// 保存
		e1 := store.SaveBTCMoveLogPageData(lastpage, lastpageData)
		if e1 != nil {
			fmt.Println("[Satoshi genesis] SaveBTCMoveLogPageData Error:", e1.Error())
			return
		}
		if len(lastpageData) == limit {
			// 满页了
			lastpageData = []string{}
			lastpage++
			continue // 立即下一页
		}
		// 休眠8小时后继续
		time.Sleep(sleepTimeEmpty)
		continue
	}
}

func readSatoshiGenesisByUrl(url string, start int, limit int) ([]string, error) {
	if len(url) == 0 {
		return nil, fmt.Errorf("url is empty") // error
	}
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	url += fmt.Sprintf("?trsno=%d&limit=%d", start, limit)
	fmt.Printf("[Satoshi genesis] load check url: %s\n", url)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("[Satoshi genesis] read Validated SatoshiGenesisByUrl return error: %s.\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	resstr := string(body)
	if len(resstr) < 32 {
		// 返回空
		fmt.Printf("[Satoshi genesis] got data count 0.\n")
		return []string{}, nil
	}
	// 有内容
	logs := strings.Split(resstr, "|")
	fmt.Printf("[Satoshi genesis] got data count %d.\n", len(logs))
	// 解析
	return logs, nil
}
