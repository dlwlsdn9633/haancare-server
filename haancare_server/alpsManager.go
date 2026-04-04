package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
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

type LoginResponse struct {
	Result      bool   `json:"result"`
	AccessToken string `json:"accessToken"`
}

func GetToken() (token string, err error) {
	targetURL := "https://partner.alps.llogis.com/auth/login"
	payload := map[string]interface{}{
		"principal":  config.HanncareId,
		"credential": config.HaancarePw,
		"systmId":    "3",
		"macAddress": "normal-browser",
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(reqBody))
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", REFERER)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Wrapf(err, "status: %d, response: %v", resp.StatusCode, respBody)
		return
	}

	var loginResp LoginResponse
	if err = json.Unmarshal(respBody, &loginResp); err != nil {
		err = errors.WithStack(err)
		return
	}

	if !loginResp.Result {
		err = errors.Wrapf(err, "login failed - response: %+v", respBody)
		return
	}
	token = loginResp.AccessToken
	return
}

func GetAlpsOrders(baseUrl, orderNo, token string) (ordResults []OrderResult, err error) {

	filterData := OrderSearchFilter{
		SrchOrdNo:  orderNo,
		SrchCustCd: HAANCARE_CODE,
	}
	//timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	filterJSON, _ := json.Marshal(filterData)

	params := url.Values{}
	params.Add("filter", string(filterJSON))
	//params.Add("_", fmt.Sprintf("%d", timestamp))

	targetURL := fmt.Sprintf("%s?%s", baseUrl, params.Encode())
	var req *http.Request
	req, err = http.NewRequest("GET", targetURL, nil)
	if err != nil {
		err = errors.Wrap(err, "error while requesting")
		return
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Referer", REFERER)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "api call error")
		return
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	body, _ := io.ReadAll(resp.Body)

	// 401 error -> refresh session id
	if code != http.StatusOK {
		err = errors.Wrapf(err, "error status - %d, %s", code, string(body))
		return
	}

	if err = json.Unmarshal(body, &ordResults); err != nil {
		err = errors.Wrapf(err, "%v", body)
		return
	}
	return
}
