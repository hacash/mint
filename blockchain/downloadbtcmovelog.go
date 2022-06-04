package blockchain

import (
	"fmt"
	"github.com/hacash/core/stores"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Download bitcoin one-way transfer records
func (bc *BlockChain) downLoadBTCMoveLog() {

	sleepTime := time.Hour * 8

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
	// Paging read
	limit := stores.SatoshiGenesisLogStorePageLimit
	lastpageData := []string{}
	if reqpage == 0 {
		reqpage = 1 // Get the first page for the first time
		lastpageData = []string{}
	} else {
		// Read last page data
		list, e := store.GetBTCMoveLogPageData(reqpage)
		if e != nil {
			fmt.Println(e)
			return
		}
		lastpageData = stores.SatoshiGenesisPageSerializeForShow(list)
	}
	// Waiting to read new
	lastpage := reqpage
	// Cyclic read
	for {
		lastdatasize := len(lastpageData)
		addgetstart := lastdatasize + ((lastpage - 1) * limit) + 1
		addgetlimit := limit - lastdatasize
		// read
		pagedata, err := readSatoshiGenesisByUrl(downloadUrl, addgetstart, addgetlimit)
		if err != nil || len(pagedata) == 0 {
			// Download is empty
			time.Sleep(sleepTime)
			continue // The data is empty. Try again after 8 hours of hibernation
		}
		lastpageData = append(lastpageData, pagedata...)
		// preservation
		e1 := store.SaveBTCMoveLogPageData(lastpage, stores.SatoshiGenesisPageParseForShow(lastpageData))
		if e1 != nil {
			fmt.Println("[Satoshi genesis] SaveBTCMoveLogPageData Error:", e1.Error())
			return
		}
		if len(lastpageData) == limit {
			// Full page
			lastpageData = []string{}
			lastpage++
			continue // Next page now
		}
		// Continue after 8 hours of sleep
		time.Sleep(sleepTime)
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
		// Return null
		fmt.Printf("[Satoshi genesis] got data count 0.\n")
		return []string{}, nil
	}
	// With content
	logs := strings.Split(resstr, "|")
	fmt.Printf("[Satoshi genesis] got data count %d.\n", len(logs))
	// analysis
	return logs, nil
}
