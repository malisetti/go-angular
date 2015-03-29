(function(){
	'use strict';

	var puppyControllers = angular.module('puppyControllers', ['ngDialog']);
	puppyControllers.controller('PuppyCtrl', ['$scope', '$http', 'ngDialog', function ($scope, $http, ngDialog) {
    $scope.main = {
      page: 1,
      pages: 1,
      loading: false
    };

    $scope.upvotes = [];
    $scope.downvotes = [];

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

    $scope.votePuppy = function(vote){
      var url = "/pups";

      $http.put(url, vote).
      success(function (data, status, headers, config) {
      }).
      error(function (data, status, headers, config) {
                // called asynchronously if an error occurs
                // or server returns response with an error status.
                $scope.main.page--;
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
      var alreadyUpVoted = $scope.upvotes.indexOf(pup.id);
      var alreadyDownVoted = $scope.downvotes.indexOf(pup.id);
      if(alreadyUpVoted * alreadyDownVoted == 1){
        $scope.upvotes.push(pup.id);
        pup.upvotes++

        var vote = {};
        vote.ID = pup.id;
        vote.VT = true;

        $scope.votePuppy(vote);
      }
    }

    $scope.downvote = function(pup) {
      var alreadyUpVoted = $scope.upvotes.indexOf(pup.id);
      var alreadyDownVoted = $scope.downvotes.indexOf(pup.id);
      if(alreadyUpVoted * alreadyDownVoted == 1){
        $scope.downvotes.push(pup.id);
        pup.downvotes++
        var vote = {};
        vote.ID = pup.id;
        vote.VT = false;

        $scope.votePuppy(vote);
      }
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

  puppyControllers.controller('TopPuppyCtrl', ['$scope', function($scope){

  }]);
})();