package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// 2026.03.03 popinvrgstinfo
// 2026.03.04 popinvrsrvinfo
const (
	HAANCARE_CODE = "982718"
	// BASE_URL      = "https://pid.alps.llogis.com:18210/pid/ftr/pacltrc/inner/popinvrgstinfo"
	BASE_URL = "https://pid.alps.llogis.com:18210/pid/ftr/pacltrc/inner/popinvrsrvinfo"
	REFERER  = "https://partner.alps.llogis.com/"
)

type OrderSearchFilter struct {
	SrchOrdNo  string `json:"srchOrdNo"`  // 주문번호
	SrchCustCd string `json:"srchCustCd"` // 고객코드 (한케어: 982718)
}

type OrderResult struct {
	PicshNm   string `json:"picshNm"`
	OrdNo     string `json:"ordNo"`     // 주문번호
	PickYmd   string `json:"pickYmd"`   // 집하일자
	JobCustCd string `json:"jobCustCd"` // 고객코드
	InvNo     string `json:"invNo"`     // 운송장번호
	JobCustNm string `json:"jobCustNm"` // 고객명
}

func GetAlpsOrders(orderNo, token string) (ordResults []OrderResult, err error) {
	filterData := OrderSearchFilter{
		SrchOrdNo:  orderNo,
		SrchCustCd: HAANCARE_CODE,
	}
	//timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	filterJSON, _ := json.Marshal(filterData)

	params := url.Values{}
	params.Add("filter", string(filterJSON))
	//params.Add("_", fmt.Sprintf("%d", timestamp))

	targetURL := fmt.Sprintf("%s?%s", BASE_URL, params.Encode())
	var req *http.Request
	req, err = http.NewRequest("GET", targetURL, nil)
	if err != nil {
		err = fmt.Errorf("failed to request: %w", err)
		return
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Referer", REFERER)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("api call error: %w", err)
		return
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	body, _ := io.ReadAll(resp.Body)

	// 401 error -> refresh session id
	if code != http.StatusOK {
		err = fmt.Errorf("api returned status: %d, %s", code, string(body))
		return
	}

	if err = json.Unmarshal(body, &ordResults); err != nil {
		err = fmt.Errorf("unmarshal error: %w", err)
		return
	}
	return
}
