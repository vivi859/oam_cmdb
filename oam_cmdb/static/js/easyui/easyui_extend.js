/**
 * 扩展validatebox验证规则.
 * 内置规则:
 *  email：匹配 email 正则表达式规则。
 *  url：匹配 URL 正则表达式规则。
 *  length[0,100]：允许从 x 到 y 个字符。
 *  remote['http://.../action.do','paramName']：发送 ajax 请求来验证值，成功时返回 'true' 。
 * 新增规则:
 *   combox选择框必填,
 *   regex正则验证,
 *   equals两元素比较,
 *   min最小值
 *   max最大值
 *   range值范围
 *   minLength最小长度,
 *   maxLength最大长度,
 *   phone电话格式,
 *   mobile手机格式,
 *   num是否为数字和精度
 *   currency货币,
 *   integer整数,
 *   chinese中文,
 *   english英文字母,
 *   unnormal是否包含空格或非法字符,
 *   idcard中国身份证,
 *   ip ip地址,
 *   zip邮政编码
 *   letterOrNumOrline[]由字母、数字或下划线组成，长度{0}-{1}字符
 *   date时间范围
 */
$.extend($.fn.validatebox.defaults.rules, {
    password:{
        validator:function(value, param){
        	return /^[\S]{6,30}$/.test(value);
        },   
        message: '密码6-30个字符,不含空格'
    },
    username:{
        validator:function(value,param){
            return /^[a-z0-9][a-z0-9@.]{1,28}[a-z0-9]$/.test(value)
        },
        message:'用户名由小写字母,数字或_@.组成,长度3-30个字符'
    },
	combox : {  
        validator : function(value, param,missingMessage) {  
        	var vals=$('#'+param).combobox('getValues');
        	//console.log(JSON.stringify(vals));
            if(vals&&vals.length>0){  
            	if(vals.length==1&&vals[0]=="")
            		return false;
                return true;  
            }  
            return false;  
        },  
        message : "{1}"  
    },
    regex: {   
        validator: function(value, param){
        	var re = new RegExp(param[0]); 
        	return re.test(value);
        },   
        message: '值不符合要求'  
    }, 
	timeCheck: {	//验证开始时间小于结束时间
		validator: function(value, param){
		startTime2 = $(param[0]).datetimebox('getValue');
		var d1 = $.fn.datebox.defaults.parser(startTime2);
		var d2 = $.fn.datebox.defaults.parser(value);
		varify=d2>d1;
		return varify;
		},
		message: '结束时间要大于开始时间！'
		},
	 numCheck: {
				validator: function(value, param){
				preVar = $(this).parents("tr").find("input[name='nominalAmount']").val();
				return parseFloat(preVar) >= parseFloat(value);
				},
				message: '零售价要小于或等于面值'
	},
    equals : {
        validator : function(value, param) {
            return value == $(param[0]).val();
        },
        message : '值不相等'
    },
    unequalsZero : {
        validator : function(value) {
            return value !=0;
        },
        message : '值不等于零'
    },
    min:{
    	validator: function(value, param){
    		if(isNaN(param[0])&&param[0].charAt(0)=="#"){
    			return value >= $(param[0]).val();
    		}else{
    			return value >= param[0];
    		}
        },
        message: '最小值{0}'
    },
    max:{
    	validator: function(value, param){
    		if(isNaN(param[0])&&param[0].charAt(0)=="#"){
    			return value <= $(param[0]).val();
    		}else{
    			return value <= param[0];
    		}
        },
        message: '最大值{0}'
    },
    range:{
    	validator: function(value, param){
            return value>=param[0]&&value<=param[1];
        },
        message: '值范围[{0}-{1}]'
    },
    minLength: {
        validator: function(value, param){
            return value.length >= param[0];
        },
        message: '至少{0}个字符.'
    },
    maxLength: {
        validator: function(value, param){
            return value.length <= param[0];
        },
        message: '最多{0}个字符.'
    },
    phone : {// 验证电话号码
        validator : function(value) {
            return /^((\(\d{2,3}\))|(\d{3}\-))?(\(0\d{2,3}\)|0\d{2,3}-)?[1-9]\d{6,7}(\-\d{1,4})?$/i.test(value);
        },
        message : '格式不正确,请使用下面格式:010-88888888'
    },
    mobile : {// 验证手机号码
        validator : function(value) {
            return /^[1][3|5|6|7|8|9]\d{9}$/.test(value);
        },
        message : '手机号码格式不正确'
    },
    num : {// 验证是否为数字
        validator : function(value,param) {
            var isNum=!isNaN(value);
            if(isNum){
            	if(param&&value.indexOf(".")>=0){
            		var vp=value.split(".")[1].length;
            		if(vp==0)
            			return false;//.后没有数字
            		var pp=parseInt(param);
            		if(vp>0&&pp>0){
            			return vp<=pp;
            		}
            	}
            	return true;
            }else{
            	return false;
            }
        },
        message : '数字格式或精度不正确'
    },
    currency : {// 验证货币
        validator : function(value) {
            return /^\d+(\.\d+)?$/i.test(value);
        },
        message : '货币格式不正确'
    },
    integer : {// 验证整数
        validator : function(value) {
            return /^[+]?[0-9]+\d*$/i.test(value);
        },
        message : '请输入整数'
    },
    chinese : {// 验证中文
        validator : function(value) {
            return /^[\Α-\￥]+$/i.test(value);
        },
        message : '请输入中文'
    },
    english : {// 验证英语
        validator : function(value) {
            return /^[A-Za-z]+$/i.test(value);
        },
        message : '请输入英文'
    },
    unnormal : {// 验证是否包含空格和非法字符
        validator : function(value) {
            return /[^\n\r\t><!%\|=\*#`]+/i.test(value);
        },
        message : '不能有空白符或其他非法字符:<>!|%=*#`'
    },
    letterOrNumOrline : {
        validator : function(value,param) {
        	return /^[a-zA-Z0-9_]+$/i.test(value)&&value.length>=param[0]&&value.length<=param[1];
        },
        message : '由字母、数字或下划线组成，长度{0}-{1}字符'
    },
    zip : {// 验证邮政编码
        validator : function(value) {
            return /^[0-9]\d{5}$/i.test(value);
        },
        message : '邮政编码格式不正确'
    },
    /*ip : {// 验证IP地址
        validator : function(value) {
            return isIpv4(value)||isIpv6(value);
        },
        message : 'IP地址格式不正确'
    },
    idcard : {// 验证身份证
            validator : function(value) {
                return isIDcard(value);
            },
            message : '不是有效的身份证'
    },*/
    safeChar:{
    	validator : function(value) {
            return safeChar(value);
        },
        message : '不能包含特殊符号'
    },
    date: {//时间大小比较 date[param[0],param[1]] 参数可以是:日期字符串,#id,now当前时间
		validator: function(value, param){
			var dateValue = $.fn.datebox.defaults.parser(value);
			var start=null,end=null;
			if(param[0]!=''){
				if(param[0].charAt(0)=='#'){
					var tmpStart=$(param[0]).datebox('getValue')
					if(tmpStart!='')
						start=$.fn.datebox.defaults.parser(tmpStart);
				}else if(param[0]=='now'){
					start=new Date();
				}else{
					start=$.fn.datebox.defaults.parser(param[0]);
				}
			}
			if(param[1]!=''){
				if(param[1].charAt(0)=='#'){
					var tmpEnd=$(param[1]).datebox('getValue')
					if(tmpEnd!='')
						end=$.fn.datebox.defaults.parser(tmpEnd);
				}else if(param[0]=='now'){
					end=new Date();
				}else{
					end=$.fn.datebox.defaults.parser(param[1]);
				}
			}
			if(start!=null&&end!=null){
				return start<=dateValue&&end>=dateValue;
			}else if(start!=null&&end==null){
				return start<=dateValue;
			}else if(start==null&&end!=null){
				return end>=dateValue;
			}
			return true;
		},
		message: '时间不正确.'
	}
});
 
/**
 * 给文本框控件增加一个清空图标
 * @param {*} idOrJQueryObj 元素id或元素jquery对象
 */
function textBoxEnableClear(idOrJQueryObj)
{
    var theObj =$.type(idOrJQueryObj)==="string"?$("#"+idOrJQueryObj):idOrJQueryObj;
    var clearIconCls='icon-clear';
    
    //根据当前值，确定是否显示清除图标
    var showIcon = function(){
        var icon = theObj.textbox('getIcon',0);
        var opts=theObj.textbox('options');
        
        if(opts.readonly||opts.disabled){
            icon.css('display','none');
            return;
        }
        if (theObj.textbox('getValue')){
            icon.css('display','inline-block');
        } else {
            icon.css('display','none');
        }
    };
 
    theObj.textbox({
        icons:[{
            iconCls:clearIconCls,
            handler: function(e){
                theObj.textbox('clear');
                theObj.textbox('textbox').focus();
            }
        }],
        //值改变时，根据值，确定是否显示清除图标
        onChange:function(newValue,oldValue){
            showIcon();
        }
 
    }); 
 
    showIcon();
}
/**
 * 给容器下所有文本框加上清除按钮功能 
 * @param {文本框或文本框窗口id} fieldId 
 */
function enableClear(fieldId){
    var field=document.getElementById(fieldId);
    if (field.nodeName=="input") {
        $('#'+fieldId).textbox().textbox('addClearBtn');
    }else{
        var inputArr=$('#'+fieldId+" :text");
        inputArr.each(function(index,item){
            var $box=$(item);
            if ($box.hasClass("easyui-textbox")||$box.hasClass("easyui-numberbox")) {
                textBoxEnableClear($box)
            }
        });
    }
}

function textoverFmt(val,row){
    if (val&&val.length>5){
        return '<span title="' + val + '">' + val + '</span>';
    } else {
        return val;
    }
}
