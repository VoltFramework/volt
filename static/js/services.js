var voltServices = angular.module('voltServices', ['ngResource']);

voltServices.factory('Tasks', function ($resource) {
    return $resource('/tasks', {}, {
                        'query': {
                            method: 'GET', 
                            isArray: true,
                            transformResponse: function (data, headers) {
                                var tasks = JSON.parse(data).tasks;
				for (id in tasks) {
				    switch (tasks[id].state) {
				    case 0:
					tasks[id].state="STARTING";
					break;
				    case 1:
					tasks[id].state="RUNNING";
					break;
				    case 2:
					tasks[id].state="FINISHED";
					break;
				    case 3:
					tasks[id].state="FAILED";
					break;
				    case 4:
					tasks[id].state="KILLED";
					break;
				    case 5:
					tasks[id].state="LOST";
					break;
				    case 6:
					tasks[id].state="STAGING";
					break;
				    }
				}
				return tasks;
                            }
                        }
    })
});