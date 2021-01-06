package gb28181

import (
	"encoding/json"
	"fmt"
	"github.com/Monibuca/plugin-gb28181/transaction"
	"github.com/Monibuca/plugin-gb28181/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/golang-module/carbon"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Resp struct {
	ErrorCode   int         `json:"ErrorCode"`
	Message     string      `json:"Message"`
	Data        interface{} `json:"Data"`
	RefreshTime int64       `json:"RefreshTime"`
}

func makeResp(errCode int, msg string, data interface{}) []byte {
	resp, _ := json.Marshal(Resp{
		ErrorCode:   errCode,
		Message:     msg,
		Data:        data,
		RefreshTime: time.Now().Unix(),
	})
	return resp
}

func makeJsonStrResp(errCode int, msg string, data string) []byte {
	resp := fmt.Sprintf(`{
    "ErrorCode": %d,
    "Message": "%s",
    "Data": "%s",
    "RefreshTime": %d
}`, errCode, msg, data, time.Now().Unix())
	return []byte(resp)
}

func ListAll(w http.ResponseWriter, r *http.Request) {
	//sse := util.NewSSE(w, r.Context())
	var list []*Device
	Devices.Range(func(key, value interface{}) bool {
		list = append(list, value.(*Device))
		return true
	})

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(makeResp(0, "ok", list))
}

func RecordInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")
	channel, err := strconv.Atoi(r.URL.Query().Get("channel"))
	if err != nil {
		w.Write(makeResp(-1, "param error", nil))
		return
	}
	start := carbon.Parse(startTime + "+08:00").Time.Unix()
	end := carbon.Parse(endTime + "+08:00").Time.Unix()
	if start <= 0 || end <= start {
		w.Write(makeResp(-1, "时间范围错误", nil))
		return
	}
	var resp string
	if v, ok := Devices.Load(id); ok {
		resp, err = v.(*Device).RecordInfo(channel, startTime, endTime)
	} else {
		w.Write(makeResp(-1, "设备不存在或未连接", nil))
		return
	}

	if err != nil {
		w.Write(makeResp(-1, "获取录像失败,查询失败:"+err.Error(), nil))
		return
	}
	w.Write(makeJsonStrResp(0, "ok", resp))
}

func Playback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")
	channel, err := strconv.Atoi(r.URL.Query().Get("channel"))

	if err != nil {
		w.Write(makeResp(-1, "param error", nil))
		return
	}
	start := carbon.Parse(startTime + "+08:00").Time.Unix()
	end := carbon.Parse(endTime + "+08:00").Time.Unix()
	if start <= 0 || end <= start {
		w.Write(makeResp(-1, "时间范围错误", nil))
		return
	}
	v, ok := Devices.Load(id)
	if ok {
		status, streamUri := v.(*Device).Playback(channel, start, end)
		if status != 200 {
			w.Write(makeResp(-1, "获取录像失败，点播失败", nil))
			return
		}
		w.Write(makeJsonStrResp(0, "ok", streamUri))
		return
	}
	w.Write(makeResp(-1, "设备不存在或未连接", nil))
	return
}

func PlayRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")
	channel, err := strconv.Atoi(r.URL.Query().Get("channel"))

	if err != nil {
		w.Write(makeResp(-1, "param error", nil))
		return
	}
	start := carbon.Parse(startTime + "+08:00").Time.Unix()
	end := carbon.Parse(endTime + "+08:00").Time.Unix()
	if start <= 0 || end <= start {
		w.Write(makeResp(-1, "时间范围错误", nil))
		return
	}

	var resp string
	v, ok := Devices.Load(id)

	if ok {
		resp, err = v.(*Device).RecordInfo(channel, startTime, endTime)
	} else {
		w.Write(makeResp(-1, "设备不存在或未连接", nil))
		return
	}

	if err != nil {
		w.Write(makeResp(-1, "获取录像失败，查询失败1:"+err.Error(), nil))
		return
	}

	var tmp struct {
		Response struct {
			SumNum string `json:"SumNum"`
		} `json:"Response"`
	}

	json.Unmarshal([]byte(resp), &tmp)
	if tmp.Response.SumNum == "0" {
		w.Write(makeResp(-1, "没有录像", nil))
		return
	}
	status, streamUri := v.(*Device).Playback(channel, start, end)
	if status != 200 {
		w.Write(makeResp(-1, "获取录像失败，点播失败2", nil))
		return
	}
	w.Write(makeResp(0, "ok", streamUri))
}

func CatelogCallback(c *transaction.Core, d *Device) {
	if c.Config.CatelogCallback != "" {
		go func() {
			//data, _ :=   json.Marshal(d.Channels)
			data, _ :=   jsoniter.Marshal(d.Channels)
			_, err := utils.Post(c.Config.CatelogCallback+"?id="+d.ID, data, "application/json")
			if err != nil {
				log.Println("notify " + c.Config.CatelogCallback + " error:" + err.Error())
			}
		}()
	}
}

func RemoveCallback(c *transaction.Core, d *Device) {
	if c.Config.RemoveCallback != "" {
		go func() {
			data, _ := json.Marshal(d)
			_, err := utils.Post(c.Config.RemoveCallback, data, "application/json")
			if err != nil {
				log.Println("notify " + c.Config.RemoveCallback + " error:" + err.Error())
			}
		}()
	}
}

func RemoveDead(c *transaction.Core, devices *sync.Map) {
	tick := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-tick.C:
			devices.Range(func(k, v interface{}) bool {
				device := v.(*Device)
				if device.UpdateTime.Sub(device.RegisterTime) > time.Duration(c.Config.RegisterValidity)*time.Second {
					devices.Delete(k)
					if c.Config.RemoveCallback != "" {
						go func() {
							data, _ := json.Marshal(device)
							_, err := utils.Post(c.Config.RemoveCallback, data, "application/json")
							if err != nil {
								log.Println("notify " + c.Config.RemoveCallback + " error:" + err.Error())
							}
						}()
					}
				}
				return true
			})
		}
	}

}
