package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type DeliveryRecord struct {
	InvNo          string `json:"invNo"`          // 운송장 번호
	GdsNm          string `json:"gdsNm"`          // 상품명
	AcperNm        string `json:"acperNm"`        // 수령인 이름
	AcperTel       string `json:"acperTel"`       // 수령인 전화번호
	AcperRdnmBadr  string `json:"acperRdnmBadr"`  // 수령인 도로명 주소
	SnperNm        string `json:"snperNm"`        // 발송인 이름
	DlvYmd         string `json:"dlvYmd"`         // 배송일자
	DlvEmpNm       string `json:"dlvEmpNm"`       // 배송기사 이름
	DlvEmpScanCpno string `json:"dlvEmpScanCpno"` // 배송기사 연락처
	TotalCnt       int    `json:"totalCnt"`       // 전체 데이터 개수
	PickYmd        string `json:"pickYmd"`        // 집하일자
}

type SearchFilter struct {
	SrchPickYmd      string      `json:"srchPickYmd"`
	SrchPickYmdStrt  string      `json:"srchPickYmdStrt"`
	SrchPickYmdEnd   string      `json:"srchPickYmdEnd"`
	CboSrchCustSctCd string      `json:"cboSrchCustSctCd"`
	SrchCustCd       string      `json:"srchCustCd"`
	SrchCustNm       string      `json:"srchCustNm"`
	CboSrchWkSctCd   string      `json:"cboSrchWkSctCd"`
	JobCustCd        interface{} `json:"jobCustCd"`
	TabIdx           string      `json:"tabIdx"`
	RowCount         int         `json:"rowCount"`
	DispCount        int         `json:"dispCount"`
	PickYmd          string      `json:"pickYmd"`
	ColNm            string      `json:"colNm"`
	UstRtgSctCd      string      `json:"ustRtgSctCd"`
	FstmIstrYmd      string      `json:"fstmIstrYmd"`
	Status           string      `json:"_STATUS_"`
}

func main() {
	InitLogger()
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	StartCronJobs()
	InitRouter(e)
	e.Logger.Fatal(e.Start(":8080"))
}
