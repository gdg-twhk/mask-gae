package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-zoo/bone"

	"github.com/cage1016/mask/model"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to mask-endpoints")
}

func StoresHandler(w http.ResponseWriter, r *http.Request) {
	bone.GetValue(r, "radius")
	bone.GetValue(r, "nw")
	bone.GetValue(r, "ne")
	bone.GetValue(r, "se")
	bone.GetValue(r, "sw")
	bone.GetValue(r, "x")
	bone.GetValue(r, "y")
	bone.GetValue(r, "max")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	code := http.StatusOK
	w.WriteHeader(code)

	stroes := []model.Stores{
		{
			Name:        "中美藥局",
			Phone:       "02 -27627468",
			Address:     "台北市松山區富錦街531號",
			MaskAdult:   42,
			MaskChild:   13,
			Updated:     `2020/02/04 18:30`,
			Available:   "星期一上午看診、星期二上午看診、星期三上午看診、星期四上午看診、星期五上午看診、星期六上午看診、星期日上午看診、星期一下午看診、星期二下午看診、星期三下午看診、星期四下午看診、星期五下午看診、星期六下午看診、星期日下午看診、星期一晚上看診、星期二晚上看診、星期三晚上看診、星期四晚上看診、星期五晚上看診、星期六晚上看診、星期日晚上看診",
			Coordinates: []float64{121.565481, 25.061285},
			Note:        "營業時間如有異動,以藥局公告為準",
		},
	}

	val, _ := json.Marshal(stroes)
	w.Write(val)
}
