(function(){

	'use strict';
	/* Services */
	var puppyServices = angular.module('puppyServices', ['ngResource']);
	phonecatServices.factory('Puppy', ['$resource',
		function($resource){
			return $resource('phones/:phoneId.json', {}, {
				query: {method:'GET', params:{phoneId:'phones'}, isArray:true}
			});
		}]);
})();