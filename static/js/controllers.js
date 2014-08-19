var voltControllers = angular.module('voltControllers', []);

voltControllers.controller('Tasks', ['$scope', 'Tasks', '$interval', '$http', function ($scope, Tasks, $interval, $http) {
  $scope.refreshInterval = 5;
  $interval(function() {
      Tasks.query(function(d) {
          $scope.tasks = d;
      });
  }, $scope.refreshInterval * 1000);
  $scope.tasks = Tasks.query();

    $scope.trash = function (id) {
      $http({method: 'DELETE', url: '/tasks/'+id}).success(function(data) {});
    };
    $scope.kill = function (id) {
      $http({method: 'PUT', url: '/tasks/'+id+'/kill'}).success(function(data) {});
    };
}]);


voltControllers.controller('Modal', function ($scope, $modal, $log) {
  $scope.task = {
    cpus:'0.5',
    mem:'512',
    cmd:'/bin/ls'
  }
  $scope.open = function (size) {

    var modalInstance = $modal.open({
      templateUrl: 'modal.html',
      controller:  ModalCtrl,
      size: size,
      resolve: {task: function() {return $scope.task;}
    }
    });
  };
});

var ModalCtrl = function ($scope, $modalInstance, $http, task) {
    $scope.task = task;
    
    $scope.send = function () {
        $http({method: 'POST', url: '/tasks', data : $scope.task, headers:{'Accept': 'application/json', 'Content-Type': 'application/json; ; charset=UTF-8'}}).success(function(data) {
    });
        $modalInstance.dismiss('send');
    };

    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
};

voltControllers.controller('File', function ($scope, $modal, $http) {
  $scope.file = {};
    $scope.refresh = function() {
	$http.get('/tasks/'+$scope.file.id+'/file/volt_'+$scope.file.name).
	    success(function(data, status, headers, config) {
		    $scope.file.content= data;
	    }).
	    error(function(data, status, headers, config) {
		    $scope.file.content= 'error';
	    });
    };
    $scope.open = function (name, id, size) {
	$scope.file.name = name;
	$scope.file.id = id;
	$scope.refresh();
    var modalInstance = $modal.open({
      templateUrl: 'file.html',
      controller:  FileCtrl,
      size: size,
      resolve: {file: function() {return $scope.file;}
    }
    });
  };
});

var FileCtrl = function ($scope, $modalInstance, $http,file) {
    $scope.file = file;
      $scope.refresh = function() {
	$http.get('/tasks/'+$scope.file.id+'/file/volt_'+$scope.file.name).
	    success(function(data, status, headers, config) {
		    $scope.file.content= data;
	    }).
	    error(function(data, status, headers, config) {
		    $scope.file.content= 'error';
	    });
    };
    $scope.close = function () {
        $modalInstance.dismiss('close');
    };
};
