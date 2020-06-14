/**
 * agileBoard - Controller for agile Board view$location.path("/2")
 */
function ahathallan_modbus_scheme_info($scope, $http,$location, SweetAlert) {

    //站号
    $scope.SlaveId = "1";

    $scope.spin_Slave = {
        min: 0,
        max: 254,
        step: 1,
        decimals: 0,
        boostat: 5,
        maxboostedstep: 10,
    };

    //波特率
    $scope.BaudRate_value="9600";

    $scope.BaudRate_list=[
        {BaudRate_no:"1200",BaudRate_Value:"1200"},
        {BaudRate_no:"2400",BaudRate_Value:"2400"},
        {BaudRate_no:"4800",BaudRate_Value:"4800"},
        {BaudRate_no:"9600",BaudRate_Value:"9600"},
        {BaudRate_no:"19200",BaudRate_Value:"19200"},
        {BaudRate_no:"38400",BaudRate_Value:"38400"},
        {BaudRate_no:"115200",BaudRate_Value:"115200"}
    ]

    //校验位
    $scope.Parity_value="0";

    $scope.Parity_list=[
        {Parity_no:"0",Parity_Value:"无"},
        {Parity_no:"1",Parity_Value:"奇数"},
        {Parity_no:"2",Parity_Value:"偶数"}
    ]

    //数据位
    $scope.DataBits_value="0";

    $scope.DataBits_list=[
        {DataBits_no:"0",DataBits_value:"8"}
    ]

    //停止位
    $scope.StopBits_value="0";

    $scope.StopBits_list=[
        {StopBits_no:"0",StopBits_value:"1"},
        {StopBits_no:"0",StopBits_value:"2"}
    ]

    $scope.scheme_id = $location.search().id

    loadInfo = function () {
        //加载事件
        $http.get("/modbus_scheme_get/" + $scope.scheme_id).success(function (result) {
            if (result["Device_Name"] != "") {

                $scope.Device_url= "web#/tables/ahathallan_modbus_scheme_address?Scheme_code=" + $scope.scheme_id,
                $scope.Device_Name = result["Device_Name"];
                $scope.Device_info = result["Device_info"];
                $scope.Device_type = result["Device_type"];
                $scope.SlaveId = result["SlaveId"];
                $scope.BaudRate_value = result["BaudRate"];
                $scope.DataBits_value = result["DataBits"];
                $scope.Parity_value = result["Parity"];
                $scope.StopBits_value = result["StopBits"];

                $http.get("/modbus_address_get_lst/" + $scope.scheme_id, {
                    params: {}
                }).success(function (result) {
                    if (result.length > 0) {
                        $scope.scheme_address_list = []
                        result.forEach((item, index, arr) => {
                            if (item["Address_Name"] != '' && item["Address_Name"] != undefined) {
                                $scope.scheme_address_list.push({
                                    Address_scheme_code: item["Address_scheme_code"],
                                    Address_No: item["Address_No"],
                                    Address_Name: item["Address_Name"],
                                    Address: item["Address"],
                                    Function_value: item["Function_value"],
                                    Function_Name: item["Function_Name"],
                                    Modbus_Data_type: item["Modbus_Data_type"],
                                    Modbus_Data_type_Name: item["Modbus_Data_type_Name"],
                                    Device_address_url: "web#/tables/ahathallan_modbus_scheme_address?Scheme_code=" + item["Address_scheme_code"]+"&Address_No=" + item["Address_No"]
                                });
                            }
                        });
                    }
                }).error(function (data) {

                });
            }
        }).error(function (data) {

        });
    }

    loadInfo();

    $scope.delete_modbus_sheme_info = function () {
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
                    $http.get("/modbus_scheme_delete/" + $scope.scheme_id).success(function (result) {
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
                                $location.path("tables/ahathallan_modbus_scheme_lst");
                            });
                    })
                } else {

                }
            });
    }

    $scope.Device_Name = ""
    //保存
    $scope.save_modbus_sheme_info = function () {
        var modbus_scheme = {
            // dev_no: $scope.equ_num,
            Scheme_code:$scope.scheme_id,
            Device_Name: $scope.Device_Name,
            Device_info: $scope.Device_info,
            Device_type: $scope.Device_type,
            SlaveId: $scope.SlaveId,
            BaudRate: $scope.BaudRate_value,
            DataBits: $scope.DataBits_value,
            Parity: $scope.Parity_value,
            StopBits: $scope.StopBits_value
        };

        //上传路径
        var url = "/modbus_scheme_save";

        //保存数据
        $http.post(url, modbus_scheme)
            .success(function (result) {
                // $scope.id = result;
                // load_all();
                SweetAlert.swal({
                    title: "保存成功！",
                    text: "当前方案已保存成功！"
                });
            })
            .error(function () {

            });
    }

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
                                loadInfo();
                            });
                    })
                } else {

                }
            });
    }
}

angular
    .module('inspinia')
    .controller('ahathallan_modbus_scheme_info', ahathallan_modbus_scheme_info)