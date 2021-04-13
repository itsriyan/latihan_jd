package modules

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	jdApi "go_open_api_sdk/api"
	jdApiDomain "go_open_api_sdk/api/domain"
	jdApiCommon "go_open_api_sdk/common"
	"log"
	"strconv"
	"strings"

	pdf "github.com/ledongthuc/pdf"
)

const (
	appKey      = "393cca8f22090dbdafd6ca40fecbba6b"
	appSecret   = "68b8a566c1c6659ddff94821651d20a7"
	accessToken = "6957c078437cde7de6bdffa6db1ca774"
	format      = "json"
	version     = "1.0"
	signMethod  = "hmacmd5"

	methodGetVariantList      = "com.jd.epi.ware.micro.action.WareQueryClient"
	getOrderIDListByCondition = "epi.popOrder.getOrderIdListByCondition"
	getOrderInfoByOrderID     = "epi.popOrder.getOrderInfoByOrderId"
	getOrderInfoListForBatch  = "epi.popOrder.getOrderInfoListForBatch"
	getSkuInfoBySpuID         = "com.jd.eptid.warecenter.api.ware.WarePlusClient.getSkuInfoBySpuId"
	getShopBrandList          = "epi.popShop.getShopBrandList"
	getSkuBySkuIDs            = "epi.ware.openapi.SkuApi.getSkuBySkuIds"
	getQueryShopTemplate      = "epi.freight.queryShopTemplateByShopIdApi"
	getWareInfoByShopID       = "com.jd.eptid.warecenter.api.ware.WarePlusClient.getWareTinyInfoListByVenderId"
	methodGetStock            = "epi.ware.openapi.warestock.queryWareStockList"
	methodGetVariantDetail    = "com.jd.eptid.warecenter.api.ware.WarePlusClient.getSkuInfoBySpuId"
	methodGetOrderPrint       = "epi.pop.order.print.uat1"
)

func constructCall(req *jdApiDomain.Request) (*jdApiDomain.Response, error) {

	fmt.Printf("%+v", req)
	fmt.Println()

	resp, err := jdApi.ApiManager.Call(req)
	if err != nil {
		return nil, err
	}
	fmt.Println("=================================================================")
	fmt.Printf("%s", resp)
	fmt.Println("=================================================================")
	return resp, nil
}

type GetOrderPrintResponse struct {
	Success bool            `json:"success"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Model   OrderPrintModel `json:"model"`
}

type OrderPrintModel struct {
	Data       []OrderPrintPDF `json:"data"`
	TemplateID int             `json:"templateId"`
}

type OrderPrintPDF struct {
	PDF           string `json:"pdf"`
	PreDeliveryID string `json:"preDeliveryId"`
	DeliveryID    string `json:"deliveryId"`
	OrderID       string `json:"orderId"`
}

type OrderPrint struct {
	Address       string     `json:"address"`
	Mail          string     `json:"mail"`
	FreightAmount JdidNumber `json:"freightAmount"`
	OrderSkuInfos []struct {
		SKUName         string     `json:"skuName"`
		SKUNumber       int        `json:"skuNumber"`
		SubtotalSubProm JdidNumber `json:"subtotalSubProm"`
		PerSKUPrice     JdidNumber `json:"perSkuPrice"`
		JDPrice         JdidNumber `json:"jdPrice"`
		SKUID           int        `json:"skuId"`
	} `json:"orderSkuinfos"`
	InstallmentFee string     `json:"installmentFee"`
	UserPhone      string     `json:"installmentFee"`
	BookTime       int        `json:"bookTime"`
	PaySubtotal    JdidNumber `json:"paySubtotal"`
	CustomerName   string     `json:"customerName"`
	PaymentType    string     `json:"paymentType"`
	AllSKUPrice    JdidNumber `json:"allSkuPrice"`
	InvoiceNoImg   string     `json:"invoiceNoImg"`
	VenderName     string     `json:"venderName"`
	ExpressCompany string     `json:"expressCompany"`
	PotonganHarga  JdidNumber `json:"potonganHarga"`
	PEO            string     `json:"peo"`
}

// JdidNumber currently only used in OrderPrint
type JdidNumber float64

func (value *JdidNumber) UnmarshalJSON(data []byte) error {
	asString := string(data)

	// Ignore null, like in the main JSON package.
	if asString == "null" || asString == `""` {
		return nil
	}

	asString = strings.Trim(asString, `"`)
	asString = strings.Replace(asString, ",", "", -1)

	parsedValue, err := strconv.ParseFloat(asString, 64)
	*value = JdidNumber(parsedValue)
	return err
}

func GetOrderPrint(orderId string, printType string, packageNum string) {
	now := jdApiCommon.NowWithTimeZone()

	req := &jdApiDomain.Request{
		Method:      methodGetOrderPrint,
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		Format:      format,
		Version:     version,
		SignMethod:  signMethod,
		Timestamp:   now,
		ParamJson:   orderId + "," + printType + "," + packageNum + "," + `"PDF"`,
	}
	// req.AddBizParameters("orderId", orderId)
	// req.AddBizParameters("printType", printType)
	// req.AddBizParameters("packageNum", packageNum)

	resp, err := constructCall(req)
	if err != nil {
		log.Println("error: ", err.Error())
	}

	orderPrintResponse := GetOrderPrintResponse{}
	err = json.Unmarshal([]byte(resp.Openapi_data), &orderPrintResponse)
	if err != nil {
		log.Println("error: ", err.Error())
	}
	log.Printf("%+v", resp.Openapi_data)
	log.Println("=================================================================================")
	// log.Printf("%+v", orderPrintResponse.Model.Data)
	log.Println("=================================================================================")
	// log.Printf("%+v", orderPrintResponse.Model.Data[0].PDF)
	log.Println("=================================================================================")

	dec, err := base64.StdEncoding.DecodeString(orderPrintResponse.Model.Data[0].PDF)
	if err != nil {
		panic(err)
	}

	// f, err := os.Create("myfilename.pdf")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	// if _, err := f.Write(dec); err != nil {
	// 	panic(err)
	// }
	// if err := f.Sync(); err != nil {
	// 	panic(err)
	// }
	// f.Seek(0, 0)
	// output file contents
	buff := bytes.NewReader(dec)

	r, err := pdf.NewReader(buff, int64(len(dec)))
	if err != nil {
		log.Println("1", err)
		return
	}
	s, err := readPdf(r)
	if err != nil {
		log.Println("2", err)
		return
	}
	fmt.Println(s)
	return

}

// s, _, err := pdf.ConvertPDF(r)
// if err != nil {
// 	panic(err)
// }
// // fmt.Println(s)

// ss := strings.Fields(s)
// awb := ""
// for i, v := range ss {
// 	if strings.Contains(v, "JXPOP") && len(v) > len(awb) {
// 		fmt.Println("index : ", i, " value : ", v)
// 		awb = v
// 	}
// }
// fmt.Println(awb)

func readPdf(r *pdf.Reader) (string, error) {
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			for _, word := range row.Content {
				fmt.Println(word.S)
			}
		}
	}
	return "", nil
}

func GetSkuBySpuIDs(spuID []int) {
	now := jdApiCommon.NowWithTimeZone()
	// sign := jdApiCommon.BuildApiSign(requestMethod, appKey, accessToken, now, format, version, signMethod, param, "", appSecret)

	req := &jdApiDomain.Request{
		Method:      getSkuInfoBySpuID,
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		Format:      format,
		Version:     version,
		SignMethod:  signMethod,
		Timestamp:   now,
	}
	req.AddBizParameters("spuIdList", spuID)

	resp, err := constructCall(req)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
}

func GetWareInfoByShopID(page int, pageSize int) {
	now := jdApiCommon.NowWithTimeZone()

	req := &jdApiDomain.Request{
		Method:      getWareInfoByShopID,
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		Format:      format,
		Version:     version,
		SignMethod:  signMethod,
		Timestamp:   now,
	}
	req.AddBizParameters("page", page)
	req.AddBizParameters("pageSize", pageSize)

	resp, err := constructCall(req)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
}

func GetQueryShopTemplate() {
	now := jdApiCommon.NowWithTimeZone()

	req := &jdApiDomain.Request{
		Method:      getQueryShopTemplate,
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		Format:      format,
		Version:     version,
		SignMethod:  signMethod,
		Timestamp:   now,
	}

	resp, err := constructCall(req)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
}

func GetSkuBySkuIDs(skus string) {
	now := jdApiCommon.NowWithTimeZone()

	req := &jdApiDomain.Request{
		Method:      getSkuBySkuIDs,
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		Format:      format,
		Version:     version,
		SignMethod:  signMethod,
		Timestamp:   now,
	}
	req.AddBizParameters("skuIds", skus)

	resp, err := constructCall(req)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
}

// MethodGetVariantDetail
func MethodGetVariantDetail(paramJson string) {
	now := jdApiCommon.NowWithTimeZone()

	req := &jdApiDomain.Request{
		Method:      methodGetVariantDetail,
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		Format:      format,
		Version:     version,
		SignMethod:  signMethod,
		Timestamp:   now,
		ParamJson:   paramJson,
	}

	resp, err := constructCall(req)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%+v", resp)
}

// MethodGetStock
func MethodGetStock(request []GetStockRequest) {
	now := jdApiCommon.NowWithTimeZone()

	req := &jdApiDomain.Request{
		Method:      methodGetStock,
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		Format:      format,
		Version:     version,
		SignMethod:  signMethod,
		Timestamp:   now,
	}
	req.AddBizParameters("jsonStr", request)

	resp, err := constructCall(req)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
}

// MethodGetVariantList
func MethodGetVariantList(status []int, scrollID string) {
	now := jdApiCommon.NowWithTimeZone()

	req := &jdApiDomain.Request{
		Method:      methodGetVariantList,
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		Format:      format,
		Version:     version,
		SignMethod:  signMethod,
		Timestamp:   now,
	}
	req.AddBizParameters("wareStatus", status)
	// req.AddBizParameters("scrollId", scrollID)

	resp, err := constructCall(req)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
}

// GetOrderIDListByCondition
func GetOrderIDListByCondition(paramJson string) {
	now := jdApiCommon.NowWithTimeZone()

	req := &jdApiDomain.Request{
		Method:      getOrderIDListByCondition,
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		Format:      format,
		Version:     version,
		SignMethod:  signMethod,
		Timestamp:   now,
		ParamJson:   paramJson,
	}

	resp, err := constructCall(req)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%+v", resp)
}

// GetOrderInfoByOrderID
func GetOrderInfoByOrderID(orderID string) {
	now := jdApiCommon.NowWithTimeZone()

	req := &jdApiDomain.Request{
		Method:      getOrderInfoByOrderID,
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		Format:      format,
		Version:     version,
		SignMethod:  signMethod,
		Timestamp:   now,
		ParamJson:   orderID,
	}
	// req.AddBizParameters("orderId", orderID)

	resp, err := constructCall(req)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
}

// GetOrderInfoListForBatch
func GetOrderInfoListForBatch(orderIDs []int64) {
	now := jdApiCommon.NowWithTimeZone()

	req := &jdApiDomain.Request{
		Method:      getOrderInfoListForBatch,
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		Format:      format,
		Version:     version,
		SignMethod:  signMethod,
		Timestamp:   now,
		ParamJson:   "[1027409981]",
	}
	// req.AddBizParameters("orderIdList", orderIDs)

	resp, err := constructCall(req)
	if err != nil {
		log.Println(err)
	}

	log.Println()
	log.Printf("%+v", resp)
}
