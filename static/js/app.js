(function(){
	'use strict';
	/* App Module */
	var pupsApp = angular.module('pupsApp', [
		'ngDialog',
		'ngRoute',
		'puppyControllers',
		'puppyServices'
		]);
	pupsApp.config(['$routeProvider',
		function($routeProvider) {
			$routeProvider.
			when('/', {
				templateUrl: 'partials/puppies.html',
				controller: 'PuppyCtrl'
			}).
			when('/top', {
				templateUrl: 'partials/top-puppies.html',
				controller: 'PuppyTopCtrl'
			}).
			otherwise({
				redirectTo: '/'
			});
		}]);
})();