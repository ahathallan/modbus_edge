/**
 * agileBoard - Controller for agile Board view$location.path("/2")
 */
function ahathallan_modbus($scope,$interval, $http,$location,$interval) {
    $scope.set_value = 55;

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
        {Function_no:"04H",Function_Value:"读输入寄存器"}
    ]

    //数据类型
    $scope.Data_type_value="3";

    $scope.Data_type_list=[
        {Data_type_no:"0",Data_type_Value:"Unsigned Integer"},
        {Data_type_no:"1",Data_type_Value:"Integer"},
        {Data_type_no:"2",Data_type_Value:"Double Presion"},
        {Data_type_no:"3",Data_type_Value:"IEEE Floating Point"},
        {Data_type_no:"4",Data_type_Value:"IEEE Reserved World"}
    ]

    //Modbus值
    $scope.Modbus_value = 0

    $scope.is_query=false

    function fresh_modbus_data(){
        var equ_Info = {
            Serial: $scope.Serial_value,
            SlaveId: $scope.SlaveId,
            BaudRate: $scope.BaudRate_value,
            Parity: $scope.Parity_value,
            Function_value: $scope.Function_value,
            Modbus_Data_type: $scope.Data_type_value,
            Address:$scope.Address
        };

        //上传路径
        var url = "/modbus_once_get";

        //保存数据
        $http.post(url, equ_Info)
            .success(function (result) {
                $scope.Modbus_value= result;
                // load_all();
            })
            .error(function () {

            });
    }

    //一次性获取modbus值
    $scope.modbus_once_get = function () {
        $scope.is_query=true
        fresh_modbus_data()
    }

    $scope.is_link=false

    $scope.is_link_show=true

    $scope.is_dis_link_show=false

    //连接
    $scope.modbus_link = function () {
        $scope.is_link=true
        $scope.is_link_show=false
        $scope.is_dis_link_show=true
    }

    //取消连接
    $scope.modbus_dis_link = function () {
        $scope.is_link=false
        $scope.is_query=false
        $scope.is_link_show=true
        $scope.is_dis_link_show=false
    }

    $scope.timer = $interval(function () {
        if($scope.is_link && $scope.is_query)
        {
            fresh_modbus_data()
        }
    },5000);


    // //连接
    // $scope.modbus_link = function () {
    //     $http.get("/modbus_link").success(function (result) {
    //         if (result) {
    //             refreshData();
    //         }
    //     })
    // }
    //
    // //取消连接
    // $scope.modbus_dis_link = function () {
    //     $http.get("/modbus_dis_link").success(function (result) {
    //         if (result) {
    //             refreshData();
    //         }
    //     })
    // }


}

function ahathallan_modbus_set($scope, $http,$location,$interval) {
    $scope.set_value = 55;

    //串口号
    $scope.Serial_value="COM1";

    $.getJSON("js/modbus_set.json", function (data) {
        // $scope.Serial_list=[
        //     {Serial_no:"COM1",Serial_Value:"COM1"},
        //     {Serial_no:"COM2",Serial_Value:"COM2"},
        //     {Serial_no:"COM3",Serial_Value:"COM3"},
        //     {Serial_no:"COM4",Serial_Value:"COM4"},
        //     {Serial_no:"COM5",Serial_Value:"COM5"},
        //     {Serial_no:"COM6",Serial_Value:"COM6"},
        //     {Serial_no:"COM7",Serial_Value:"COM7"},
        //     {Serial_no:"COM8",Serial_Value:"COM8"}
        // ]
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
    $scope.Function_value="10H";

    $scope.Function_list=[
        {Function_no:"05H",Function_Value:"写单个线圈"},
        {Function_no:"06H",Function_Value:"写单个保持寄存器"},
        {Function_no:"0EH",Function_Value:"写多个线圈"},
        {Function_no:"10H",Function_Value:"写多个保持寄存器"}
    ]

    //数据类型
    $scope.Data_type_value="0";

    $scope.Data_type_list=[
        {Data_type_no:"0",Data_type_Value:"Unsigned Integer"},
        {Data_type_no:"1",Data_type_Value:"Integer"},
        {Data_type_no:"2",Data_type_Value:"Double Presion"},
        {Data_type_no:"3",Data_type_Value:"IEEE Floating Point"},
        {Data_type_no:"4",Data_type_Value:"IEEE Reserved World"}
    ]

    //Modbus值
    $scope.Modbus_value = 0

    $scope.is_query=false

    function set_modbus_data(){
        var equ_Info = {
            Serial: $scope.Serial_value,
            SlaveId: $scope.SlaveId,
            BaudRate: $scope.BaudRate_value,
            Parity: $scope.Parity_value,
            Function_value: $scope.Function_value,
            Modbus_Data_type: $scope.Data_type_value,
            Address:$scope.Address,
            Modbus_value:$scope.Modbus_value
        };

        //上传路径
        var url = "/modbus_once_set";

        //保存数据
        $http.post(url, equ_Info)
            .success(function (result) {
                //$scope.Modbus_value= result;
                // load_all();
            })
            .error(function () {

            });
    }

    //一次性设置modbus值
    $scope.modbus_once_get = function () {
        set_modbus_data()
    }

    $scope.is_float=false

    $scope.is_int=true

    $scope.select_change = function () {
        $scope.Modbus_value = 0

        if($scope.Data_type_value=="0" || $scope.Data_type_value=="1"){
            $scope.is_float=false
            $scope.is_int=true
        }else{
            $scope.is_float=true
            $scope.is_int=false
        }
    }

    $scope.spin_int = {
        min: 0,
        max: 100000,
        step: 1,
        decimals: 0,
        boostat: 5,
        maxboostedstep: 10,
    };

    $scope.spin_float = {
        min: 0,
        max: 100000,
        step: 0.01,
        decimals: 2,
        boostat: 5,
        maxboostedstep: 10,
    };
}

angular
    .module('inspinia')
    .controller('ahathallan_modbus', ahathallan_modbus)
    .controller('ahathallan_modbus_set', ahathallan_modbus_set)