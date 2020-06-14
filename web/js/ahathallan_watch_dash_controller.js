/**
 * agileBoard - Controller for agile Board view$location.path("/2")
 */
function ahathallan_watch_dash($scope,$interval, $http,$location) {
    dev_No=$location.search().id

    $scope.modbus_device_json = {}

    loadInfo = function (dev_No) {
        //加载事件
        $http.get("/modbus_Device_data/"+dev_No).success(function (result) {
            $scope.modbus_device_json=result
        }).error(function (data) {

        });
    }

    loadInfo(dev_No);

    $scope.timer = $interval(function () {
        loadInfo(dev_No);
    },5000);

    $scope.$on('$destroy',function(){
        $interval.cancel($scope.timer);
    })

}

angular
    .module('inspinia')
    .controller('ahathallan_watch_dash', ahathallan_watch_dash)