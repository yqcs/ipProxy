package collection

//采集代理返回的参数
type proxyResult struct {
	Ip           string `json:"ip"`           //ip
	Port         int    `json:port`           //端口
	Agreement    string `json:agreement`      //请求协议
	Anonymous    string `json:anonymous`      //透明度
	Region       string `json:region`         //地区
	Speed        string `json:"speed"`        //响应速度
	Source       string `json:"source"`       //来源（采集资源站）
	Verification string `json:"verification"` //验证时间
}

//采集代理所需的参数
type proxyParamet struct {
	ipIndex           int `json:"ipIndex"`           //ip下标
	portIndex         int `json:"portIndex"`         //端口下标
	agreementIndex    int `json:"agreementIndex"`    //请求协议下标
	anonymousIndex    int `json:"anonymousIndex"`    //透明度下标
	regionIndex       int `json:"regionIndex"`       //地区下标
	speedIndex        int `json:"speedIndex"`        //响应速度下标
	sourceIndex       int `json:"sourceIndex"`       //来源（采集资源站）下标
	verificationIndex int `json:"verificationIndex"` //验证时间下标
}
