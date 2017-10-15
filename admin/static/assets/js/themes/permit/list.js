function deletePermit(act,g){
	$.post(act,null,function(r){
		o=eval('('+r+')');
		alert(o.message);	
		if(parseInt(o.code)==0){
			alert(o.message);		
		}else{
			
		}
	});
	return;
}
function checkboxManager(){
	//iCheck for checkbox and radio inputs
        $('input[type="checkbox"]').iCheck({
            checkboxClass: 'icheckbox_minimal-blue',
            radioClass: 'iradio_minimal-blue'
        });

        //When unchecking the checkbox
        $("#check-all").on('ifUnchecked', function (event) {
            //Uncheck all checkboxes
            $("input[type='checkbox']", ".table-mailbox").iCheck("uncheck");
        }).on('ifChecked', function (event) {//When checking the checkbox
            //Check all checkboxes
            $("input[type='checkbox']", ".table-mailbox").iCheck("check");
        });
}

$(function() {

	checkboxManager();
	
//     $("#example1").dataTable();
//     $('#example2').dataTable({
//         "bPaginate": true,
//         "bLengthChange": false,
//         "bFilter": false,
//         "bSort": true,
//         "bInfo": true,
//         "bAutoWidth": false
//     });
 });