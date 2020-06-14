/**
 * agileBoard - Controller for agile Board view$location.path("/2")
 */
function ahathallan_modbus_device_grid($scope, $http,$location) {
    //加载事件
    load_all=function(){
        $http.get("/modbus_device_get_lst",{
            params: {
            }
        }).success(function (result) {
            if (result.length > 0) {
                $scope.equ_lst=[]
                result.forEach((item, index, arr) => {
                    $scope.equ_lst.push({Dev_No:item["Dev_No"],Dev_Name:item["Dev_Name"],device_url:"web#/dashboards/ahathallan_modbus_device_info?id="+item["Dev_No"]});
                });
            }
        }).error(function (data) {

        });
    }

    $scope.navitage_to_info = function (id){
        $location.path("dashboards/ahathallan_watch_dash").search({id:id});
    }

    load_all();

    $checkThree=true;
}

function ahathallan_modbus_device_info($scope, $http,$location) {
    //加载事件
    $scope.scheme_list = []

    $scope.scheme_value= ""

    load_all_shceme = function () {
        $http.get("/modbus_scheme_get_lst", {
            params: {}
        }).success(function (result) {
            if (result.length > 0) {
                result.forEach((item, index, arr) => {
                    if(item["Device_Name"]!='' && item["Device_Name"]!=undefined){
                        $scope.scheme_list.push({
                            Scheme_code: item["Scheme_code"],
                            Device_Name: item["Device_Name"],
                            Serial: item["Serial"]
                        });
                    }
                });

                $scope.scheme_value=result[0]["Scheme_code"]
            }
        }).error(function (data) {

        });
    }

    //串口号
    $scope.Serial_value="COM1";

    $.getJSON("js/modbus_set.json", function (data) {
        $scope.Serial_list=data.Serial_list
    })

    load_all_shceme();

    dev_No=$location.search().id

    $scope.loadInfo = function (dev_No) {
        //加载事件
        $http.get("/modbus_device_get/"+dev_No).success(function (result) {
            if(result["dev_no"]!="") {

                $scope.Dev_No = result["Dev_No"];
                $scope.Dev_Name = result["Dev_Name"];
                $scope.scheme_value = result["Scheme_code"];
                $scope.Serial = result["Serial"];
            }else {

            }
        }).error(function (data) {

        });
    }

    $scope.loadInfo(dev_No);

    //保存
    $scope.submitInfo = function () {
        var equ_Info = {
            // dev_no: $scope.equ_num,
            Dev_No: $scope.Dev_No,
            Dev_Name: $scope.Dev_Name,
            Scheme_code: $scope.scheme_value,
            Serial: $scope.Serial
        };

        //上传路径
        var url = "/modbus_device_save";

        //保存数据
        $http.post(url, equ_Info)
                    .success(function (result) {
                        // $scope.id = result;
                        // load_all();
                    })
                    .error(function () {

                    });
    }
}

angular
    .module('inspinia')
    .controller('ahathallan_modbus_device_info', ahathallan_modbus_device_info)
    .controller('ahathallan_modbus_device_grid', ahathallan_modbus_device_grid)