var voltApp = angular.module('voltApp', [
  'voltControllers',
  'voltServices',
  'ui.bootstrap',
  'angularMoment'
]);

voltApp.config(['$httpProvider', function($httpProvider) {
        $httpProvider.defaults.useXDomain = true;
        delete $httpProvider.defaults.headers.common['X-Requested-With'];
        }
]);
