/**
 * agileBoard - Controller for agile Board view$location.path("/2")
 */
function ahathallan_modbus_scheme_address($scope, $http,$location) {
    $scope.scheme_id = $location.search().Scheme_code
    $scope.address_No = $location.search().Address_No

    $http.get("/modbus_scheme_get/" + $scope.scheme_id).success(function (result) {
        if (result["Device_Name"] != "") {

            $scope.Device_Name = result["Device_Name"];
            $scope.Device_info = result["Device_info"];
            $scope.Device_type = result["Device_type"];
        }
    }).error(function (data) {

    });

    //Modbus地址名称
    $scope.Address_No = ""

    //Modbus地址
    $scope.Address="1";

    $scope.spin_address = {
        min: 0,
        max: 100000,
        step: 1,
        decimals: 0,
        boostat: 5,
        maxboostedstep: 10,
    };

    //指令代码
    $scope.Function_value="03H";

    $scope.Function_list=[
        {Function_no:"01H",Function_Value:"读线圈状态"},
        {Function_no:"02H",Function_Value:"读离散输入状态"},
        {Function_no:"03H",Function_Value:"读保持寄存器"},
        {Function_no:"04H",Function_Value:"读输入寄存器"},
        {Function_no:"05H",Function_Value:"写单个线圈"},
        {Function_no:"06H",Function_Value:"写单个保持寄存器"},
        {Function_no:"0EH",Function_Value:"写多个线圈"},
        {Function_no:"10H",Function_Value:"写多个保持寄存器"}
    ]

    //数据类型
    $scope.Modbus_Data_type="3";

    $scope.Data_type_list=[
        {Data_type_no:"0",Data_type_Value:"Unsigned Integer"},
        {Data_type_no:"1",Data_type_Value:"Integer"},
        {Data_type_no:"2",Data_type_Value:"Double Presion"},
        {Data_type_no:"3",Data_type_Value:"IEEE Floating Point"},
        {Data_type_no:"4",Data_type_Value:"IEEE Reserved World"},
        {Data_type_no:"5",Data_type_Value:"多选框"}
    ]

    $http.get("/modbus_address_get/" + $scope.address_No).success(function (result) {
        if (result["Address_Name"] != "") {
            $scope.Address_Name = result["Address_Name"];
            $scope.Address = result["Address"];
            $scope.Address_Json = result["Address_Json"];
            $scope.Function_value = result["Function_value"];
            $scope.Data_type_value = result["Data_type_value"];
            $scope.Modbus_Data_type = result["Modbus_Data_type"];
            $scope.Option_context = result["Option_context"];
            $scope.UOM = result["UOM"];

            if($scope.Modbus_Data_type=="0" || $scope.Modbus_Data_type=="1"){
                $scope.is_float=false
                $scope.is_int=true
                $scope.is_option=false
            }else if($scope.Modbus_Data_type=="2" || $scope.Modbus_Data_type=="3"|| $scope.Modbus_Data_type=="4"){
                $scope.is_float=true
                $scope.is_int=false
                $scope.is_option=false
            }else {
                $scope.is_float=true
                $scope.is_int=false
                $scope.is_option=true
            }
        }
    }).error(function (data) {

    });

    $scope.is_float=true

    $scope.is_int=false

    $scope.is_option=false

    $scope.select_change = function () {
        if($scope.Modbus_Data_type=="0" || $scope.Modbus_Data_type=="1"){
            $scope.is_float=false
            $scope.is_int=true
            $scope.is_option=false
        }else if($scope.Modbus_Data_type=="2" || $scope.Modbus_Data_type=="3"|| $scope.Modbus_Data_type=="4"){
            $scope.is_float=true
            $scope.is_int=false
            $scope.is_option=false
        }else {
            $scope.is_float=true
            $scope.is_int=false
            $scope.is_option=true
        }
    }

    //保存
    $scope.save_modbus_address_info = function () {
        var modbus_scheme_address = {
            // dev_no: $scope.equ_num,
            Address_scheme_code: $scope.scheme_id,
            Address_No: $scope.address_No,
            Address_Name:$scope.Address_Name,
            Address_Json:$scope.Address_Json,
            Address:$scope.Address,
            Function_value: $scope.Function_value,
            Modbus_Data_type: $scope.Modbus_Data_type,
            Option_context: $scope.Option_context,
            UOM: $scope.UOM
        };

        //上传路径
        var url = "/modbus_address_save";

        //保存数据
        $http.post(url, modbus_scheme_address)
            .success(function (result) {
                $location.path("tables/ahathallan_modbus_scheme_lst");
            })
            .error(function () {

            });
    }

    //返回列表
    $scope.navitage_to_list = function () {
        $location.path("tables/ahathallan_modbus_scheme_lst");
    }
}

angular
    .module('inspinia')
    .controller('ahathallan_modbus_scheme_address', ahathallan_modbus_scheme_address)