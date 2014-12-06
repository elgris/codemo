angular.module("myapp", []).controller("CodeFormController", function($scope, $http) {
    $scope.codeForm = {};
    $scope.codeForm.src = "input := []int{1, 4, 2, 6, 4, 2, 4, 6}\n" + "for i := range input {\n" + "	for j := i + 1; j < len(input); j++ {\n" + "		if input[i] > input[j] {\n" + "			t := input[j]\n" + "			input[j] = input[i]\n" + "			input[i] = t\n" + "		}\n" + "	}\n" + "}\n" + "fmt.Printf(\"%v\", input)\n";

    $scope.submitCodeForm = function(item, event) {
        console.log("--> Submitting form");
        // TODO: show loader
        $http({
                method: 'POST',
                url: 'src',
                data: $.param($scope.codeForm), // pass in data as strings
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                } // set the headers so angular passing info as form data (not request payload)
            })
            .success(function(data) {
                console.log(data);
                $scope.successMessage = data
                // TODO: hide loader
            })
            .error(function(data, status, headers, config) {
                $scope.errorMessage = "HTTP error code: " + status
            });;
        return false;
    };

});
