package rabbit

type QueueName map[string]string

//审核响应消息结构
type AuditBackMsg struct {
	SiteCode    string //站点
	BussUuid    string //唯一业务ID(default uuid=0)
	AuditStatus int    //审核状态(1.审核中 2.审核通过 3.审核不通过)
	AuditRemark string //审核备注
	AuditUid    int    //审核人ID
	AuditUser   string //审核人（默认为系统）
	AuditTime   int    //发起时间
}

//审核消息结构
type AuditMsg struct {
	SiteCode   string       //站点
	AuditMark  string       //消息审核模板类型
	BussUuid   string       //唯一业务ID
	BussData   BusinessData //业务审核基础数据
	Module     string       //消息来源模块
	CreateUid  int          //消息创建者UID
	CreateUser string       //消息创建者用户
	CreateTime int          //消息创建时间
}

//价格审核业务数据结构
type BusinessData struct {
	ChargePrice    float64 //计费价格
	PriceLoss      float64 //亏损金额，SOA无此数据，需审核中心计算
	CatId          int     //分类ID
	VirWhCode      string  //销售仓库
	PipelineCode   string  //网站渠道
	GoodSn         string  //SKU
	SysLabelId     int     //价格系统标签（价格类型）
	CalculatePrice float64 //计算结果价格
	Rate           float64 //利润率
	ChangeType     int     //价格变更类型：1：人工 2：系统
}