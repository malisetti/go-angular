(function () {
    'use strict';

    var app = angular.module('pups', []);

    app.controller('PuppyCtrl', ['$scope', '$http', function ($scope, $http) {
            $scope.main = {
              page: 1,
              pages: 1
            };

            $scope.response = {}
            ,$scope.response.page = 0
            ,$scope.response.pages = 0
            ,$scope.response.perpage = 0
            ,$scope.response.total = 0
            ,$scope.response.images = [];

            $scope.getPuppies = function(){
              var url = "/pups";
              url = ($scope.main.page === undefined || $scope.main.page === 0) ? url : url + "/" + $scope.main.page;

              $http.get(url).
              success(function (data, status, headers, config) {
                $scope.response = data;
                $scope.main.pages = data.pages;
              }).
              error(function (data, status, headers, config) {
                // called asynchronously if an error occurs
                // or server returns response with an error status.
              });
            };

            $scope.getPuppies();

            $scope.nextPage = function() {
              if ($scope.main.page < $scope.main.pages) {
                $scope.main.page++;
                $scope.getPuppies();
              }
            };

            $scope.previousPage = function() {
              if ($scope.main.page > 1) {
                $scope.main.page--;
                $scope.getPuppies();
              }
            };

        }]);
})();
