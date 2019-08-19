package godingtalk

import (
	"errors"
	"fmt"
	"net/url"
)

type User struct {
	OAPIResponse
	Userid          string
	IsSys           bool `json:"is_sys"`
	SysLevel        int  `json:"sys_level"`
	Order           int
	IsLeader        bool
	Unionid         string      `json:"unionid"`         // 员工在当前开发者企业账号范围内的唯一标识，系统生成，固定值，不会改变
	Name            string      `json:"name"`            // 员工名字
	TEL             string      `json:"tel"`             // 分机号（仅限企业内部开发调用）
	WorkPlace       string      `json:"workPlace"`       // workPlace 办公地点
	Remark          string      `json:"remark"`          // remark 备注
	Mobile          string      `json:"mobile"`          // mobile 手机号码
	Email           string      `json:"email"`           // email 员工的电子邮箱
	OrgEmail        string      `json:"orgEmail"`        // orgEmail 员工的企业邮箱，如果员工已经开通了企业邮箱，接口会返回，否则不会返回
	Active          bool        `json:"active"`          // active 是否已经激活，true表示已激活，false表示未激活
	OrderInDepts    string      `json:"orderInDepts"`    // orderInDepts 在对应的部门中的排序，Map结构的json字符串，key是部门的Id，value是人员在这个部门的排序值
	IsAdmin         bool        `json:"isAdmin"`         // isAdmin 是否为企业的管理员，true表示是，false表示不是
	IsBoss          bool        `json:"isBoss"`          // isBoss 是否为企业的老板，true表示是，false表示不是
	IsLeaderInDepts string      `json:"isLeaderInDepts"` // isLeaderInDepts 在对应的部门中是否为主管：Map结构的json字符串，key是部门的Id，value是人员在这个部门中是否为主管，true表示是，false表示不是
	IsHide          string      `json:"isHide"`          // isHide 是否号码隐藏，true表示隐藏，false表示不隐藏
	Department      []int       `json:"department"`      // department 成员所属部门id列表
	Position        string      `json:"position"`        // position 职位信息
	Avatar          string      `json:"avatar"`          // avatar 头像url
	HiredDate       string      `json:"hiredDate"`       // hiredDate 入职时间。Unix时间戳 （在OA后台通讯录中的员工基础信息中维护过入职时间才会返回)
	Jobnumber       string      `json:"jobnumber"`       // jobnumber 员工工号
	Extattr         interface{} `json:"extattr"`         // extattr 扩展属性，可以设置多种属性（但手机上最多只能显示10个扩展属性，具体显示哪些属性，请到OA管理后台->设置->通讯录信息设置和OA管理后台->设置->手机端显示信息设置）
	IsSenior        bool        `json:"isSenior"`        // isSenior 是否是高管
	StateCode       string      `json:"stateCode"`       // stateCode 国家地区码
	Roles           []Role      `json:"roles"`           // roles 用户所在角色列表
}

//Role 角色
type Role struct {
	ID        string `json:"id"`        //角色id
	Name      string `json:"name"`      //角色名称
	GroupName string `json:"groupName"` //角色组名称
}

type UserList struct {
	OAPIResponse
	HasMore  bool
	Userlist []User
}

type Department struct {
	OAPIResponse
	Id                    int
	Name                  string
	ParentId              int
	Order                 int
	DeptPerimits          string
	UserPerimits          string
	OuterDept             bool
	OuterPermitDepts      string
	OuterPermitUsers      string
	OrgDeptOwner          string
	DeptManagerUseridList string
	SourceIdentifier      string `json:"sourceIdentifier"` //	部门标识字段，开发者可用该字段来唯一标识一个部门，并与钉钉外部通讯录里的部门做映射
}

type DepartmentList struct {
	OAPIResponse
	Departments []Department `json:"department"`
}

// DepartmentList is 获取部门列表
func (c *DingTalkClient) DepartmentList() (DepartmentList, error) {
	var data DepartmentList
	err := c.httpRPC("department/list", nil, nil, &data)
	return data, err
}

//DepartmentDetail is 获取部门详情
func (c *DingTalkClient) DepartmentDetail(id int) (Department, error) {
	var data Department
	params := url.Values{}
	params.Add("id", fmt.Sprintf("%d", id))
	err := c.httpRPC("department/get", params, nil, &data)
	return data, err
}

//UserList is 获取部门成员
func (c *DingTalkClient) UserList(departmentID, offset, size int) (UserList, error) {
	var data UserList
	if size > 100 {
		return data, fmt.Errorf("size 最大100")
	}

	params := url.Values{}
	params.Add("department_id", fmt.Sprintf("%d", departmentID))
	params.Add("offset", fmt.Sprintf("%d", offset))
	params.Add("size", fmt.Sprintf("%d", size))
	err := c.httpRPC("user/listbypage", params, nil, &data)
	return data, err
}
//DeptMember is 获取部门人员列表id
func (c *DingTalkClient) DeptMember(id int) ([]string, error) {
	params := url.Values{}
	params.Add("deptId", fmt.Sprintf("%d", id))
	var data struct {
		OAPIResponse
		UserIds []string `json:"userIds"`
	}
	err := c.httpRPC("/user/getDeptMember", params, nil, &data)
	if err != nil {
		return nil, err
	}
	if data.ErrCode != 0 {
		return nil, errors.New(data.ErrMsg)
	}
	return data.UserIds, nil
}

//CreateChat is
func (c *DingTalkClient) CreateChat(name string, owner string, useridlist []string) (string, error) {
	var data struct {
		OAPIResponse
		Chatid string
	}
	request := map[string]interface{}{
		"name":       name,
		"owner":      owner,
		"useridlist": useridlist,
	}
	err := c.httpRPC("chat/create", nil, request, &data)
	return data.Chatid, err
}

//UserInfoByCode 校验免登录码并换取用户身份
func (c *DingTalkClient) UserInfoByCode(code string) (User, error) {
	var data User
	params := url.Values{}
	params.Add("code", code)
	err := c.httpRPC("user/getuserinfo", params, nil, &data)
	return data, err
}

//UseridByUnionId 通过UnionId获取玩家Userid
func (c *DingTalkClient) UseridByUnionId(unionid string) (string, error) {
	var data struct {
		OAPIResponse
		UserID string `json:"userid"`
	}

	params := url.Values{}
	params.Add("unionid", unionid)
	err := c.httpRPC("user/getUseridByUnionid", params, nil, &data)
	if err != nil {
		return "", err
	}

	return data.UserID, err
}
