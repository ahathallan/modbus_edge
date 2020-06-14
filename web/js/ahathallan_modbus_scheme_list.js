/**
 * agileBoard - Controller for agile Board view$location.path("/2")
 */
function ahathallan_modbus_scheme_list($scope, $http,$location,SweetAlert) {
    //加载事件
    load_all = function () {
        $http.get("/modbus_scheme_get_lst", {
            params: {}
        }).success(function (result) {
            if (result.length > 0) {
                $scope.scheme_list = []
                result.forEach((item, index, arr) => {
                    if(item["Device_Name"]!='' && item["Device_Name"]!=undefined){
                        $http.get("/modbus_address_get_lst/" + item["Scheme_code"], {
                            params: {}
                        }).success(function (result_address) {
                            scheme_address_list = []
                            if (result_address.length > 0) {
                                result_address.forEach((item_address, index, arr) => {
                                    if (item_address["Address_Name"] != '' && item_address["Address_Name"] != undefined) {
                                        scheme_address_list.push({
                                            Address_scheme_code: item_address["Address_scheme_code"],
                                            Address_No: item_address["Address_No"],
                                            Address_Name: item_address["Address_Name"],
                                            Address: item_address["Address"],
                                            Function_value: item_address["Function_value"],
                                            Function_Name: item_address["Function_Name"],
                                            Modbus_Data_type: item_address["Modbus_Data_type"],
                                            Modbus_Data_type_Name: item_address["Modbus_Data_type_Name"],
                                            Device_address_url: "web#/tables/ahathallan_modbus_scheme_address?Scheme_code=" + item_address["Address_scheme_code"]+"&Address_No=" + item_address["Address_No"]
                                        });
                                    }
                                });
                            }
                            Img_url=""
                            if(item["Device_type"]=="半导体"){
                                Img_url="./img/181_3.jpg"
                            }else if(item["Device_type"]=="地理信息"){
                                Img_url="./img/181_5.jpg"
                            }else{
                                Img_url="./img/181_1.jpg"
                            }
                            $scope.scheme_list.push({
                                Scheme_code: item["Scheme_code"],
                                Device_Name: item["Device_Name"],
                                Img_url:Img_url,
                                Device_info: item["Device_info"],
                                Device_type: item["Device_type"],
                                Scheme_address_list: scheme_address_list,
                                Device_url: "web#/tables/ahathallan_modbus_scheme_info?id=" + item["Scheme_code"],
                                Device_info_url:"web#/tables/ahathallan_modbus_device?id=" + item["Scheme_code"],
                                Device_address_url:"web#/tables/ahathallan_modbus_scheme_address?Scheme_code=" + item["Scheme_code"]
                            });
                        }).error(function (data) {

                        });
                    }
                });
            }
        }).error(function (data) {

        });
    }

    $scope.navitage_to_info = function (id) {
        $location.path("tables/ahathallan_modbus_scheme_info").search({id: id});
    }

    load_all();

    $scope.delete_modbus_address_info = function (address_id) {
        SweetAlert.swal({
                title: "是否确认删除?",
                text: "您将删除这条信息!",
                type: "warning",
                showCancelButton: true,
                confirmButtonColor: "#DD6B55",
                confirmButtonText: "是的,确认删除!",
                cancelButtonText: "取消",
                closeOnConfirm: false,
                closeOnCancel: true },
            function (isConfirm) {
                if (isConfirm) {
                    $http.get("/modbus_address_delete/" + address_id).success(function (result) {
                        params: {}
                    }).success(function (result) {
                        SweetAlert.swal({
                                title: "删除成功！",
                                text: "点击回到主列表",
                                type: "success",
                                confirmButtonText: "确认",
                                closeOnConfirm: true
                            },
                            function () {
                                load_all();
                            });
                    })
                } else {

                }
            });
    }

    $checkThree = true;
}

angular
    .module('inspinia')
    .controller('ahathallan_modbus_scheme_list', ahathallan_modbus_scheme_list)