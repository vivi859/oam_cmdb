
function saveApp(okFun){
    var isValid=$('#appFrm').form('validate');
    if (!isValid){
        return false;
    }
  
    var appForm={};
    var hid=$('#appId').val();
    if(hid!=""){
      appForm.appId=parseInt(hid);
    }
    
    appForm.appName=$('#appName').textbox('getValue');
    appForm.appUrl=$('#appUrl').textbox('getValue');
    appForm.appDir=$('#appDir').textbox('getValue');
    appForm.sourcecodeRepo=$('#sourcecodeRepo').textbox('getValue');
    appForm.devLang=$('#devLang').combobox('getText');
    appForm.appType=$('#appType').combobox('getText');
    
    var hport=$('#appPort').numberbox('getValue');
    if(hport!=""){
      appForm.appPort=parseInt(hport);
    }
    var pid=$('#projId').combobox('getValue');
    if(pid!=""){
      appForm.projId=parseInt(pid);
    }
    var hostIds=$('#hostIds').combogrid('getValues');
    if(hostIds.length>0){
      appForm.hostIds=toIntArray(hostIds)
    }
    appForm.desc=$('#desc').textbox('getValue');
    $.ajax({
          type: "POST",
          url: "/appinfo/saveapp",
          contentType: "application/json; charset=utf-8",
          data:  JSON.stringify(appForm),
          dataType: "json",
          success: function (result) {
            if(result.status==200){
              toast('保存成功');
              if(okFun)
                okFun();
            }else{
              showError(result.message);
            }
          }
      });
  }

//保存主机
function saveHost(okFun){
  var isValid=$('#hostFrm').form('validate');
  if (!isValid){
      return false;
  }
  var isEdit=false;
  var hostForm={};
  var hid=$('#hostId').val();
  if(hid!=""){
    hostForm.hostId=parseInt(hid);
    if(hostForm.hostId>0)
      isEdit=true;
  }
  hostForm.hostName=$('#hostName').textbox('getValue');
  var selectedProjs=$('#projId').combobox('getValues');
  selectedProjs.remove(0);//删除'无'项
  if(selectedProjs.length>0){
    hostForm.projIds=toIntArray(selectedProjs);
  }

  var i=0;
  var accts=[];
  for(i=0;i<3;i++){
    var tmpacct={};
    tmpacct.fieldUser=$('#hostAccountName'+i).val();
    tmpacct.fieldPwd=$('#hostAccountPwd'+i).passwordbox('getValue');
    var id=$('#hostAccountId'+i).val();
    
    tmpacct.accountId=id==""?0:parseInt(id);
    if(tmpacct.fieldUser==""){
      if(isEdit&&tmpacct.accountId>0){
          showError("已有用户名不能修改为空,如需删除请到账号管理中操作")
          return;
      }else{
        continue;
      }
    }
    //新增的,生成个名字
    if(tmpacct.accountId==0){
      tmpacct.accountName=hostForm.hostName+"账号"+tmpacct.fieldUser;
    }
    accts.push(tmpacct);
  }
  if(accts.length>0){
    hostForm.accounts=accts;
  }

  var hType=$('#hostType').combobox('getValue');
  hostForm.hostType=hType==""?0:parseInt(hType);
  hostForm.publicIp=$('#publicIp').textbox('getValue');
  hostForm.internalIp=$('#internalIp').textbox('getValue');
  var hport=$('#sshPort').numberbox('getValue');
  if(hport!=""){
    hostForm.sshPort=parseInt(hport);
  }
  var hos=$('#osName').combobox('getValue');
  if(hos!=""){
    hostForm.osName=hos;
  }

  hostForm.desc=$('#desc').textbox('getValue');
  $.ajax({
        type: "POST",
        url: "/host/savehost",
        contentType: "application/json; charset=utf-8",
        data:  JSON.stringify(hostForm),
        dataType: "json",
        success: function (result) {
          if(result.status==200){
            toast('保存成功');
           if(okFun){
            okFun();
           }
          }else{
            showError(result.message);
          }
        }
    });
}
