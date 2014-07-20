var voltServices = angular.module('voltServices', ['ngResource']);

voltServices.factory('Tasks', function ($resource) {
    return $resource('http://198.27.68.58:8080/tasks', {}, {
                        'query': {
                            method: 'GET', 
                            isArray: true,
                            transformResponse: function (data, headers) {
                                return JSON.parse(data).tasks;
                            }
                        }
    })
});