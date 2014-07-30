var voltControllers = angular.module('voltControllers', []);

voltControllers.controller('Tasks', ['$scope', 'Tasks', '$interval', function ($scope, Tasks, $interval) {
  $scope.refreshInterval = 5;
  $interval(function() {
      var update = Tasks.query(function() {
          $scope.tasks = update;
      });
  }, $scope.refreshInterval * 1000);
  $scope.tasks = Tasks.query();
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