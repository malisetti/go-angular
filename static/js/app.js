(function () {
    var app = angular.module('pups', ['simplePagination']);
    
    app.controller('PuppyCtrl', ['$scope', '$http', 'Pagination', function ($scope, $http, Pagination) {
            $scope.response = {};
            $scope.response.page = 0;
            $scope.response.pages = 0;
            $scope.response.perpage = 0;
            $scope.response.total = 0;
            $scope.response.images = [];
            $scope.pagination = Pagination.getNew(10);
            var url = $scope.pagination.page === undefined ?  "/pups" : "/pups/" +  $scope.pagination.page;
            
            $http.get(url).
                    success(function (data, status, headers, config) {
                        $scope.response = data;
                        $scope.pagination = Pagination.getNew(data.perpage);
                        $scope.pagination.numPages = data.pages;
                    }).
                    error(function (data, status, headers, config) {
                        // called asynchronously if an error occurs
                        // or server returns response with an error status.
                    });
        }]);
})();
