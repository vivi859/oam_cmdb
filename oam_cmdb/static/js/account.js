var fieldCache={};
//保存账号
function saveAccount(){
    var act=$('#act');
    var isValid=$('#accountFrm').form('validate');
    if(!isValid){
        return;
    }

    var formData={};
    var selectedType=$('#typeId').combobox('getValue');
    formData.typeId=selectedType==""?0:parseInt(selectedType);
    var selectedProjs=$('#projId').combobox('getValues');
    selectedProjs.remove(0);//删除'无'项
    if(selectedProjs.length>0){
        formData.projIds=toIntArray(selectedProjs);
    }
    formData.accountName=$('#accountName').textbox('getValue');
    formData.fieldUser=$('#fieldUser').textbox('getValue');
    formData.fieldPwd=$.trim($('#fieldPwd').textbox('getValue'));
    var hostId=$('#hostId').val();
    if(hostId!=""){
        formData.hostId=parseInt(hostId);
      //console.log("host:"+formData.hostId)
    }
    if(act="edit"){
        formData.accountId=parseInt($('#accountId').val());
        var oldPwd=$.trim($('#oldFieldPwd').val());
        if(oldPwd==formData.fieldPwd){
            formData.fieldRePwd=formData.fieldPwd;//没改过密码
        }else{
            formData.fieldRePwd=$.trim($('#fieldRePwd').textbox('getValue'));
        }
    }else{
        formData.fieldRePwd=$.trim($('#fieldRePwd').textbox('getValue'));
    }
    //console.log("pwd:"+formData.fieldPwd+"|repwd:"+formData.fieldRePwd);
    if(formData.fieldPwd!=""&&formData.fieldPwd!=formData.fieldRePwd){
        showError("两次密码不相同");
        return;
    }

    formData.fieldUrl=$('#fieldUrl').textbox('getValue');
    formData.fieldRemark=$('#fieldRemark').textbox('getValue');
    var relAcctIds=$('#relAccountIds').tagbox('getValues');
    if(relAcctIds.length>0){
        formData.relAccountIds=toIntArray(relAcctIds);
    }
    if(selectedType>0){
        var fields=fieldCache[formData.typeId];
        if(fields&&fields.length>0){
            var fieldOther={}
            $.each(fields,function(index,elem){
                if (elem.valueType==1) {
                    fieldOther[elem.fieldKey]=isCheck(elem.fieldKey);
                }else if(elem.valueType==2){
                    fieldOther[elem.fieldKey]=$('#'+elem.fieldKey).numberbox('getValue');
                }else if(elem.valueType==3){
                    fieldOther[elem.fieldKey]=$('#'+elem.fieldKey).combobox('getValue');
                }else{
                    fieldOther[elem.fieldKey]=$('#'+elem.fieldKey).textbox('getValue');
                }
            });
            formData.fieldOther=JSON.stringify(fieldOther);
       }
    }
    $.ajax({
        type: "POST",
        url: "/account/saveaccount",
        contentType: "application/json; charset=utf-8",
        data:  JSON.stringify(formData),
        dataType: "json",
        success: function (result) {
            handleJsonResult(result,function(data){
                parent.closeAccountDlg(true);
            })
        }
    });
}
const opsCache={}
//创建动态属性表单组件
function createOtherField(fields,selectorId){
    var wrapper=$('#'+selectorId);
    wrapper.empty();
    if (fields.length>0){
        $.each(fields,function(index,elem){
            var dataoptions="label:\""+elem.fieldName+"\",labelWidth:\"110\"";
            var eleHtml='<div class="inputctl"> <input name="'+elem.fieldKey+'" id="'+elem.fieldKey+'" class="';
            if(elem.isRequired){
                dataoptions=dataoptions+", required:true"
                eleHtml+="required "
            } 
            if(elem.valueType==1){
                eleHtml+='easyui-checkbox" value="true"';
               // console.log(elem.fieldValue);
                if(elem.fieldValue==true){
                   // console.log("---");
                    dataoptions+=',checked:true';
                }
            }else if(elem.valueType==2){
                eleHtml+='easyui-numberbox"';
                if (elem.maxLen!=0) {
                    dataoptions=dataoptions+",max:"+elem.maxLen
                }
                dataoptions+=",width:'90%'";
            }else if(elem.valueType==3){
               // opsCache[elem.fieldKey]=elem.valueRule;
                eleHtml+='easyui-combobox"';
                dataoptions+=',panelHeight:"150px",width:"90%",valueField:"id",textField:"text",data:'+elem.valueRule
            }else if(elem.valueType==4){
                eleHtml+='easyui-textbox"';
                dataoptions=dataoptions+', readonly:true,width:"65%"'
                eleHtml+=" data-options='"+dataoptions+"'"
                if(!isEmptyStr(elem.fieldValue)){
                    eleHtml+=' value="'+elem.fieldValue+'">&nbsp;<a class="btn btn-green" href="'+elem.fieldValue+'">下载</a';
                }
                eleHtml+='>&nbsp;<a class="btn btn-blue" href="#" onclick="openUploadFileWin(\''+elem.valueRule+'\',\''+elem.fieldKey+'\')">上传</a';
            }else{
                eleHtml+=(elem.isCiphertext?"easyui-passwordbox\"":"easyui-textbox\"");
                dataoptions+=',width:"90%"';
                if(elem.valueRule!=""){
                    try {
                        var reg=JSON.parse(elem.valueRule);
                        if(reg.pattern&&reg.pattern.length()>=1){
                            dataoptions=dataoptions+ ",validType:\"regex[\\\""+reg.pattern+"\\\"]\"";
                        }
                    } catch (error) {
                        console.error(error)
                    }
                }else{
                    if(elem.maxLen!=0){
                        dataoptions=dataoptions+ ",validType:\"maxLength["+elem.maxLen+"]\"";
                    }
                }
                if(elem.maxLen>=250){
                    dataoptions=dataoptions+", multiline:true"
                }
            }
            if(elem.valueType!=4)
                eleHtml+=" data-options='"+dataoptions+"'"
            if (elem.valueType!=1&&elem.valueType!=4&&!isEmptyStr(elem.fieldValue)){
                eleHtml+=' value="'+elem.fieldValue+'"'
            }
            if(elem.maxLen>250){
                eleHtml+=' style="height:60px"'
            }
            eleHtml+="></div>"
           
            wrapper.append(eleHtml)
            $.parser.parse('#'+selectorId);
        });
    }
}


function acctActionFmt(value,row,index){
    var link='';
    link=permStr("ACCOUNT_MGT:VIEW",'<a href="javascript:void(0)" onclick="showDetail('+row.accountId+')" class="btn btn-green btn-xs">查看</a>');
    if(!row.isDeleted){
        link+= permStr("ACCOUNT_MGT:EDIT",'&nbsp;<a href="javascript:void(0)" onclick="openEditAccountDlg('+row.accountId+')" class="btn btn-blue btn-xs">编辑</a>');
    }else{
        link+=permStr("ACCOUNT_MGT:EDIT",'&nbsp;<a href="#" perm="ACCOUNT_MGT:EDIT" onclick="recover('+row.accountId+')" class="btn btn-blue btn-xs">恢复</a>');
    }
    if(row.fieldPwd&&row.fieldPwd!=""){
        link+=permStr("ACCOUNT_MGT:COPYPWD",'&nbsp;<a href="javascript:void(0)" onclick="copyPwd('+row.accountId+')" class="btn btn-kermesinus btn-xs">复制密码</a>');
    }
    return link;
}
  //打开新增账号窗口
function openNewAccountDlg(projId){
    var url='/account/toaddaccount';
    if(projId&&projId>0){
        url=url+"?projId="+projId
    }
    document.getElementById("dlgframe").src=url;
    $('#accountDlg').dialog('setTitle',"新增账号");
    $('#accountDlg').dialog('open');
}
  //打开编辑账号窗口
function openEditAccountDlg(aid){
    document.getElementById("dlgframe").src='/account/toeditaccount?id='+aid;
    $('#accountDlg').dialog('setTitle',"编辑账号");
    $('#accountDlg').dialog('open');
}

function closeAccountDlg(reload){
    $('#accountDlg').dialog('close');
    if(reload){
        $('#accountGrid').datagrid('reload');
    }
}
function openUploadFileWin(accept,uploadFileFieldId){
    if (accept!=""){
        $('#uploadFileWin input[type="file"]').prop("accept",accept)
    }
    $('#uploadFileWin').dialog({
        title: '上传文件',
        width: '70%',
        height: 150,
        closed: false,
        cache: false,
        modal: true,
        top:'30%',
        buttons:[
            {text:'确定上传',iconCls:'icon-ok',handler:function(){
                uploadFieldFile(uploadFileFieldId);
               }
            },{text:'取消',handler:function(){
                $('#uploadFileWin').dialog('close');
            }
        }]
    });
    //$('#uploadFileWin').dialog('open');
    
}
function uploadFieldFile(uploadFileFieldId){
    var $fileField=document.getElementById('uploadfile');
    if($fileField.files.length==0){
        alert("请选择要上传的文件");
        return false;
    }
    var fd=new FormData();
    fd.append("uploadfile",$fileField.files[0])
    uploadFile("/uploadfile",fd,function(data){
        toast('文件上传成功');
        $('#'+uploadFileFieldId).textbox('setValue',data);
        $('#uploadFileWin').dialog('close');
    });
}
function decryptAndCopy(cryptStr){
    var pwd;
    try {
        pwd=decrypt(cryptStr)
    } catch (error) {
        console.log(error);
        toast("解密异常");
        return;
    }
    if(pwd==null){
        toast("解密异常");
        return;
    }
    try {
        var isOk= cp(pwd);
        if(isOk){
            toast("已复制");
            setTimeout(function () {
            cp("");
            }, 30000);
        }
    } catch (error) {
        showError("复制失败");
        return;
    }
}
function copyPwd(accid){
    if(!permMap['ACCOUNT_MGT:COPYPWD']){
        return;
    }
    $('#accountGrid').datagrid('selectRecord',accid);
    var selectedRow=$('#accountGrid').datagrid('getSelected');
    if(selectedRow.fieldPwd!=""){
        decryptAndCopy(selectedRow.fieldPwd);
    }
}

function showDetail(id){
    $('#relViewDlg').window('open')
    $('#relViewDlg').window('refresh', '/account/toviewaccount?id='+id);
}
