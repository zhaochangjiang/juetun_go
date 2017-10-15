function changeUppid(o){
	var afterObj= o.parent();
	afterObj.nextAll().remove();
	var upid=o.val();
	if(upid=='-1'){return;}
	$.load($("#loadUppid").attr("loadUppid"),{pid:upid},function(obj){
		if(parseInt(obj.code)>0){
			$.alert(obj.message);
			return;
		}
		if(typeof(obj.data[upid])=='undefined'||obj.data[upid].length==0){
			return;
		}
		var s= '<div class="col-xs-2"><select class="form-control" name="uppid[]"  onchange="changeUppid($(this));return;"><option value="-1">--请选择--</option>';
		for(var i in obj.data[upid]){
			s+='<option value="'+obj.data[upid][i].Id+'">'+obj.data[upid][i].Name+'</option>';
		}
		s+='</select></div>';
		afterObj.after(s);
		return;
	});
}

