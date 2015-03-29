(function(){
	'use strict';

	var puppyControllers = angular.module('puppyControllers', []);

	puppyControllers.controller('PuppyCtrl', ['$scope', 'Phone',
		function($scope, Phone) {
			$scope.phones = Phone.query();
			$scope.orderProp = 'age';
		}]);
	
	puppyControllers.controller('PuppyTopCtrl', ['$scope', '$routeParams', 'Phone',
		function($scope, $routeParams, Phone) {
			$scope.phone = Phone.get({phoneId: $routeParams.phoneId}, function(phone) {
				$scope.mainImageUrl = phone.images[0];
			});
			$scope.setImage = function(imageUrl) {
				$scope.mainImageUrl = imageUrl;
			}
		}]);

})();