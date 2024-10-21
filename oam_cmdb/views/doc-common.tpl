<!-- 导入文档框 -->
<div class="easyui-dialog formDiv" id="docImportDlg" style="width:400px;height:380px" data-options="title:'导入文档',closed: true,modal:true,
buttons:[{text:'确定',iconCls:'icon-ok',handler:docImport},
{text:'取消',handler:function(){$('#docImportDlg').dialog('close');}
}]">
    <form id="docImportFrm">
    <div class="inputctl">
     <select id="docProjId" style="width:98%"><option value="">-选择项目-</option>
        {{range $key, $val := .projItems}} <option value="{{$key}}">{{$val}}</option>{{end}}
      </select>
    </div>
    <div class="inputctl"><input name="docTitle" type="text" class="textbox" id="docTitle" style="width:98%" placeholder="请输入标题" maxlength="80" required></div>
    <div class="inputctl"><input type="file" name="docFile" id="docFile" 
    accept=".pdf,.doc,.docx,.md,application/pdf,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document" required>
    </div>
    </form>
    <div style="margin-top:20px;color:#666;">导入说明:
        <ol style="line-height:150%">
        <li>支持md,pdf,doc,docx格式文档</li>
        <li>pdf格式可在线查看,但不支持在线编辑</li>
        <li>doc,docx格式不支持在线查看和编辑</li>
        </ol>
    </div>
</div>

<div class="easyui-dialog" id="docEditDlg" style="width:1024px;height:650px" data-options="maximized:true,title:'编辑文档',closed: true,modal:true,
buttons:[{text:'保存文档',iconCls:'icon-ok',handler:function(){
    document.getElementById('docframe').contentWindow.saveArticle(false);
}
},{text:'取消',handler:function(){$('#docEditDlg').dialog('close');}
}]">
     <iframe id="docframe" style="width:100%;height:99%;border:0;"></iframe>
</div>