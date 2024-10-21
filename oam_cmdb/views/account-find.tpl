<!-- 账号筛选窗口 -->
<div title="账号查询" class="easyui-dialog" id="accountFindDlg" style="width: 550px; height: 460px;padding: 5px;" 
data-options="resizable:true,modal:true,closed:true,cache:false,buttons: [{
    text:'确认选择',iconCls:'icon-ok',handler:function(){
        selectRelAccounts();
        $('#accountFindDlg').dialog('close');
    }},
    {text:'取消', handler:function(){
       $('#accountFindDlg').dialog('close');
   }}]">
    <div id="toolbar" class="clear" style="margin-bottom: 10px;">
        <select id="condi-proj"><option value="0">-选择项目-</option>
        {{range $key, $val := .projs}} <option value="{{$key}}">{{$val}}</option>{{end}}
        </select>
        <select id="condi-atype"><option value="-1">-选择账号类型-</option>
        <option value="0">普通账号</option>
        {{range $key, $val := .accountTypes}}<option value="{{$key}}">{{$val}}</option>{{end}}
        </select>
        <input type="text" id="condi-name" enabled-clear class="textbox" maxlength="50">
        <button class="btn btn-green" onclick="searchAccount()">查询</button>
    </div>
    <table id="accountGrid" style="width: 98%;height: 320px;">
    </table>
</div>
<script>
    $(document).ready(function() {
      $('#accountGrid').datagrid({
            idField:'accountId',
            fitColumns:true,
            border:false,
            loadFilter:dataGridFilter,
            autoRowHeight:false,
            pageSize:15,
            pageList:[15,30,50],
            pagination:true,
            scrollbarSize:0,
            columns:[[
                {field:'accountId',title:'ID',checkbox:true},
                {field:'accountName',title:'账号名称',width:'25%'},
                {field:'fieldUser',title:'用户名',width:'30%'},
                {field:'projName',title:'所属项目',width:'25%'},
                {field:'typeName',title:'类型',width:'15%'}
            ]]
        });

       $('#accountGrid').datagrid('getPanel').addClass("lines-bottom");
       var pager = $('#accountGrid').datagrid('getPager');
       pager.pagination({showPageList:false,showRefresh: false});
    });

    function searchAccount(){
        var projId=parseInt($('#condi-proj').val());
        var typeId=parseInt($('#condi-atype').val());
        var aname=$('#condi-name').val();
        var param={"notDeleted":1,"base":true};
        if(projId>0){
            param.projId=projId;
        }
        if(typeId>=0){
            param.typeId=typeId;
        }
        if(aname!=""){
            param.accountName=aname
        }
        var grid=$('#accountGrid');
        if(isEmptyStr(grid.datagrid('options').url)){
            grid.datagrid('options').url='/account/accountpage';
        }
        $('#accountGrid').datagrid('load',param);
    }

    function selectRelAccounts(){
        var rows=$('#accountGrid').datagrid('getChecked');
        var relAccountTagbox=$('#relAccountIds');
     
        if(rows.length>0){
            var exist=false;
            var tagData=relAccountTagbox.tagbox('getData');
            var tagValues=relAccountTagbox.tagbox('getValues');;
           // console.log(tagData);
            var oid=$('#accountId').val()||0
            for(let row of rows){
                if(oid==row.accountId){
                    continue;
                }
                exist=false;
                for(let ele of tagData){
                    if(row.accountId==ele.id){
                        exist=true;
                        break;
                    }
                }
                if(!exist){
                   // console.log("row:"+row)
                    tagData.push({"id":row["accountId"],"text":row["accountName"]});
                    tagValues.push(row["accountId"])
                }
            }

            relAccountTagbox.tagbox('loadData',tagData);
            relAccountTagbox.tagbox('setValues',tagValues);
           // for(let val of tagValues){relAccountTagbox.tagbox('select',val)}
        }else{
            toast("未选择任何账号");
        }
    }
    </script>