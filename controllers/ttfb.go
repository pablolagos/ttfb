package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/tcnksm/go-httpstat"
	"gopkg.in/macaron.v1"
)

func CheckTTFB(ctx *macaron.Context) {
	var theUrl string

	ctx.Req.ParseForm()
	site := ctx.Req.Form.Get("url")

	// Parse url
	parsedURL, err := url.Parse(site)
	if err != nil {
		ctx.Redirect("/")
		return
	}

	if parsedURL.Scheme == "" {
		theUrl = "https"
	}

	theUrl += parsedURL.Scheme + "://" + parsedURL.Host + parsedURL.Path

	result, err := checkTTFB(theUrl)
	if err != nil {
		ctx.Redirect("/")
		//ctx.Write([]byte("Error: " + err.Error()))
		return
	}

	ttfb := result.Result.ServerProcessing / time.Millisecond
	var q string
	var stars int
	if ttfb >= 1000 {
		q = "Poor performance!<br><br>The server takes more than 1 second to respond, which penalizes its position in search engine and delays the loading of the rest of the website's elements."
		stars = 1
	} else if ttfb >= 600 {
		q = "Too slow!<br><br>The server takes more than 600ms to respond. This negatively affects the website position in search rankings and makes visitors disappointed with the speed of the website. This must be fixed."
		stars = 2
	} else if ttfb >= 500 {
		q = "Slow site<br><br>.The server takes more than 500ms to respond. The site is considered slow by search engines, which lowers its position in the rankings. This must be fixed."
		stars = 3
	} else if ttfb >= 200 {
		stars = 4
		q = "The site takes more than 200 milliseconds to respond. It can still be considered a slow site and is unlikely to reach the first places in the rankings."
	} else {
		stars = 5
		q = "Congratulations!<br><br>The website loads very fast, which improves the search engine ranking."
	}

	var totaltime = result.Result.DNSLookup + result.Result.Connect + result.Result.TLSHandshake + result.Result.ServerProcessing + result.TransferTime

	//var r strings.Builder
	//r.WriteString("Checking url:" + theUrl + "\n")
	//r.WriteString(fmt.Sprintf("%+v\n", result))
	ctx.Data["res"] = fmt.Sprintf("%+v\n", result)
	ctx.Data["result"] = result
	ctx.Data["res_dns"] = result.Result.DNSLookup / time.Millisecond
	ctx.Data["res_connect"] = result.Result.Connect / time.Millisecond
	ctx.Data["res_tls"] = result.Result.TLSHandshake / time.Millisecond
	ctx.Data["res_wait"] = result.Result.ServerProcessing / time.Millisecond
	ctx.Data["res_transfer"] = result.TransferTime / time.Millisecond
	ctx.Data["w1"] = result.Result.DNSLookup.Seconds() / totaltime.Seconds() * 100
	ctx.Data["w2"] = result.Result.Connect.Seconds() / totaltime.Seconds() * 100
	ctx.Data["w3"] = result.Result.TLSHandshake.Seconds() / totaltime.Seconds() * 100
	ctx.Data["w4"] = result.Result.ServerProcessing.Seconds() / totaltime.Seconds() * 100
	ctx.Data["w5"] = result.TransferTime.Seconds() / totaltime.Seconds() * 100

	ctx.Data["stars"] = stars
	ctx.Data["qualification"] = q
	ctx.Data["url"] = theUrl

	// Get site type
	//tp, _ := sitetype.GetSiteType(theUrl, result.Body)
	ctx.Data["isWordPress"] = bytes.Contains(result.Body, []byte("/wp-content/"))
	//ctx.Data["isWooCommerce"] = tp.IsWooCommerce

	ctx.HTML(200, "ttfb/results")

}

type TTFBResult struct {
	Result       httpstat.Result
	TransferTime time.Duration
	Body         []byte
}

func checkTTFB(url string) (result TTFBResult, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	context := httpstat.WithHTTPStat(req.Context(), &result.Result)
	req = req.WithContext(context)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	start := time.Now()
	result.Body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return result, err
	}
	// _, err := io.Copy(ioutil.Discard, res.Body)
	result.TransferTime = time.Now().Sub(start)
	res.Body.Close()
	result.Result.End(time.Now())
	//tt = result.ContentTransfer(time.Now())

	return result, nil
}
