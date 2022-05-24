package controllers

import (
	"fmt"
	"github.com/tcnksm/go-httpstat"
	"gopkg.in/macaron.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
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
		theUrl = "https://"
	}

	theUrl += parsedURL.Scheme + parsedURL.Host + parsedURL.Path

	result, transferTime, err := checkTTFB(theUrl)
	if err != nil {
		ctx.Write([]byte("Error: " + err.Error()))
		return
	}

	ttfb := result.ServerProcessing / time.Millisecond
	var q string
	var stars int
	if ttfb >= 600 {
		q = "Too slow!<br><br>The server takes more than 600ms to respond. This negatively affects the website position in search rankings and makes visitors disappointed with the speed of the website. This must be fixed."
		stars = 1
	} else if ttfb >= 500 {
		q = "Slow site<br><br>.The server takes more than 500ms to respond. The site is considered slow by search engines, which lowers its position in the rankings. This must be fixed."
		stars = 2
	} else if ttfb >= 200 {
		q = "Needs Improvement"
		stars = 3
	} else {
		stars = 5
		q = "Congratulations!<br><br>The website loads very fast, which improves the search engine ranking."
	}

	var totaltime = result.DNSLookup + result.Connect + result.TLSHandshake + result.ServerProcessing + transferTime

	//var r strings.Builder
	//r.WriteString("Checking url:" + theUrl + "\n")
	//r.WriteString(fmt.Sprintf("%+v\n", result))
	ctx.Data["res"] = fmt.Sprintf("%+v\n", result)
	ctx.Data["result"] = result
	ctx.Data["res_dns"] = result.DNSLookup / time.Millisecond
	ctx.Data["res_connect"] = result.Connect / time.Millisecond
	ctx.Data["res_tls"] = result.TLSHandshake / time.Millisecond
	ctx.Data["res_wait"] = result.ServerProcessing / time.Millisecond
	ctx.Data["res_transfer"] = transferTime / time.Millisecond
	ctx.Data["w1"] = result.DNSLookup.Seconds() / totaltime.Seconds() * 100
	ctx.Data["w2"] = result.Connect.Seconds() / totaltime.Seconds() * 100
	ctx.Data["w3"] = result.TLSHandshake.Seconds() / totaltime.Seconds() * 100
	ctx.Data["w4"] = result.ServerProcessing.Seconds() / totaltime.Seconds() * 100
	ctx.Data["w5"] = transferTime.Seconds() / totaltime.Seconds() * 100

	ctx.Data["stars"] = stars
	ctx.Data["qualification"] = q
	ctx.Data["url"] = theUrl
	ctx.HTML(200, "ttfb/results")

}

func checkTTFB(url string) (result httpstat.Result, transferTime time.Duration, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	context := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(context)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return result, 0, err
	}
	start := time.Now()
	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		return result, 0, err
	}
	tt := time.Now().Sub(start)
	res.Body.Close()
	result.End(time.Now())
	//tt = result.ContentTransfer(time.Now())

	return result, tt, nil
}
