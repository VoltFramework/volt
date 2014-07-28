var voltServices = angular.module('voltServices', ['ngResource']);

voltServices.factory('Tasks', function ($resource) {
    return $resource('/tasks', {}, {
                        'query': {
                            method: 'GET', 
                            isArray: true,
                            transformResponse: function (data, headers) {
                                return JSON.parse(data).tasks;
                            }
                        }
    })
});