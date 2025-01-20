package mop

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)
/*
These constants define URLs and a user-agent for interacting with Yahoo Finance:  
- `crumbURL`: URL to fetch a crumb for requests.  
- `cookieURL`: URL to retrieve cookies for authentication.  
- `userAgent`: User-agent string for identifying the client .  
- `euConsentURL`: URL for collecting user consent with a session ID.
*/
const crumbURL = "https://query1.finance.yahoo.com/v1/test/getcrumb"
const cookieURL = "https://finance.yahoo.com/"
const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/113.0"
const euConsentURL = "https://consent.yahoo.com/v2/collectConsent?sessionId="
// This function fetches a crumb (a specific piece of data) from a remote server. It sends an HTTP GET request to the `crumbURL` with custom headers, including cookies for authentication. The response body is read and returned as a string. If any error occurs during the request or reading the response, the function panics.
func fetchCrumb(cookies string) string {
	client := http.Client{}
	request, err := http.NewRequest("GET", crumbURL, nil)
	if err != nil {
		panic(err)
	}

	request.Header = http.Header{
		"Accept":          {"*/*"},
		"Accept-Encoding": {"gzip, deflate, br"},
		"Accept-Language": {"en-US,en;q=0.5"},
		"Connection":      {"keep-alive"},
		"Content-Type":    {"text/plain"},
		"Cookie":          {cookies},
		"Host":            {"query1.finance.yahoo.com"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-site"},
		"TE":              {"trailers"},
		"User-Agent":      {userAgent},
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return string(body[:])
}
// This function fetches the necessary cookies for making requests to Yahoo Finance. It first sends a GET request to `cookieURL` with specified headers, extracts a session ID and CSRF token from the response, and uses them to send a POST request to consent Yahoo's terms. It then retrieves and processes cookies from the second response, specifically looking for an A1 cookie. If the A1 cookie is found, it is returned; otherwise, the function panics.
func fetchCookies() string {

	client := http.Client{}
	var cookies []*http.Cookie
	request, err := http.NewRequest("GET", cookieURL, nil)
	if err != nil {
		panic(err)
	}

	request.Header = http.Header{
		"Authority":                 {"finance.yahoo.com"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"Accept-Encoding":           {"gzip, deflate, br"},
		"Accept-Language":           {"en-US,en;q=0.9"},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-User":            {"?1"},
		"Upgrade-Insecure-Requests": {"1"},
		"User-Agent":                {userAgent},
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	cookieA1 := getA1Cookie(response.Cookies())
	if cookieA1 != "" {
		return cookieA1
	}
	sessionRegex := regexp.MustCompile("sessionId=(?:([A-Za-z0-9_-]*))")
	sessionID := sessionRegex.FindStringSubmatch(response.Request.URL.RawQuery)[1]

	csrfRegex := regexp.MustCompile("gcrumb=(?:([A-Za-z0-9_]*))")
	csrfToken := csrfRegex.FindStringSubmatch(response.Request.Response.Request.URL.RawQuery)[1]

	gucsCookie := response.Request.Response.Request.Response.Cookies()
	var gucsCookieString string = ""
	for _, cookie := range gucsCookie {
		gucsCookieString += cookie.Name + "=" + cookie.Value + "; "
	}
	gucsCookieString = strings.TrimSuffix(gucsCookieString, "; ")

	if len(gucsCookie) == 0 {
		panic(err)
	}
	form := url.Values{}
	form.Add("csrfToken", csrfToken)
	form.Add("sessionId", sessionID)
	form.Add("namespace", "yahoo")
	form.Add("agree", "agree")
	request2, err := http.NewRequest("POST", euConsentURL+sessionID, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}

	contentLength := strconv.FormatInt(int64(len(form.Encode())), 10)

	request2.Header = http.Header{
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"Accept-Encoding":           {"gzip, deflate, br"},
		"Accept-Language":           {"en-US,en;q=0.9"},
		"Connection":                {"keep-alive"},
		"Cookie":                    {gucsCookieString},
		"Content-Length":            {contentLength},
		"Content-Type":              {"application/x-www-form-urlencoded"},
		"DNT":                       {"1"},
		"Host":                      {"consent.yahoo.com"},
		"Origin":                    {"https://consent.yahoo.com"},
		"Referer":                   {euConsentURL + sessionID},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"same-origin"},
		"Sec-Fetch-User":            {"?1"},
		"Upgrade-Insecure-Requests": {"1"},
		"User-Agent":                {userAgent},
	}

	response2, err := client.Do(request2)
	if err != nil {
		panic(err)
	}
	defer response2.Body.Close()
	cookies = response2.Request.Response.Request.Response.Request.Response.Cookies()
	cookieA1 = getA1Cookie(cookies)
	if cookieA1 != "" {
		return cookieA1
	} else {
		panic(err)
	}
}
// This function checks the provided cookies for one with the name "A1". If found, it returns the "A1" cookie in the format `Name=Value;`. If the "A1" cookie is not present, it returns an empty string.
// The "A1" cookie is a session or authentication cookie used by a web service - Yahoo Finance
func getA1Cookie(cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if cookie.Name == "A1" {
			return cookie.Name + "=" + cookie.Value + "; "
		}
	}
	return ""
}
