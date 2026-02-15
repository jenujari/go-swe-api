package lib

import (
	"fmt"
	"testing"
	"time"

	baselib "github.com/jenujari/planets-lib"
	"github.com/stretchr/testify/assert"
)

func Test_UTCToSiderealTime(T *testing.T) {
	expected := 2460716.931356 // Replace with the expected value for the given UTC time

	t1 := time.Date(2025, 2, 10, 10, 20, 0, 0, time.UTC)

	result, err := UTCToSiderealTime(t1)

	assert.NoError(T, err, "Expected no error, got %v", err)
	assert.InDelta(T, expected, result, 0.0001, "Expected %f, got %f", expected, result)
}

func Test_SiderealTimeToUTC(T *testing.T) {
	expected := time.Date(2025, 2, 10, 10, 20, 0, 0, time.UTC)

	result, err := SiderealTimeToUTC(2460716.931356)

	assert.NoError(T, err, "Expected no error, got %v", err)
	assert.Equal(T, expected, result, "Expected %v, got %v", expected, result)
}

func Test_LongDiff(T *testing.T) {
	t1 := time.Date(2026, 1, 21, 13, 0, 0, 0, time.UTC)
	jd, _ := UTCToSiderealTime(t1)

	diff, err := LongDiff(jd, baselib.PLANET_LIB_MAP["Sun"], baselib.PLANET_LIB_MAP["Mercury"])
	assert.NoError(T, err, "Expected no error, got %v", err)
	assert.InDelta(T, 0.0770, diff, 0.0001, "Expected %f, got %f", 0.0770, diff)
}

func Test_FindConjunctionRange(T *testing.T) {
	t1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedStartT := time.Date(2026, 1, 20, 4, 0, 0, 0, time.UTC)
	expectedEndT := time.Date(2026, 1, 23, 3, 0, 0, 0, time.UTC)

	startT, endT, inConj := FindConjunctionRange(t1, t1.Add(time.Duration(24*time.Hour*30)), 1, 1.0/24.0, baselib.PLANET_LIB_MAP["Sun"], baselib.PLANET_LIB_MAP["Mercury"])
	fmt.Println("Start time : ", startT)
	fmt.Println("End time : ", endT)

	assert.True(T, inConj, "Expected conjunction to be found")
	assert.Equal(T, expectedStartT, startT, "Expected start time %v, got %v", expectedStartT, startT)
	assert.Equal(T, expectedEndT, endT, "Expected end time %v, got %v", expectedEndT, endT)
}

func Test_GetPlanetCalculation(T *testing.T) {
	defer SweClear()

	// https://www.drikpanchang.com/planet/position/planetary-positions-sidereal.html?date=14/01/2026&time=13:45:30
	t1 := time.Date(2026, 1, 14, 13, 45, 30, 0, time.UTC)

	table := map[string]struct {
		expected baselib.PlanetCord
	}{
		"Moon": {
			expected: baselib.PlanetCord{Longitude: 222.80, Latitude: -5.11},
		},
		"Sun": {
			expected: baselib.PlanetCord{Longitude: 270.17, Latitude: 0.00},
		},
		"Mercury": {
			expected: baselib.PlanetCord{Longitude: 265.71, Latitude: -1.76},
		},
		"Mars": {
			expected: baselib.PlanetCord{Longitude: 268.92, Latitude: -0.97},
		},
		"Venus": {
			expected: baselib.PlanetCord{Longitude: 272.05, Latitude: -0.97},
		},
		"Jupiter": {
			expected: baselib.PlanetCord{Longitude: 85.31, Latitude: 0.27},
		},
		"Saturn": {
			expected: baselib.PlanetCord{Longitude: 332.87, Latitude: -2.22},
		},
		"Rahu": {
			expected: baselib.PlanetCord{Longitude: 317.22, Latitude: 0.00},
		},
		"Ketu": {
			expected: baselib.PlanetCord{Longitude: 137.22, Latitude: 0.00},
		},
	}

	siderealTime, err := UTCToSiderealTime(t1)
	if err != nil {
		T.Fatalf("Error converting UTC to sidereal time: %v", err)
	}

	for name, tc := range table {
		result, err := GetPlanetCalculation(siderealTime, name)

		assert.NoError(T, err, "%s: Expected no error, got %v", name, err)
		assert.NotNil(T, result, "%s: Expected non-nil result, got nil", name)

		assert.InDelta(T, tc.expected.Longitude, result.Longitude, 0.01, "%s: Expected Longitude %f, got %f", name, tc.expected.Longitude, result.Longitude)
		assert.InDelta(T, tc.expected.Latitude, result.Latitude, 0.01, "%s: Expected Latitude %f, got %f", name, tc.expected.Latitude, result.Latitude)
	}
}
