function eCharts($http, $window) {
    function link($scope, element, attrs) {
        var echarts_theme=$scope.$eval(attrs.eTheme);

        var myChart = echarts.init(element[0],echarts_theme);
        $scope.$watch(attrs['eData'], function() {
            var option = $scope.$eval(attrs.eData);
            if (angular.isObject(option)) {
                myChart.setOption(option);
            }
        }, true);
        $scope.getDom = function() {
            return {
                'height': element[0].offsetHeight,
                'width': element[0].offsetWidth
            };
        };
        $scope.$watch($scope.getDom, function() {
            // resize echarts图表
            myChart.resize();
        }, true);

        myChart.on('dataZoom', function (params) {
            var f_eZoom=$scope.$eval(attrs['eZoom'])
            if (f_eZoom) {
                var handler = f_eZoom(params);
                handler && handler(params)
            }

            // START_TIME=getLocalTime(params.batch[0].startValue)
            // END_TIME=getLocalTime(params.batch[0].endValue)
            //
            // console.log(START_TIME.toLocaleString())
            // console.log(END_TIME.toLocaleString())
        })

        myChart.on('click', function (param) {
            if ($scope.ngClick) {
                var handler = $scope.ngClick(param);
                handler && handler(param)
            }
        })
    }
    return {
        restrict: 'A',
        link: link,
    };
};

/**
*
* Pass all functions into module
*/
angular
    .module('inspinia')
    .directive('eChart', eCharts)