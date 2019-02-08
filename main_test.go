package main

import (
	"testing"
	"errors"
)

func TestGetUserIdFromToken(t *testing.T) {
	passFunc := func(string) (string, error) {
		return "foo", nil
	}
	failFunc := func(string) (string, error) {
		return "", errors.New("bar")
	}

	tests := []struct{
		name string
		function userGetter
		token string
		actual string
		errString string
	}{
		{
			"Sad: Missing token",
			passFunc,
			"",
			"",
			"Missing Authorization token",
		},
		{
			"Sad: Cannot get user",
			failFunc,
			"ff",
			"",
			"Invalid token",
		},
		{
			"Happy: Returns user",
			passFunc,
			"fff",
			"foo",
			"",
		},
	}

	for _, tc := range tests {
		actual, errString := getUserIdFromToken(tc.function, tc.token)

		if actual != tc.actual {
			t.Errorf("%v:\n  Expected:%v\n  Actual:%v", tc.name, tc.actual, actual)
		}

		if errString != tc.errString {
			t.Errorf("%v:\n  Expected:%v\n  Actual:%v", tc.name, tc.errString, errString)
		}
	}
}
