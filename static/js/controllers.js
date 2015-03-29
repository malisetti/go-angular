(function(){
	'use strict';

	var puppyControllers = angular.module('puppyControllers', ['puppyFactory', 'ngDialog']);
	
	puppyControllers.controller('PuppyCtrl', ['$scope', 'Puppy', 'ngDialog', function ($scope, Puppy, ngDialog) {
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

		var getPuppies = function(){
			Puppy.getPuppies($scope.main.page).
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

		var votePuppy = function(vote){
			Puppy.votePuppy(vote).
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

				votePuppy(vote);
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

				votePuppy(vote);
			}
		}


		getPuppies();

		$scope.nextPage = function() {
			if ($scope.main.page < $scope.main.pages) {
				$scope.main.page++;
				getPuppies();
			}
		};

		$scope.previousPage = function() {
			if ($scope.main.page > 1) {
				$scope.main.page--;
				getPuppies();
			}
		};

		$scope.clickToOpen = function () {
			ngDialog.open({ template: 'popupTmpl.html' });
		};

	}]);

puppyControllers.controller('TopPuppyCtrl', ['$scope', 'Puppy', 'ngDialog', function ($scope, Puppy, ngDialog){
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


	var getTopPuppies = function(){
			Puppy.getTopPuppies($scope.main.page).
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

	var votePuppy = function(vote){
			Puppy.votePuppy(vote).
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

				votePuppy(vote);
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

				votePuppy(vote);
			}
		}


		getTopPuppies();

		$scope.nextPage = function() {
			if ($scope.main.page < $scope.main.pages) {
				$scope.main.page++;
				getTopPuppies();
			}
		};

		$scope.previousPage = function() {
			if ($scope.main.page > 1) {
				$scope.main.page--;
				getTopPuppies();
			}
		};

		$scope.clickToOpen = function () {
			ngDialog.open({ template: 'popupTmpl.html' });
		};

}]);
})();