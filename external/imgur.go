package imgur

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type POSTResponse struct {
	Data struct {
		Link string `json:"link"`
	} `json:"data"`
}

func Upload(file *bytes.Buffer) (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://api.imgur.com/3/image", file)
	req.Header.Add("Authorization", "Client-ID "+viper.GetString("IMGUR_CLINET_ID"))
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("Imgur client error")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("Read response error")
	}
	var t POSTResponse
	err = json.Unmarshal(body, &t)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("JSON parse error")
	}
	log.Println(t.Data.Link)
	return t.Data.Link, nil
}
