var voltServices = angular.module('voltServices', ['ngResource']);

voltServices.factory('Tasks', function ($resource) {
    return $resource('/tasks', {}, {
                        'query': {
                            method: 'GET', 
                            isArray: false,
                            transformResponse: function (data, headers) {
				if (data == '') {
				    return data;
				}

				data = JSON.parse(data);
				for (id in data.tasks) {
				    switch (data.tasks[id].state) {
				    case 0:
					data.tasks[id].state="STARTING";
					data.tasks[id].class="info";
					break;
				    case 1:
					data.tasks[id].state="RUNNING";
					data.tasks[id].class="info";
					break;
				    case 2:
					data.tasks[id].state="FINISHED";
					data.tasks[id].class="success";
					break;
				    case 3:
					data.tasks[id].state="FAILED";
					data.tasks[id].class="danger";
					break;
				    case 4:
					data.tasks[id].state="KILLED";
					data.tasks[id].class="danger";
					break;
				    case 5:
					data.tasks[id].state="LOST";
					data.tasks[id].class="danger";
					break;
				    case 6:
					data.tasks[id].state="STAGING";
					data.tasks[id].class="default";
					break;
				    }
				}

				return data;
                            }
                        }
    })
});

voltServices.factory('Metrics', function ($resource) {
    return $resource('/metrics', {}, {
                        'query': {
                            method: 'GET', 
                            isArray: false
			}
    })
});