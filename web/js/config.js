/**
 * Ahathallan  -  Angularjs Admin
 */
function config($stateProvider, $urlRouterProvider, $ocLazyLoadProvider, IdleProvider, KeepaliveProvider) {

    // Configure Idle settings
    IdleProvider.idle(5); // in seconds
    IdleProvider.timeout(120); // in seconds

    $urlRouterProvider.otherwise("/tables/ahathallan_modbus_scheme_lst");

    $ocLazyLoadProvider.config({
        // Set to true if you want to see what and when is dynamically loaded
        debug: false
    });

    $stateProvider

        .state('dashboards', {
            abstract: true,
            url: "/dashboards",
            templateUrl: "views/common/content.html",
        })
         .state('layouts', {
             url: "/layouts",
             templateUrl: "views/layouts.html",
             data: { pageTitle: 'Layouts' },
         })
         .state('forms', {
            abstract: true,
            url: "/forms",
            templateUrl: "views/common/content.html",
        })
        .state('tables', {
            abstract: true,
            url: "/tables",
            templateUrl: "views/common/content.html"
        })
        .state('tables.ahathallan_modbus_value', {
            url: "/ahathallan_modbus_value",
            templateUrl: "views/ahathallan_modbus_value.html",
            data: { pageTitle: '通用Modbus数据点' },
            resolve: {
                loadPlugin: function ($ocLazyLoad) {
                    return $ocLazyLoad.load([
                        {
                            files: ['css/plugins/iCheck/custom.css','js/plugins/iCheck/icheck.min.js','js/plugins/pace/pace.min.js','css/plugins/ionRangeSlider/ion.rangeSlider.css','css/plugins/ionRangeSlider/ion.rangeSlider.skinFlat.css','js/plugins/ionRangeSlider/ion.rangeSlider.min.js']
                        },
                        {
                            files: ['css/plugins/touchspin/jquery.bootstrap-touchspin.min.css', 'js/plugins/touchspin/jquery.bootstrap-touchspin.min.js']
                        }
                    ]);
                }
            }
        })
        .state('tables.ahathallan_modbus_set', {
            url: "/ahathallan_modbus_set",
            templateUrl: "views/ahathallan_modbus_set.html",
            data: { pageTitle: '通用Modbus数据设置' },
            resolve: {
                loadPlugin: function ($ocLazyLoad) {
                    return $ocLazyLoad.load([
                        {
                            files: ['css/plugins/iCheck/custom.css','js/plugins/iCheck/icheck.min.js','js/plugins/pace/pace.min.js','css/plugins/ionRangeSlider/ion.rangeSlider.css','css/plugins/ionRangeSlider/ion.rangeSlider.skinFlat.css','js/plugins/ionRangeSlider/ion.rangeSlider.min.js']
                        },
                        {
                            files: ['css/plugins/touchspin/jquery.bootstrap-touchspin.min.css', 'js/plugins/touchspin/jquery.bootstrap-touchspin.min.js']
                        }
                    ]);
                }
            }
        })
        .state('tables.ahathallan_modbus_scheme_lst', {
            url: "/ahathallan_modbus_scheme_lst",
            templateUrl: "views/ahathallan_modbus_scheme_lst.html",
            data: { pageTitle: 'Modbus方案列表' },
            resolve: {
                loadPlugin: function ($ocLazyLoad) {
                    return $ocLazyLoad.load([
                        {
                            files: ['js/plugins/pace/pace.min.js']
                        },
                        {
                            files: ['js/plugins/sweetalert/sweetalert.min.js', 'css/plugins/sweetalert/sweetalert.css']
                        },
                        {
                            name: 'oitozero.ngSweetAlert',
                            files: ['js/plugins/sweetalert/angular-sweetalert.min.js']
                        }
                    ]);
                }
            }
        })
        .state('tables.ahathallan_modbus_scheme_info', {
            url: "/ahathallan_modbus_scheme_info",
            templateUrl: "views/ahathallan_modbus_scheme_info.html",
            data: { pageTitle: 'Modbus方案信息' },
            resolve: {
                loadPlugin: function ($ocLazyLoad) {
                    return $ocLazyLoad.load([
                        {
                            files: ['js/plugins/pace/pace.min.js']
                        },
                        {
                            files: ['css/plugins/touchspin/jquery.bootstrap-touchspin.min.css', 'js/plugins/touchspin/jquery.bootstrap-touchspin.min.js']
                        },
                        {
                            files: ['js/plugins/sweetalert/sweetalert.min.js', 'css/plugins/sweetalert/sweetalert.css']
                        },
                        {
                            name: 'oitozero.ngSweetAlert',
                            files: ['js/plugins/sweetalert/angular-sweetalert.min.js']
                        }
                    ]);
                }
            }
        })
        .state('tables.ahathallan_modbus_scheme_address', {
            url: "/ahathallan_modbus_scheme_address",
            templateUrl: "views/ahathallan_modbus_scheme_address.html",
            data: { pageTitle: 'Modbus方案信息地址信息' },
            resolve: {
                loadPlugin: function ($ocLazyLoad) {
                    return $ocLazyLoad.load([
                        {
                            files: ['js/plugins/pace/pace.min.js']
                        },
                        {
                            files: ['css/plugins/touchspin/jquery.bootstrap-touchspin.min.css', 'js/plugins/touchspin/jquery.bootstrap-touchspin.min.js']
                        }
                    ]);
                }
            }
        })
        .state('tables.ahathallan_modbus_device', {
            url: "/ahathallan_modbus_device",
            templateUrl: "views/ahathallan_modbus_device.html",
            data: { pageTitle: 'Modbus设备数据' },
            resolve: {
                loadPlugin: function ($ocLazyLoad) {
                    return $ocLazyLoad.load([
                        {
                            files: ['css/plugins/iCheck/custom.css','js/plugins/iCheck/icheck.min.js','js/plugins/pace/pace.min.js','css/plugins/ionRangeSlider/ion.rangeSlider.css','css/plugins/ionRangeSlider/ion.rangeSlider.skinFlat.css','js/plugins/ionRangeSlider/ion.rangeSlider.min.js']
                        },
                        {
                            files: ['css/plugins/touchspin/jquery.bootstrap-touchspin.min.css', 'js/plugins/touchspin/jquery.bootstrap-touchspin.min.js']
                        },
                        {
                            files: ['js/plugins/sweetalert/sweetalert.min.js', 'css/plugins/sweetalert/sweetalert.css']
                        },
                        {
                            name: 'oitozero.ngSweetAlert',
                            files: ['js/plugins/sweetalert/angular-sweetalert.min.js']
                        }
                    ]);
                }
            }
        })
        .state('dashboards.ahathallan_modbus_device_grid', {
            url: "/ahathallan_modbus_device_grid",
            templateUrl: "views/ahathallan_modbus_device_grid.html",
            data: { pageTitle: 'Modbus设备列表' },
            resolve: {
                loadPlugin: function ($ocLazyLoad) {
                    return $ocLazyLoad.load([
                        {
                            files: ['js/plugins/pace/pace.min.js']
                        }
                    ]);
                }
            }
        })
        .state('dashboards.inspinia_device_info', {
            url: "/ahathallan_modbus_device_info",
            templateUrl: "views/ahathallan_modbus_device_info.html",
            data: { pageTitle: 'Modbus设备信息' },
            resolve: {
                loadPlugin: function ($ocLazyLoad) {
                    return $ocLazyLoad.load([
                        {
                            files: ['css/plugins/iCheck/custom.css','js/plugins/iCheck/icheck.min.js','js/plugins/pace/pace.min.js','css/plugins/ionRangeSlider/ion.rangeSlider.css','css/plugins/ionRangeSlider/ion.rangeSlider.skinFlat.css','js/plugins/ionRangeSlider/ion.rangeSlider.min.js']
                        }
                    ]);
                }
            }
        })
        .state('dashboards.ahathallan_watch_dash', {
            url: "/ahathallan_watch_dash",
            templateUrl: "views/ahathallan_watch_dash.html",
            data: { pageTitle: 'Modbus设备看板' },
            resolve: {
                loadPlugin: function ($ocLazyLoad) {
                    return $ocLazyLoad.load([
                        {
                            serie: true,
                            name: 'angular-flot',
                            files: ['js/plugins/flot/jquery.flot.js', 'js/plugins/flot/jquery.flot.time.js', 'js/plugins/flot/jquery.flot.tooltip.min.js', 'js/plugins/flot/jquery.flot.spline.js', 'js/plugins/flot/jquery.flot.resize.js', 'js/plugins/flot/jquery.flot.pie.js', 'js/plugins/flot/curvedLines.js', 'js/plugins/flot/angular-flot.js']
                        },
                        {
                            serie: true,
                            files: ['js/plugins/jvectormap/jquery-jvectormap-2.0.2.min.js', 'js/plugins/jvectormap/jquery-jvectormap-2.0.2.css']
                        },
                        {
                            serie: true,
                            files: ['js/plugins/jvectormap/jquery-jvectormap-world-mill-en.js']
                        },
                        {
                            name: 'ui.checkbox',
                            files: ['js/bootstrap/angular-bootstrap-checkbox.js']
                        }
                    ]);
                }
            }
        })
}
angular
    .module('inspinia')
    .config(config)
    .run(function($rootScope, $state) {
        $rootScope.$state = $state;
    });
