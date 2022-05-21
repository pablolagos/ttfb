package controllers

import (
	"fmt"
	"github.com/tcnksm/go-httpstat"
	"gopkg.in/macaron.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func CheckTTFB(ctx *macaron.Context) {

	ctx.Req.ParseForm()
	url := ctx.Req.Form.Get("url")
	ctx.Resp.Header().Set("Content-Type", "text/plain")
	result, err := checkTTFB(url)
	if err != nil {
		ctx.Write([]byte("Error: " + err.Error()))
		return
	}
	ctx.Write([]byte("Checking url:" + url + "\n"))
	ctx.Write([]byte(fmt.Sprintf("%+v\n", result)))

}

func checkTTFB(url string) (result httpstat.Result, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	context := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(context)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}

	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		return result, err
	}
	res.Body.Close()
	result.End(time.Now())

	return result, nil
}
