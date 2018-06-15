package storageprove

import (
	"time"
	"fmt"
	"net/http"
	"io/ioutil"
	"context"
)

const URL string = "http://wallet-api-test.launchain.org:50000"// update online status to server
func Startheart(node string, ctx context.Context){

	t1 := time.NewTimer(time.Minute * 60)
	fmt.Println(node, time.Now())
	registerNode(node )

	go func() {
		for{
			select{
			case <-t1.C:
				fmt.Println(node, time.Now())
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
		fmt.Print(resp.Status)
		fmt.Println(string(body))

	}

	if err != nil {
		fmt.Println(err)

	}
}
