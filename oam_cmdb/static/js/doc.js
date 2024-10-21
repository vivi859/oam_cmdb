 //文档名称列内容,加入文件图标
 function docNameFmt(value,row,index){
    var icon;
    if(row.docType=="pdf")
        icon="fa-file-pdf-o"
    else if(row.docType=="doc"||row.docType=="docx")
        icon="fa-file-word-o"
    else
        icon="fa-file-o"
    return '<i class="fa '+icon+' margin-r5"></i>'+value;
}

function docActionFmt(value,row,index){
    var link="";
    if(row.docType=="md"){
        link+= '<a href="/doc/view/'+row['docId']+'" target="_blank" perm="DOC_MGT:VIEW" class="btn btn-green btn-xs">查看</a>&nbsp;';
        link+='<a href="#" onclick="openEditDocDlg('+row['docId']+')" perm="DOC_MGT:EDIT" class="btn btn-blue btn-xs">编辑</a>&nbsp;';
    }else if(row.docType=="doc"||row.docType=="docx"){
        link= '<a href="/doc/view/'+row['docId']+'" target="_blank" perm="DOC_MGT:VIEW" class="btn btn-green btn-xs">下载</a>&nbsp;';
    }else{
        link= '<a href="/doc/view/'+row['docId']+'" target="_blank" perm="DOC_MGT:VIEW" class="btn btn-green btn-xs">查看</a>&nbsp;';
    }
    return link+='<a href="#" onclick="deleteDoc('+row['docId']+')" perm="DOC_MGT:DEL" class="btn btn-kermesinus btn-xs">删除</a>';
}

function deleteDoc(id){
    $.messager.confirm('删除确认', '删除后不可恢复，确认要删除该文档吗?', function(r){
        if (r){
           $.getJSON('/doc/deletedoc?docId='+id,function(result){
                handleJsonResult(result,function(data){
                    var rowIndex= $('#doc-table').datagrid('getRowIndex',id);
                    $('#doc-table').datagrid('deleteRow',rowIndex);
                })
           })
        }
    });
}

function docImport(){
    var fd=new FormData(document.getElementById('docImportFrm'))
    if(fd.get("docTitle")==""){
        showError("请输入标题")
        $('#docTitle').focus();
        return;
    }
    if(!fd.get("docFile")){
        showError("请选择要导入的文件")
        return;
    }
    var projId=$('#docProjId').val();
    fd.append("projId",projId);
    uploadFile('/doc/importdoc',fd,function(newrow){
        $('#doc-table').datagrid('appendRow',newrow);
        $('#docImportDlg').dialog('close');
    });
}

function openEditDocDlg(docId){
    if(docId==-1){
        var projIdEle=$('#projId');
        var projId=0;
        if(projIdEle){
            projId=projIdEle.val(); 
        }
        $('#docEditDlg').dialog('setTitle',"新增文档")
        $('#docEditDlg').dialog('open');
        document.getElementById("docframe").src='/doc/toadddoc?projId='+projId;
    }else{
        $('#docEditDlg').dialog('setTitle',"编辑文档")
        $('#docEditDlg').dialog('open');
        document.getElementById("docframe").src='/doc/toeditdoc?docId='+docId;
    }
}

function openDocImportDlg(){
    //document.getElementById('docImportFrm').reset();
    var selectedProj=$('#projId');
    if(selectedProj){
        var projId=$('#selectedProj').val();
        $('#docProjId').val(projId);
    }
    
    $('#docImportDlg').dialog('open');
}

function closeEditDocDlg(){
    $('#docEditDlg').dialog('close');
}

function updateDocGrid(newrow){
    //console.log(JSON.stringify(newrow))
    var rowIndex= $('#doc-table').datagrid('getRowIndex',newrow.docId);
    //console.log("="+rowIndex)
    if (rowIndex>=0){
        $('#doc-table').datagrid('updateRow',{index: rowIndex,row: newrow});
    }else{
        $('#doc-table').datagrid('appendRow',newrow)
    }    
}