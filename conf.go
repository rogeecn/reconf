package reconf

import (
	"github.com/levigross/grequests"
	"strconv"
	"fmt"
	"os"
	"time"
)

var c map[string]string
var appName string

const conf_url = "http://token.qoofan.com/api/app-config"

func Init(app_name string) error {
	fmt.Println("begin: ", time.Now().Unix())
	appName = app_name

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		for {
			<-time.After(time.Second * 10)

			fmt.Println("updating config...")
			err := Update()
			if err != nil {
				panic("config fetch error")
			}
			fmt.Println("update config complete...")
		}
	}()
	return Update()
}

func Update() error {
	hostname, _ := os.Hostname()
	params := map[string]string{
		"app": appName,
		"hostname": hostname,
	}

	resp, err := grequests.Get(conf_url, &grequests.RequestOptions{
		Params:params,
	})
	if err != nil {
		return err
	}
	fmt.Println(">>>>> remote string configure data: ", resp.String())

	var jsonResp struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data map[string]string `json:"data"`
	}
	err = resp.JSON(&jsonResp)
	if err != nil {
		return err
	}

	if jsonResp.Code != 0 {
		return fmt.Errorf("Remote data error, Err: %s", jsonResp.Msg)
	}
	c = jsonResp.Data

	fmt.Println("=========================================================")
	for key, value := range c {
		fmt.Printf("%30s : %s\n", key, value)
	}
	fmt.Println("=========================================================")

	return nil
}

func getValue(key string) (string, error) {
	value, ok := c[key]
	if !ok {
		return "", fmt.Errorf("Key %s not exist", key)
	}
	return value, nil
}

func String(key string) (string, error) {
	mapValue, err := getValue(key)
	if err != nil {
		return "", err
	}
	return mapValue, nil
}
func Bool(key string) (bool, error) {
	mapValue, err := getValue(key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(mapValue)
}

func Float64(key string) (float64, error) {
	value, err := getValue(key)
	if err != nil {
		return 0.0, err
	}
	return strconv.ParseFloat(value, 64)
}

func Int(key string) (int, error) {
	value, err := getValue(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

func Int64(key string) (int64, error) {
	value, err := getValue(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(value, 10, 64)
}

func DefaultString(key, defaultValue string) string {
	if value, err := String(key); err == nil {
		return value
	}
	return defaultValue
}

func DefaultBool(key string, defaultValue bool) bool {
	if value, err := Bool(key); err == nil {
		return value
	}
	return defaultValue
}

func DefaultFloat64(key string, defaultValue float64) float64 {
	if value, err := Float64(key); err == nil {
		return value
	}
	return defaultValue
}

func DefaultInt(key string, defaultValue int) int {
	if value, err := Int(key); err == nil {
		return value
	}
	return defaultValue
}

func DefaultInt64(key string, defaultValue int64) int64 {
	if value, err := Int64(key); err == nil {
		return value
	}
	return defaultValue
}
