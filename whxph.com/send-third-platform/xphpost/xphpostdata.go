package xphpostdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"unsafe"

	"github.com/sirupsen/logrus"
	"whxph.com/send-third-platform/xphapi"
)

var (
	usersToken   string
	usersDevices []xphapi.Device
)

func Start() {
	logrus.Info("XphPostData  Start-----")
	getUsersToken()
	usersUpdateDevices()
	for {

		xphpostdata()
	}
	//c := cron.New()
	//_ = c.AddFunc("0 0/10 * * * *", xphpostdata)
	//c.Start()
	//defer c.Stop()
	//select {}
}

func getUsersToken() {
	usersToken = xphapi.NewGetToken("test", "123456")

}
func usersUpdateDevices() {
	usersDevices = xphapi.NewGetDevices("test", usersToken)
}

func xphpostdata() {

	for _, item := range usersDevices {
		rclient := &http.Client{Timeout: 10 * time.Second}
		resp, err := http.NewRequest("GET", "http://47.105.215.208:8005/data/"+strconv.Itoa(item.DeviceID), nil)
		if err != nil {
			logrus.Error("获取数据异常")
			return
		}
		//resp.Header.Set("Content-Type", "application/json")
		resp.Header.Set("token", usersToken)
		respa, err := rclient.Do(resp)
		if err != nil {
			return
		}

		//fmt.Println(respa.Request.URL)
		//fmt.Println(usersToken)
		defer respa.Body.Close()

		buffer, err := ioutil.ReadAll(respa.Body)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		logrus.Info("ReadAll:", string(buffer))

		body, err := json.Marshal(buffer)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		reader := bytes.NewReader(body)
		request, err := http.NewRequest("POST", "http://47.105.215.208:8005//data/"+strconv.Itoa(item.DeviceID), reader)
		if err != nil {
			fmt.Printf("http.NewRequest%v", err)
			return
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("token", usersToken)
		client := http.Client{Timeout: 10 * time.Second}
		_, err = client.Do(request)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		str := (*string)(unsafe.Pointer(&respBytes))
		fmt.Println(*str)
		time.Sleep(30 * time.Second)
	}
}
