var voltServices = angular.module('voltServices', ['ngResource']);

voltServices.factory('Tasks', function ($resource) {
    return $resource('/tasks', {}, {
                        'query': {
                            method: 'GET', 
                            isArray: true,
                            transformResponse: function (data, headers) {
				if (data == '') {
				    return data;
				}
                                var tasks = JSON.parse(data).tasks;
				for (id in tasks) {
				    switch (tasks[id].state) {
				    case 0:
				    case "0":
					tasks[id].state="STARTING";
					tasks[id].class="info";
					break;
				    case 1:
				    case "1":
					tasks[id].state="RUNNING";
					tasks[id].class="info";
					break;
				    case 2:
				    case "2":
					tasks[id].state="FINISHED";
					tasks[id].class="success";
					break;
				    case 3:
				    case "3":
					tasks[id].state="FAILED";
					tasks[id].class="danger";
					break;
				    case 4:
				    case "4":
					tasks[id].state="KILLED";
					tasks[id].class="danger";
					break;
				    case 5:
				    case "5":
					tasks[id].state="LOST";
					tasks[id].class="danger";
					break;
				    case 6:
				    case "6":
					tasks[id].state="STAGING";
					tasks[id].class="default";
					break;
				    }
				}
				return tasks;
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
