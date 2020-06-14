/**
 * Ahathallan  -  Angularjs Admin
 *
 */
(function () {
    angular.module('inspinia', [
        'ui.router',                    // Routing
        'oc.lazyLoad',                  // ocLazyLoad
        'ui.bootstrap',                 // Ui Bootstrap
        'pascalprecht.translate',       // Angular Translate
        'ngIdle',                       // Idle timer
        'ngSanitize'                    // ngSanitize
    ])
})();

// Other libraries are loaded dynamically in the config.js file using the library ocLazyLoad

//setInterval(function () {
//		$.ajax({
//				type:'get',//jquey是不支持post方式跨域的
//				async:true,
//				url:"/api/v1/get_flow_column_option", //跨域请求的URL
//				dataType:'json',
//				success : function(result){
//					if (result) {
//						$('#myModalLabel').html(result[0][0]);
//						$('#modelText').html(result[0][0]);
//						$('#btnWatch').click(function(){
//                            $(location).attr('href', 'http://localhost/#/ui/inspinia_notifications');
//						});
//						$('#myModal').modal();
//
//						setTimeout(function () {
//							 $('#myModal').modal('hide');
//						}, 2000);
//					};
//				},
//				error:function(){
//					//alert('fail');
//				}
//			});
//	}, 10000);