package main

import "latihan_jd_id/modules"

func main() {
	// modules.GetSkuBySpuIDs([]int{607623544})
	// modules.GetWareInfoByShopID(1, 20)
	// modules.GetQueryShopTemplate()
	// modules.GetSkuBySkuIDs("607624117")

	// request := []modules.GetStockRequest{
	// 	{SkuId: 607892813},
	// 	{SkuId: 607918951},
	// }
	// modules.MethodGetStock(request)

	// modules.MethodGetVariantList([]int{1}, "")

	// modules.MethodGetVariantDetail("[607623544]")

	// modules.GetOrderPrint("1034095924", "1", "1")

	// businessParam := make(map[string]int, 0)
	// businessParam["startRow"] = 0
	// businessParam["createdTimeBegin"] = 1587382904087
	// businessParam["createdTimeEnd"] = 1587642104087
	// marshaledBusinessParam, err := json.MarshalIndent(businessParam, "", " ")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// modules.GetOrderIDListByCondition(string(marshaledBusinessParam))

	modules.GetOrderInfoByOrderID("1038883858")
	// modules.GetOrderInfoListForBatch([]int64{1007265274})
}
