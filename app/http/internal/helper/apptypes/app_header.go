package apptypes

// headerDevice 设备相关参数
type headerDevice struct {
	DeviceModel string `header:"X-Device-Model" form:"h_adm"` // 设备型号。如：iPhone 14 Pro
	DeviceBand  string `header:"X-Device-Band" form:"h_adb"`  // 设备厂商品牌。如：xiaomi, huawei
	OSVersion   string `header:"X-Os-Version" form:"h_aos"`   // 操作系统版本。如：ios 16.5, android 10.8
	AppVersion  string `header:"X-App-Version" form:"h_av"`   // APP 版本号，格式：app版本/RN版本，如：1.0/0.1.202039271
}

// APPHeader  通用请求 header
type APPHeader struct {
	headerI18N
	headerDevice
	headerUTM

	APPChannelFlag string `header:"X-App-Channel" form:"h_ac"` // 来源渠道
}

func (c APPHeader) ChannelID() uint {
	// todo 解析来源渠道
	return 0
}
