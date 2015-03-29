(function(){
	'use strict';
	/* Services */
	var puppyFactory = angular.module('puppyFactory', []);

	puppyFactory.factory('Puppy', ['$http', function($http){
		var puppyUrlBase = '/pups';
		var topPuppyUrlBase = '/top';
    	var factory = {};

    	factory.getPuppies = function(pageId){
    		var url = (pageId === undefined || pageId === 0) ? puppyUrlBase : puppyUrlBase + "/" + pageId;
	        return $http.get(url);
    	};

    	factory.votePuppy = function(vote){
  			return $http.put(puppyUrlBase, vote);
    	};

    	factory.getTopPuppies = function(pageId){
    		var url = (pageId === undefined || pageId === 0) ? topPuppyUrlBase + "/1" : topPuppyUrlBase + "/" + pageId;
    		return $http.get(url);
    	};

    	return factory;

	}]);
})();