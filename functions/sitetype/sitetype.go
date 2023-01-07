package sitetype

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
)

type SiteType struct {
	IsWordPress   bool
	IsWooCommerce bool
}

// GetSiteType Obtains site Type
//  url: site url with initial http(s)://
//  body: index body already read
func GetSiteType(siteurl string, body []byte) (st SiteType, err error) {
	st.IsWordPress = bytes.Contains(body, []byte("/wp-content/"))

	// Parse url
	var theUrl string
	parsedURL, err := url.Parse(siteurl)
	if err != nil {
		return st, err
	}

	if parsedURL.Scheme == "" {
		theUrl = "https"
	}

	theUrl += parsedURL.Scheme + "://" + parsedURL.Host + parsedURL.Path
	if !strings.HasSuffix(theUrl, "/") {
		theUrl += "/"
	}

	wcs, _ := Head(theUrl + "wp-content/plugins/woocommerce/")
	st.IsWooCommerce = wcs == 200

	return st, nil
}

// Head performs a head method and return the ststus code.
// Head follows up to 10 redirections.
func Head(url string) (status int, err error) {
	res, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	return res.StatusCode, nil
}
