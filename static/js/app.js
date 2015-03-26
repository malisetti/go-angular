(function () {
    'use strict';

    var app = angular.module('pups', ['ngDialog']);

    app.controller('PuppyCtrl', ['$scope', '$http', 'ngDialog', function ($scope, $http, ngDialog) {
            $scope.main = {
              page: 1,
              pages: 1,
              loading: false
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
                $scope.main.page--;
              }).
              finally(function(){
                $scope.main.loading = false;
              });
            };

            $scope.open = function (pup) {
              var newScope = $scope.$new();
              newScope.pup = pup;

              ngDialog.open({
                template: '<img ng-src="{{ pup.large }}" alt="{{ pup.title }}"/>',
                plain: true,
                scope: newScope
              });
            };

            $scope.upvote = function(pup) {

            }

            $scope.downvote = function() {

            }


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

            $scope.clickToOpen = function () {
              ngDialog.open({ template: 'popupTmpl.html' });
            };

        }]);
})();
