package test

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"OAM/models"
	fn "OAM/util"

	jsoniter "github.com/json-iterator/go"
)

func TestRandom(t *testing.T) {
	s := fn.RandomStr(10)
	fmt.Printf("随机字符:%s", s)

	d := fn.RandomAscii(20)
	fmt.Printf("\n随机ASCII:%s", d)

	n := fn.RandomNumStr(14)
	fmt.Printf("\n随机数字串:%s", n)

	nums := fn.RandomNumbers(0, 30, 10)
	for _, v := range nums {
		fmt.Println(v)
	}
	fmt.Printf("\n随机数字:%s", fn.JoinInteger(",", nums...))

	st := fn.RandomStr(3)
	fmt.Println("\n:" + st)
	fmt.Println(len(st))
}

func TestRegx(t *testing.T) {
	var regx_uname = `^[a-z0-9][\w@.]{1,28}[a-z0-9]$`
	isOk, _ := regexp.MatchString(regx_uname, "Aadmin")
	fmt.Println(isOk)

	isOk, _ = regexp.MatchString(regx_uname, "ad")
	fmt.Println(isOk)

	isOk, _ = regexp.MatchString(regx_uname, "ad@ghg.com")
	fmt.Println(isOk)

	isip, _ := regexp.MatchString("^\\d{1,3}\\.[\\d|\\.]*", "120.24.175.215")
	fmt.Println(isip)
	mdTxt := `1. 先用 aTrust登录VPN(当前账号huanggh,密码请参见账号信息)
    ![atrust1.png](/static_ext/images/6/1656409756000npkvl.png)
2. 选择堡垒机登录(账号zhangtq)
    ![baoleiji.png](/static_ext/images/6/1656410112845infdr.png)
    
    ![baoleiji2.png](/static_ext/images/6/1656410228982affex.png)
3. 输入"二次验证码"获取的认证码`
	reg, _ := regexp.Compile(`\[.+\.\w{3,4}\]\(/static_ext/([\w/]+\.\w{3,4})\)`)
	mdR := reg.FindAllStringSubmatch(mdTxt, -1)
	if mdR != nil {
		for _, v := range mdR {
			fmt.Println(v[1])

		}
	}
}

func TestStr(t *testing.T) {
	/* strs := [...]string{"UserInfo", "USerInfo", "user", "123123user", "中23", "User_info"}
	for _, v := range strs {
		fmt.Println(fn.ToFirstLetterLower(v))
	} */

	elems := []interface{}{3, 2}
	var ftm = "%d,%d"
	str := fmt.Sprintf(ftm, elems...)
	fmt.Println(str)

	balance1 := []int{999000000004, 2, 2, 7, 49}
	fmt.Println("int join :" + fn.JoinInteger(",", balance1...))

	ex := "/login?user=x"
	fmt.Println(ex, "? before:", fn.SubBefore(ex, "?"))

}

func TestJson(t *testing.T) {
	var js = `{"Now":"2007-01-02 15:04:01","Age":24,"Passwd":"asfssecret"}`
	var d, d1 Demo
	err := json.Unmarshal([]byte(js), &d)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("原生json解析:%s \n", d.Now)

	err = jsoniter.Unmarshal([]byte(js), &d1)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("jsoniter解析:%s \n", d1.Now)
	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", "2007-01-02 15:04:01", time.Local)
	var d2 = Demo{Now: t1, Age: 23, Passwd: "secret"}
	b1, _ := jsoniter.MarshalToString(d2)
	fmt.Println(b1)

	var d3 Demo
	err = jsoniter.UnmarshalFromString(b1, &d3)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d3)

	var accType = `{"TypeId":1,"TypeName":"tname",
		"Fields":[{"fieldKey":"wx_appid","fieldName":"appid","fieldId":-1,"IsRequired":true,"IsCiphertext":false},
		{"fieldKey":"app_secret","fieldName":"密钥","fieldId":-2,"IsRequired":false,"IsCiphertext":false}]}`
	var dtype models.AccountType
	jsoniter.UnmarshalFromString(accType, &dtype)
	fmt.Println("dtype:" + dtype.TypeName)
}

func TestSlice(t *testing.T) {
	var s = ",abc,efs,"
	fmt.Println(s[1 : len(s)-1])

	var s1 = "abc,a,3,b,a"
	arr := strings.Split(s1, ",")
	t.Log(arr)
	arr1 := fn.RemoveSlice(arr, "a")
	t.Log(arr)
	t.Log(arr1)
}

type Demo struct {
	Now    time.Time `format:"2006-01-02 15:04"`
	Age    int16
	Passwd string `valid:"x" json:"-"`
}

func TestRegxStr(t *testing.T) {
	regx, _ := regexp.Compile(`[%&*=<>!'-]`)
	s := regx.ReplaceAllString("谢谢a=1 and b<>c", "")
	fmt.Println(s)
}

func TestBit(t *testing.T) {
	i := 1 << 26
	fmt.Print(i)
}

func TestHomeDir(t *testing.T) {
	storePath, err := fn.UserHomeDir()
	if err != nil {
		panic(err)
	}
	fmt.Println(path.Join(storePath, ".oam"))
	fmt.Println(path.Join(storePath, `.oam`))
	fmt.Println(filepath.Join(storePath, ".oam"))

}

func TestStrFFmt(t *testing.T) {
	tpl := `<tr><td><input type="hidden" id="hostAccountId%d" value="%s">
	<input type="text" id="hostAccountName%d" maxlength="50" class="textbox" style="width:99%%" value="%s"></td>
	<td><input data-options="validType:'maxLength[50]'" class="easyui-passwordbox" id="hostAccountPwd%d" value="%s" style="width:99%%"></td>
	</tr>`
	var html, tmpId, tmpName, tmpPwd string
	for i := 0; i < 3; i++ {
		tmpName = "name"
		tmpPwd = "pwd"
		tmpId = "id"
		html = html + fmt.Sprintf(tpl, i, tmpId, i, tmpName, i, tmpPwd)
		//html = html + fmt.Sprintf(tpl, i, tmpId)
	}
	fmt.Println(html)
}
func TestJSONMAP(t *testing.T) {
	result := make(map[string]string)
	result["isuse"] = "a"
	//b, _ := json.MarshalIndent(result, "", " ")
	b, _ := jsoniter.MarshalToString(result)
	//b, _ := jsoniter.MarshalIndent(result, "", " ")
	fmt.Println(string(b))
}
