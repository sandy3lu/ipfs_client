package storageprove

import (
	"time"
	"fmt"
	"net/http"
	"io/ioutil"
	"context"
	"encoding/json"
	"strings"
	"net/url"
)

const URL string = "http://wallet-api-test.launchain.org:50000"// update online status to server
//const URL string = "http://10.0.0.116:50000"// update online status to server
var SNList  []string

func Startheart(node string, ctx context.Context){

	t1 := time.NewTimer(time.Minute * 60)

	registerNode(node )
	var err1 error
	SNList, err1= GetSuperNodeList()
	if err1!=nil {
		SNList, err1= GetSuperNodeList()// try again
	}
	go func() {
		for{
			select{
			case <-t1.C:
				if len(SNList)==0 {
					SNList, err1= GetSuperNodeList()// try again
				}
				registerNode(node )
				t1.Reset(time.Minute * 60)
			case <-ctx.Done():
				return
			}
		}
	}()

}


func registerNode(node string){
	url := URL + "/v1/ipfs/sub-nodes/" + node
	req, err := http.NewRequest("PUT", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode!=200 {
		fmt.Print("registerNode ",resp.Status)
		fmt.Println(string(body))

	}

	if err != nil {
		fmt.Println("registerNode ",err)

	}
}

type superNodeList struct{
	Count int `json:"count"`
	Info  []superNodeInfo `json:"info"`
}

type superNodeInfo struct {

	Id string `json:"_id"`//"_id": "5b20ccc198504f6a3ab83f71",
	Address string `json:"address"`//"address": "112345",
	Created_at string `json:"created_at"`//"created_at": "2018-06-13T16:05:28.709+08:00",
	Status int `json:"status"`//"status": 1,
	Updated_at string `json:"updated_at"`//"updated_at": "2018-06-13T16:05:28.709+08:00"

}

func GetSuperNodeList()  ([]string, error){
	url := URL + "/v1/ipfs/super-nodes?page=0&limit=100"
	req, err := http.NewRequest("GET", url, nil)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode!=200 {
		fmt.Print(resp.Status)
		fmt.Println(string(body))
		return nil,fmt.Errorf("status code %d",resp.StatusCode )
	}

	if err != nil {
		fmt.Println(err)
		return nil,err
	}

	//dat:=make(map[string]interface{})
	var overDue superNodeList
	var fileName = make([]string,0)
	//fmt.Println(string(body))
	if err := json.Unmarshal(body, &overDue); err == nil {
		num:=overDue.Count
		fmt.Println("superNode List......")
		for  i:=0;i<num;i++ {
			ipfsfile := overDue.Info[i]

			if strings.ContainsAny(ipfsfile.Address,"Qm") {

				//fmt.Println(ipfsfile.Address)
				fileName = append(fileName, ipfsfile.Address)
			}
		}
		return fileName, nil
	} else {

		return nil,fmt.Errorf("json str to struct error")
	}
	return nil,err
}



func ReportNodeStatus(node, stat string){
	var clusterinfo = url.Values{}
	clusterinfo.Add("node_address", node)
	clusterinfo.Add("info", stat)
	url_addr := URL + "/v1/ipfs/flow"
	data := clusterinfo.Encode()

	req, err := http.NewRequest("POST", url_addr, strings.NewReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")


	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode!=200 {
		fmt.Print("registerNode ",resp.Status)
		fmt.Println(string(body))

	}

	if err != nil {
		fmt.Println("registerNode ",err)

	}
}
