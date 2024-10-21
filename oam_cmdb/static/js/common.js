/* 通用js */
//ajax标记(部分浏览器没有所以此处统一加入)
$.ajaxSetup({"x-requested-with":"XMLHttpRequest"});
$(document).ajaxComplete(function(event,xhr,options){
	if(xhr.status==401){
		alert("未登录或登录已失效，请重新登录");
		top.location.href=top.location.href;
	}
	else if(xhr.status==403){
		alert("无权限访问");
		return;
	}else if(xhr.status>=500){
		alert("操作失败,服务器异常");
		return;
	}
});
/**
 * 分页datagrid数据转换
 * @param result
 * @returns {___anonymous128_171}
 */
function dataGridFilter(result){
	if(result && result.status==200){
		var page = result.data;
		//console.log("filter:"+JSON.stringify(result));
		if(page){
			return {total:page.totalRow,rows:page.rows||[]};
		}else{
			return {total:0,rows:[]};
		}
	}else{
		showWarn(result.message==""?"无加载数据":result.message);
	}
}
/**
 * 无分页datagrid数据转换
 * @param result
 * @returns
 */
function dataGridNoPagerFilter(result){
	if(result && result.status==200){
		return result.data==null?[]:result.data;
	}else{
		showWarn(result.message==""?"无加载数据":result.message);
	}
}
/**
 * combo数据转换
 * @param result
 * @returns 
 */
function dataComboFilter(result){
	if(result && result.status==200){
		var pageData = result.data;
		return pageData;
	}else{
		showWarn(result.message==""?"无加载数据":result.message);
	}
}

/**
 * 显示普通消息
 * @param message
 */
function showMsg(message){
	$.messager.alert('消息',message);
}
/**
 * 显示错误消息
 * @param message
 */
function showError(message,fn){
	$.messager.alert('错误',message,'error',fn);
}
/**
 * 显示提示消息
 * @param message
 */
function showInfo(message){
	$.messager.alert('提示',message,'info');
}

function isCheck(id){
	return $('#'+id).checkbox('options').checked
}

//弹出个toast消息
function toast(message,timeout){
	var _timeout=timeout?timeout:1000;
	$.messager.show({
		id:'toast-inner',
		msg:message,
		timeout:_timeout,
		border:false,
		showType:'fade',
		showSpeed:300,
		width:'auto',height:'auto',
		style:{
			right:'',
			top:document.body.scrollTop+document.documentElement.scrollTop+10,
			bottom:''
		}
	});
}

/**
 * 显示提示消息
 * @param message 消息
 * @param top 设置面板距离顶部的位置（即Y轴位置）。如果isPercent为true,则该值为小数
 */
function showInfo(message,top,isPercent){
	//$(window).height()为浏览器当前窗口可视区域高度
	var relTop = winTopCoord(top);
	$.messager.alert({
		title : '提示',
		msg : message,
		icon : 'info',
		top : relTop
	});
}
/**
 * 显示警告消息
 * @param message
 */
function showWarn(message){
	$.messager.alert('警告',message,'warning');
}
/**
 * 显示警告消息
 * @param message 消息
 * @param top 设置面板距离顶部的位置（即Y轴位置）
 */
function showWarn(message,top){
	//$(window).height()为浏览器当前窗口可视区域高度
	var relTop = winTopCoord(top);
	$.messager.alert({
		title : '警告',
		msg : message,
		icon : 'warning',
		top : relTop
	});
}

function winTopCoord(top){
	//使用百分比
	if (typeof(top)=="string"&&top.charAt(top.length-1)=='%'){
		($(window).height()) *parseFloat(top.slice(0,-1))
	}
	return top;
}

//数组是否有重复元素
function isRepeatArray(arr) {
   var hash = {};
   for(var i in arr) {
       if(hash[arr[i]]==1)
       {
           return true;
       }
       hash[arr[i]] = 1;
    }
   return false;
}

//给数组添加一个删除某元素的方法
Array.prototype.remove = function(val) {
	var index=-1;
	for (var i = 0; i < this.length; i++) {
		if (this[i] == val) {
			index=i;
			break;
		}
	}
	if (index > -1) {
		this.splice(index, 1);
	}
};
Array.prototype.joinstr = function(separator) {
	var tmp="";
	for (var i = 0; i < this.length; i++) {
		if (this[i]&&this[i] != "") {
			tmp=tmp+this[i]
			if(i<this.length-1){
				tmp=tmp+separator;
			}
		}
	}
	return tmp;
};
function ajaxForm(selector,url,okHandler){
	$(selector).form("submit",{
		url:url,
		onSubmit: function(){
			var r=$(this).form('validate');
			//console.log("验证结果:"+r);
	        return $(this).form('validate');
	    },
		success:function(result){
			//console.log("success:"+result);
			handleJsonResult(result,okHandler)
		}
	});
}

function handleJsonResult(jsonResp,okHandler){
	//console.log(typeof jsonResp)
	if(jsonResp){
		var result=typeof(jsonResp)=="string"?$.parseJSON(jsonResp):jsonResp;
		if(result.status==200){
			if(okHandler){
				if(typeof(okHandler)=="string"){
					location.href=okHandler
				}else{
					okHandler(result.data);
				}
			}else{
				toast("执行成功");
			}
		}else{
			showError(result.message);
		}
	}else{
		alert("执行异常");
	}
}

function uploadFile(url,formData,callBack){
    $.ajax({
		url: url,
		type: 'post',
		processData: false,
		contentType: false,
		data: formData,
		success: function(result) {
			//console.log(typeof result)
			handleJsonResult(result,callBack)
		}
    })
}

function base64Encode(str){
	return window.btoa(encodeURIComponent(str))
}
function base64Dencode(str){
	return decodeURIComponent(window.atob(str));
}
//下载公钥临时存本地
function loadpwdkey(isPubkey){
	var storeKey=isPubkey?"_pubkey":"_prikey";
	var url=isPubkey?"/pubkey":"/prikey";
	var key=sessionStorage.getItem(storeKey)
	if(key==null){
		$.ajax({
			url:url,
			method:'POST',
			async:false,
			success: function(result){
				handleJsonResult(result,function(data){
					if(data!=null&&data!=""){
						key=data
						sessionStorage.setItem(storeKey,data)
					}else{
						throw "下载密钥失败"
					}
				});
			},
		})
	}
	return key
}

function decrypt(pwd){
	var encrypt=new JSEncrypt(); 
	encrypt.setPrivateKey(loadpwdkey(false));
	return encrypt.decrypt(pwd);
}
//密码简单编码
function pwdencode(pwd){
	//var encrypt=new JSEncrypt(); 
	//pubkey=loadpwdkey() 
	//encrypt.setPublicKey(pubkey);  
	//return encrypt.encrypt($.trim(pwd)); 
	var tmpStr=base64Encode(pwd);
	var tmpArr=tmpStr.split("");
	tmpStr=tmpArr.reverse().join("")+randomString(5);
	return base64Encode(tmpStr);
}

function randomString(length) {
	var str = 'abcdefghijklmnopqrstuvwxyz0123456789';
	var result = '';
	for (var i = length; i > 0; --i) 
	  result += str[Math.floor(Math.random() * str.length)];
	return result;
}
function nowTime(){
	var d = new Date(),str = '';
	str+= d.getFullYear()+'-';
	str+= d.getMonth() + 1+'-';
	str+= d.getDate()+' ';
	str+= d.getHours()+':'; 
	str+= d.getMinutes()+':'; 
	str+= d.getSeconds(); 
	return str;
}

function boolFmt(value,row,index){
	return value?"是":"否";
}
//超长列格式,鼠标放上提示全部效果
function tipFmt(value,row,index){
	if(value.length>10)
		return '<span title="'+value+'">'+value+'</span>';
	return value;
}
function ipFmt(value,row,index){
	var arr=[value,row.internalIp]
	 return arr.joinstr("<br>");
  }

//复制内容到粘贴板
function cp(text){
	var w = document.createElement('input');
	w.value = text;
	document.body.appendChild(w);
	w.select();
	// 调用浏览器的复制命令
	document.execCommand("Copy");
	w.style.display='none';
	return true;
}

function toIntArray(strArray){
	var intArr=[];
	for(var i = 0; i < strArray.length; i++) {
		intArr.push(parseInt(strArray[i]))
	}
	return intArr;
}
//不包含特殊字符
function safeChar(value){
	var re = /[~#\^\$><%!*=`]/gi;  
	return !re.test(value); 
}
function resetForm(formId){
	//document.getElementById(formId).reset();
	$('#'+formId).form('reset');
}
var hostTypes=[{"value":0,"text":"通用服务器"},
    {"value":2,"text":"Web应用服务器"},
    {"value":1,"text":"数据库服务器"},
    {"value":3,"text":"文件服务器"},
    {"value":4,"text":"保垒机"}];

function isEmptyStr(str){
	return str==null||str=="";
}

var permMap={};
function loadPerms(funcodes){
	if(funcodes.length==0){
		return;
	}
	$.ajax({
		url:"/fun/getloginuserperms",
		method:'POST',
		data:{"perms":funcodes.join(",")},
		async:false,
		success: function(result){
			handleJsonResult(result,function(data){
				if(data!=null){
					permMap=data
				}
			});
		},
		error: function(){
			alert("加载权限失败");
		}
	});
}

function filterPermElements(){
	$("button[perm]").each(function(){
		if(!permMap[this.getAttribute("perm")]){
			this.setAttribute('disabled','disabled');
		}
	});
	$("a[perm]").each(function(){
		if(!permMap[this.getAttribute("perm")]){
			$(this).hide();
		}
	});
}
function permStr(funcode, defaultStr){
	return permMap[funcode]?defaultStr:"";
}

/**
 * 获取下拉框数据
 * @param {*} selectId 下拉框id
 * @returns 
 */
function getSelectorData(selectId){
	var selectEle=document.getElementById(selectId);
	var ds={};
	for(var i=0;i<selectEle.options.length;i++){
		var opt=selectEle.options[i];
		ds[opt.value]=opt.text;
	}
	return ds;
}