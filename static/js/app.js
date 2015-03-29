(function () {
  'use strict';

  var app = angular.module('pupsApp', ['ngRoute', 'puppyControllers']);

  app.config(['$routeProvider',
    function($routeProvider) {
      $routeProvider.
      when('/', {
        templateUrl: 'partials/puppies.html',
        controller: 'PuppyCtrl'
      }).
      when('/top', {
        templateUrl: 'partials/top-puppies.html',
        controller: 'TopPuppyCtrl'
      }).
      otherwise({
        redirectTo: '/'
      });
    }]);
})();
