package ebayprice

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	GLOBAL_ID_EBAY_US = "EBAY-US"
	GLOBAL_ID_EBAY_FR = "EBAY-FR"
	GLOBAL_ID_EBAY_DE = "EBAY-DE"
	GLOBAL_ID_EBAY_IT = "EBAY-IT"
	GLOBAL_ID_EBAY_ES = "EBAY-ES"
)

type Item struct {
	ItemId        string   `xml:"itemId"`
	Title         string   `xml:"title"`
	Location      string   `xml:"location"`
	CurrentPrice  float64  `xml:"sellingStatus>currentPrice"`
	ShippingPrice float64  `xml:"shippingInfo>shippingServiceCost"`
	BinPrice      float64  `xml:"listingInfo>buyItNowPrice"`
	ShipsTo       []string `xml:"shippingInfo>shipToLocations"`
	ListingUrl    string   `xml:"viewItemURL"`
	ImageUrl      string   `xml:"galleryURL"`
	Site          string   `xml:"globalId"`
}

type FindItemsByKeywordResponse struct {
	XmlName   xml.Name `xml:"findItemsByKeywordsResponse"`
	Items     []Item   `xml:"searchResult>item"`
	Timestamp string   `xml:"timestamp"`
}

type ErrorMessage struct {
	XmlName xml.Name `xml:"errorMessage"`
	Error   Error    `xml:"error"`
}

type Error struct {
	ErrorId   string `xml:"errorId"`
	Domain    string `xml:"domain"`
	Severity  string `xml:"severity"`
	Category  string `xml:"category"`
	Message   string `xml:"message"`
	SubDomain string `xml:"subdomain"`
}

type EBay struct {
	appid  string
	client *http.Client
}

func NewEBay(appid string, client *http.Client) *EBay {
	return &EBay{appid, client}
}

func (e *EBay) build_search_url(global_id string, keywords string, entries_per_page int) (string, error) {
	var u *url.URL
	u, err := url.Parse("http://svcs.ebay.com/services/search/FindingService/v1")
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("OPERATION-NAME", "findItemsByKeywords")
	params.Add("SERVICE-VERSION", "1.0.0")
	params.Add("SECURITY-APPNAME", e.appid)
	params.Add("GLOBAL-ID", global_id)
	params.Add("RESPONSE-DATA-FORMAT", "XML")
	params.Add("REST-PAYLOAD", "")
	params.Add("keywords", keywords)
	params.Add("paginationInput.entriesPerPage", strconv.Itoa(entries_per_page))
	params.Add("itemFilter(0).name", "ListingType")
	params.Add("itemFilter(0).value(0)", "FixedPrice")
	params.Add("itemFilter(0).value(1)", "AuctionWithBIN")
	u.RawQuery = params.Encode()
	return u.String(), err
}

func (e *EBay) FindItemsByKeywords(global_id string, keywords string, entries_per_page int) (FindItemsByKeywordResponse, error) {
	var response FindItemsByKeywordResponse
	url, err := e.build_search_url(global_id, keywords, entries_per_page)
	if err != nil {
		return response, err
	}
	headers := make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11"
	resp, err := e.client.Get(url)
	if err != nil {
		return response, err
	}

	status_code := resp.StatusCode
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if status_code != 200 {
		var em ErrorMessage
		err = xml.Unmarshal([]byte(body), &em)
		if err != nil {
			return response, err
		}
		return response, errors.New(em.Error.Message)
	} else {
		err = xml.Unmarshal([]byte(body), &response)
		if err != nil {
			return response, err
		}
	}
	return response, err
}
