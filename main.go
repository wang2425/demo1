package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// 定义一个APN结构体，包含APN名称、流量上限和流量使用量2
type APN struct {
	Name       string    // APN名称
	DataLimit  int       // 流量上限（单位：KB）
	DataUsage  int       // 流量使用量（单位：KB）
	Expiration time.Time // 到期时间
}

// 定义一个Sim结构体，包含三码信息、流量数据、到期时间、卡状态、以及APN信息
type Sim struct {
	ICCID      string    // ICCID信息
	IMSI       string    // IMSI信息
	MSISDN     string    // MSISDN信息
	DataUsage  int       // 流量使用量（单位：KB）
	DataLimit  int       // 流量上限（单位：KB）
	Expiration time.Time // 到期时间
	Status     string    // 卡状态（未激活，激活，停用）
	APNs       []APN     // SIM卡APN信息
}

// 更改SIM状态
func (s *Sim) Changestatus(NewStatus string) error {
	if s.Status == "激活" {
		return errors.New("该卡已激活，请勿重新激活")
	}
	if s.Status == "未激活" && NewStatus == "激活" {
		s.Status = NewStatus
		return nil
	}
	return errors.New("无法激活")
}

// 更改流量上限
func (s *Sim) UpdateDateLimit(NewDateLimit map[string]int) {
	for i, apn := range s.APNs {
		if NewDateLimit, ok := NewDateLimit[apn.Name]; ok {
			s.APNs[i].DataLimit = NewDateLimit
		}
	}

	maxDataLimit := 0
	for _, apn := range s.APNs {
		if apn.DataLimit > maxDataLimit {
			maxDataLimit = apn.DataLimit
		}
	}
	s.DataLimit = maxDataLimit
}

// 变更流量使用
func (s *Sim) UpdateDateusage(NewUpDateusage map[string]int) {
	s.DataUsage = 0
	for i, apn := range s.APNs {
		if NewUpDateusage, ok := NewUpDateusage[apn.Name]; ok {
			s.APNs[i].DataUsage = NewUpDateusage
			if s.APNs[i].DataUsage > s.APNs[i].DataLimit {
				s.Status = "停用"
				fmt.Printf("该卡流量已经达到最大使用量,该卡已%s\n", s.Status)
			}
			s.DataUsage = s.DataUsage + s.APNs[i].DataUsage
		}
	}

}

// 修改套餐到期时间
func (s *Sim) UpadteExpiration(NewExpiration map[string]time.Time) {
	maxExpiration := s.Expiration
	for i, apn := range s.APNs {
		if NewExpiration, ok := NewExpiration[apn.Name]; ok {
			s.APNs[i].Expiration = NewExpiration
		}
	}
	for _, apn := range s.APNs {
		if apn.Expiration.After(maxExpiration) {
			maxExpiration = apn.Expiration
		}
	}
	s.Expiration = maxExpiration
}

// 判断SIM卡是否到期
func (s *Sim) DetermineExpiration(AExpiration time.Time) {
	if AExpiration.After(s.Expiration) {
		s.Status = "停用"
		fmt.Printf("该卡流量已经到期,该卡已%s\n", s.Status)
	}
}

func main() {
	// 创建一个SIM卡实例
	simCard := Sim{
		ICCID:      "1234567890",
		IMSI:       "20112001238",
		MSISDN:     "+8619556234196",
		DataUsage:  1024,                                          // 当前流量使用量
		DataLimit:  10240,                                         // 流量上限
		Expiration: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC), // 到期时间
		Status:     "未激活",                                         // 卡状态
		APNs: []APN{
			{
				Name:      "apn1",
				DataLimit: 5120,
				DataUsage: 256,
			},
			{
				Name:      "apn2",
				DataLimit: 1024,
				DataUsage: 128,
			},
		},
	}

	for {
		var test string
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("  \n  打印SIM卡信息1\n  打印APN信息2\n  更改SIM卡状态3\n  更新流量上限4\n  更改已使用流量并判断是否达到最大流量上限5\n  更改的SIM到期时间6\n  判断Sim卡是否到期7\n  请输入指令: \n")
		test, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("输入出错:", err)
			return
		}
		test = strings.TrimSpace(test)
		switch test {
		case "1":
			{
				// 打印SIM卡信息
				fmt.Printf("SIM卡信息: \n")
				fmt.Printf("ICCID: %s\n", simCard.ICCID)
				fmt.Printf("IMSI: %s\n", simCard.IMSI)
				fmt.Printf("MSISDN: %s\n", simCard.MSISDN)
				fmt.Printf("Data Usage: %d KB\n", simCard.DataUsage)
				fmt.Printf("Data Limit: %d KB\n", simCard.DataLimit)
				fmt.Printf("Expiration: %s\n", simCard.Expiration.Format("2006-01-02"))
				fmt.Printf("Status: %s\n", simCard.Status)
			}

		case "2":
			{
				// 打印APN信息
				for _, apn := range simCard.APNs {
					fmt.Printf("APN Name: %s, Data Limit: %d KB, Data Usage: %d KB\n", apn.Name, apn.DataLimit, apn.DataUsage)
				}
			}

		case "3":
			{
				//更改SIM卡状态
				var getchangestatus string
				fmt.Print("请输入新的状态：")
				fmt.Scan(&getchangestatus)
				err := simCard.Changestatus(getchangestatus)
				if err != nil {
					fmt.Println("更改状态失败：", err)
				} else {
					fmt.Printf("SIM卡状态已更改为：%s\n", simCard.Status)
				}
			}
		case "4":
			{
				//更新流量上限
				var t1, t2 int
				fmt.Print("请输入新的流量上限:\n")
				fmt.Scan(&t1, &t2)
				//fmt.Print(t1)
				//fmt.Print(t2)
				newLimits := map[string]int{
					"apn1": t1, // 新的流量上限
					"apn2": t2, // 新的流量上限
				}
				// 更新流量上限
				simCard.UpdateDateLimit(newLimits)
				// 打印新的SIM卡和APN的流量上限
				fmt.Printf("SIM卡新的流量上限：%d KB\n", simCard.DataLimit)
				for _, apn := range simCard.APNs {
					fmt.Printf("APN %s 的流量上限：%d KB\n", apn.Name, apn.DataLimit)
				}
			}

		case "5":
			{
				var t1, t2 int
				fmt.Print("请输入已使用的流量:\n")
				fmt.Scan(&t1, &t2)
				//fmt.Print(t1)
				//fmt.Print(t2)
				newUpDateusage := map[string]int{
					"apn1": t1, // 新的流量
					"apn2": t2, // 新的流量
				}
				// 更新流量
				simCard.UpdateDateusage(newUpDateusage)
				// 打印新的SIM卡和APN的已使用流量
				fmt.Printf("SIM卡已使用的流量: %d KB\n", simCard.DataUsage)
				for _, apn := range simCard.APNs {
					fmt.Printf("APN %s 的已使用流量：%d KB\n", apn.Name, apn.DataUsage)
				}
			}

		case "6":
			{
				var year1, year2 int
				var month1, month2 int
				var day1, day2 int
				var hour1, hour2 int
				var min1, min2 int
				var sec1, sec2 int
				var nsec1, nsec2 int
				fmt.Print("请输入要更改的时间:\n")
				fmt.Scan(&year1, &month1, &day1, &hour1, &min1, &sec1, &nsec1)
				fmt.Scan(&year2, &month2, &day2, &hour2, &min2, &sec2, &nsec2)
				month3 := time.Month(month1)
				month4 := time.Month(month2)
				newExpirations := map[string]time.Time{
					"apn1": time.Date(year1, month3, day1, hour1, min1, sec1, nsec1, time.UTC),
					"apn2": time.Date(year2, month4, day2, hour2, min2, sec2, nsec2, time.UTC),
				}
				// 更新APN的到期时间，并更新SIM卡的最大到期时间
				simCard.UpadteExpiration(newExpirations)
				fmt.Printf("SIM卡的到期时间: %s\n", simCard.Expiration.Format("2006-01-02"))
				for _, apn := range simCard.APNs {
					fmt.Printf("APN %s 的到期时间：%s\n", apn.Name, apn.Expiration.Format("2006-01-02"))
				}
			}
		case "7":
			{
				var s string
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("请输入时间 (格式: 年 月 日 时 分 秒 纳秒):\n")
				s, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("输入出错:", err)
					return
				}

				s = strings.TrimSpace(s)
				slice := strings.Fields(s)

				// 检查输入的切片长度是否足够
				if len(slice) < 7 {
					fmt.Println("输入格式错误，请输入: 年 月 日 时 分 秒 纳秒")
					return
				}

				// 依次将切片中的字符串转换为对应的整型
				year1, err := strconv.Atoi(slice[0])
				if err != nil {
					fmt.Println("年份转换错误:", err)
					return
				}
				month1, err := strconv.Atoi(slice[1])
				if err != nil {
					fmt.Println("月份转换错误:", err)
					return
				}
				day1, err := strconv.Atoi(slice[2])
				if err != nil {
					fmt.Println("日期转换错误:", err)
					return
				}
				hour1, err := strconv.Atoi(slice[3])
				if err != nil {
					fmt.Println("小时转换错误:", err)
					return
				}
				min1, err := strconv.Atoi(slice[4])
				if err != nil {
					fmt.Println("分钟转换错误:", err)
					return
				}
				sec1, err := strconv.Atoi(slice[5])
				if err != nil {
					fmt.Println("秒钟转换错误:", err)
					return
				}
				nsec1, err := strconv.Atoi(slice[6])
				if err != nil {
					fmt.Println("纳秒转换错误:", err)
					return
				}

				month2 := time.Month(month1)
				aExpiration := time.Date(year1, month2, day1, hour1, min1, sec1, nsec1, time.UTC)
				fmt.Printf("转换后的时间为: %v\n", aExpiration)
			}
		}
	}
}
