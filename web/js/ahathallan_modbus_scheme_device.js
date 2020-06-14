/**
 * agileBoard - Controller for agile Board view$location.path("/2")
 */
function ahathallan_modbus_scheme_device($scope,$interval, $http,$location,SweetAlert) {
    //串口号
    $scope.Serial_value="COM1";

    $.getJSON("js/modbus_set.json", function (data) {
        $scope.Serial_list=data.Serial_list
    })

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

    $scope.spin_int = {
        min: 0,
        max: 10000000000,
        step: 1,
        decimals: 0,
        boostat: 5,
        maxboostedstep: 10,
    };

    $scope.spin_float = {
        min: 0,
        max: 10000000000,
        step: 0.01,
        decimals: 2,
        boostat: 5,
        maxboostedstep: 10,
    };

    $scope.scheme_id = $location.search().id

    $scope.scheme_address_list = []

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
                        $scope.scheme_address_list_value=[]
                        $scope.scheme_address_list_set=[]

                        result.forEach((item, index, arr) => {
                            if (item["Address_Name"] != '' && item["Address_Name"] != undefined) {
                                is_int=false
                                is_float=false
                                is_option=false

                                var option_arr=[]

                                if(item["Modbus_Data_type"]=="0" ||item["Modbus_Data_type"]=="1"){
                                    is_int=true
                                }else if(item["Modbus_Data_type"]=="2" ||item["Modbus_Data_type"]=="3"||item["Modbus_Data_type"]=="4"){
                                    is_float=true
                                }else {
                                    is_option = true
                                    arr = item["Option_context"].split(',');

                                    option_arr=[]

                                    for (var i=0;i<arr.length;i++)
                                    {
                                        option_arr.push({arr_index:i+ '',arr_value:arr[i]});
                                    }
                                }

                                if(item["Function_value"]=="03H") {
                                    $scope.scheme_address_list_value.push({
                                        Address_scheme_code: item["Address_scheme_code"],
                                        Address_No: item["Address_No"],
                                        Address_Name: item["Address_Name"],
                                        Address: item["Address"],
                                        Function_value: item["Function_value"],
                                        Function_Name: item["Function_Name"],
                                        Modbus_Data_type: item["Modbus_Data_type"],
                                        Modbus_Data_type_Name: item["Modbus_Data_type_Name"],
                                        Device_address_url: "web#/tables/ahathallan_modbus_scheme_address?Scheme_code=" + item["Address_scheme_code"] + "&Address_No=" + item["Address_No"],
                                        UOM: item["UOM"],
                                        Value: 0
                                    });
                                }

                                if(item["Function_value"]=="10H") {
                                    $scope.scheme_address_list_set.push({
                                        Address_scheme_code: item["Address_scheme_code"],
                                        Address_No: item["Address_No"],
                                        Address_Name: item["Address_Name"],
                                        Address: item["Address"],
                                        Function_value: item["Function_value"],
                                        Function_Name: item["Function_Name"],
                                        Modbus_Data_type: item["Modbus_Data_type"],
                                        Modbus_Data_type_Name: item["Modbus_Data_type_Name"],
                                        Device_address_url: "web#/tables/ahathallan_modbus_scheme_address?Scheme_code=" + item["Address_scheme_code"] + "&Address_No=" + item["Address_No"],
                                        UOM: item["UOM"],
                                        Value: 0,
                                        Is_int: is_int,
                                        Is_float: is_float,
                                        Is_option: is_option,
                                        Option_arr: option_arr
                                    });
                                }
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

    //========================刷新Modbus数据=========================
    function fresh_modbus_data(){
        var equ_Info = {
            Serial: $scope.Serial_value,
            SlaveId: $scope.SlaveId,
            BaudRate: $scope.BaudRate_value,
            Parity: $scope.Parity_value,
            DataBits: $scope.DataBits_value,
            StopBits: $scope.StopBits_value,
            Scheme_id: $scope.scheme_id
        };

        //获取数据类型
        var url = "/modbus_address_and_value_lst";

        //获取Modbus数据
        $http.post(url, equ_Info)
            .success(function (result) {
                if (result.length > 0) {
                    $scope.scheme_address_list_value=[]
                    $scope.scheme_address_list_set=[]

                    result.forEach((item, index, arr) => {
                        if (item["Address_Name"] != '' && item["Address_Name"] != undefined) {

                            is_int=false
                            is_float=false
                            is_option=false

                            var option_arr=[]

                            if(item["Modbus_Data_type"]=="0" ||item["Modbus_Data_type"]=="1"){
                                is_int=true
                            }else if(item["Modbus_Data_type"]=="2" ||item["Modbus_Data_type"]=="3"||item["Modbus_Data_type"]=="4"){
                                is_float=true
                            }else {
                                is_option = true
                                arr = item["Option_context"].split(',');

                                option_arr=[]

                                for (var i=0;i<arr.length;i++)
                                {
                                    option_arr.push({arr_index:i+ '',arr_value:arr[i]});
                                }
                            }

                            if(item["Function_value"]=="03H") {
                                $scope.scheme_address_list_value.push({
                                    Address_scheme_code: item["Address_scheme_code"],
                                    Address_No: item["Address_No"],
                                    Address_Name: item["Address_Name"],
                                    Address: item["Address"],
                                    Function_value: item["Function_value"],
                                    Function_Name: item["Function_Name"],
                                    Modbus_Data_type: item["Modbus_Data_type"],
                                    Modbus_Data_type_Name: item["Modbus_Data_type_Name"],
                                    Device_address_url: "web#/tables/ahathallan_modbus_scheme_address?Scheme_code=" + item["Address_scheme_code"] + "&Address_No=" + item["Address_No"],
                                    UOM: item["UOM"],
                                    Value: item["Modbus_value"]
                                });
                            }
                        }

                        if(item["Function_value"]=="10H") {
                            $scope.scheme_address_list_set.push({
                                Address_scheme_code: item["Address_scheme_code"],
                                Address_No: item["Address_No"],
                                Address_Name: item["Address_Name"],
                                Address: item["Address"],
                                Function_value: item["Function_value"],
                                Function_Name: item["Function_Name"],
                                Modbus_Data_type: item["Modbus_Data_type"],
                                Modbus_Data_type_Name: item["Modbus_Data_type_Name"],
                                Device_address_url: "web#/tables/ahathallan_modbus_scheme_address?Scheme_code=" + item["Address_scheme_code"] + "&Address_No=" + item["Address_No"],
                                UOM: item["UOM"],
                                Value: item["Modbus_value"],
                                Is_int: is_int,
                                Is_float: is_float,
                                Is_option: is_option,
                                Option_arr: option_arr
                            });
                        }
                    });
                }
            })
            .error(function () {

            });
    }

    $scope.is_link=false

    $scope.is_link_show=true

    $scope.is_dis_link_show=false

    //连接Modbus值
    $scope.modbus_link = function () {
        $scope.is_link=true
        $scope.is_link_show=false
        $scope.is_dis_link_show=true

        fresh_modbus_data()

        $scope.timer = $interval(function () {
            if($scope.is_link) {
                fresh_modbus_data()
            }
        },15000);
    }

    $scope.modbus_dis_link = function () {
        $scope.is_link=false
        $scope.is_link_show=true
        $scope.is_dis_link_show=false
    }

    $scope.Device_Name = ""
    //保存
    $scope.modbus_scheme_address_and_value_set = function () {
        var modbus_scheme = {
            // dev_no: $scope.equ_num,
            Scheme_code:$scope.scheme_id,
            Serial: $scope.Serial_value,
            Device_Name: $scope.Device_Name,
            Device_info: $scope.Device_info,
            Device_type: $scope.Device_type,
            SlaveId: $scope.SlaveId,
            BaudRate: $scope.BaudRate_value,
            DataBits: $scope.DataBits_value,
            Parity: $scope.Parity_value,
            StopBits: $scope.StopBits_value,
            modbus_address_and_value:$scope.scheme_address_list_set,
            Scheme_id: $scope.scheme_id
        };

        //上传路径
        var url = "/modbus_scheme_address_and_value_set";

        //保存数据
        $http.post(url, modbus_scheme)
            .success(function (result) {
                // $scope.id = result;
                // load_all();
                SweetAlert.swal({
                    title: "保存成功！",
                    text: "当前数据已保存成功！"
                });
            })
            .error(function () {

            });
    }
}

angular
    .module('inspinia')
    .controller('ahathallan_modbus_scheme_device', ahathallan_modbus_scheme_device)