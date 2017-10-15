$.extend({
    load: function(url,params,callback) {
       $.post(url,params,function(r){
		callback(r);
		});
	  retrun;
    }
})