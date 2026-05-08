package lib

import (
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
	t1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).UTC()
	t2 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(24 * time.Hour * 60)).UTC()
	expectedStartT := time.Date(2026, 1, 20, 4, 0, 0, 0, time.UTC)
	expectedEndT := time.Date(2026, 1, 23, 3, 0, 0, 0, time.UTC)

	startT, endT, inConj, err := FindConjunctionRange(t1, t2, 1, 1.0/24.0, baselib.PLANET_LIB_MAP["Sun"], baselib.PLANET_LIB_MAP["Mercury"])

	assert.NoError(T, err, "Expected no error, got %v", err)
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

func Test_CalcTithy(T *testing.T) {
	defer SweClear()

	w1 := "Wednesday"
	n1 := "Anuradha"
	w2 := "Friday"
	n2 := "Chitra"

	// Based on the same reference timestamp used in Test_GetPlanetCalculation:
	// Moon ~= 222.80, Sun ~= 270.17 => delta ~= 312.63 => floor(312.63/12)+1 = 27
	t1 := time.Date(2026, 1, 14, 13, 45, 30, 0, time.UTC)

	tithy, w, n, err := CalcTithy(t1)

	assert.NoError(T, err, "Expected no error, got %v", err)
	assert.Equal(T, int32(27), tithy, "Expected Tithy %d, got %d", 27, tithy)
	assert.Equal(T, w1, w, "Expected Weekday %s, got %s", w1, w)
	assert.Equal(T, n1, n, "Expected Nakshatra %s, got %s", n1, n)

	t2 := time.Date(2026, 3, 6, 6, 0, 0, 0, time.UTC)

	tithy, w, n, err = CalcTithy(t2)

	assert.NoError(T, err, "Expected no error, got %v", err)
	assert.Equal(T, int32(18), tithy, "Expected Tithy %d, got %d", 18, tithy)
	assert.Equal(T, w2, w, "Expected Weekday %s, got %s", w2, w)
	assert.Equal(T, n2, n, "Expected Nakshatra %s, got %s", n2, n)
}

func Test_GetAllPlanetsBalas(T *testing.T) {
	t1 := time.Date(2026, 1, 14, 13, 45, 30, 0, time.UTC)
	result, err := GetAllPlanetsBalas(t1)

	assert.NoError(T, err, "Expected no error, got %v", err)
	assert.NotNil(T, result, "Expected non-nil result, got nil")

	planetsMap, ok := result.(map[string]PlanetBalas)
	assert.True(T, ok, "Expected result to be map[string]PlanetBalas")

	// Expect all 9 planets (Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn, Rahu, Ketu)
	expectedPlanets := []string{"Sun", "Moon", "Mercury", "Venus", "Mars", "Jupiter", "Saturn", "Rahu", "Ketu"}
	assert.Len(T, planetsMap, len(expectedPlanets), "Expected %d planets, got %d", len(expectedPlanets), len(planetsMap))

	for _, name := range expectedPlanets {
		_, exists := planetsMap[name]
		assert.True(T, exists, "Planet %s should be present in result", name)
	}

	// Table-driven tests for all planet bala values and cord fields
	type expectedBala struct {
		// Bala values
		UdayBala     float64
		UchchaBala   float64
		VakraBala    float64
		KshetraBala  float64
		NavamshaBala float64
		// PlanetCord fields
		Longitude     float64
		Latitude      float64
		Distance      float64
		SpeedLong     float64
		SpeedLat      float64
		SpeedDist     float64
		SpeedCategory string
		Vedha         string
		Sign          string
		NakshatraName string
		NakshatraPada int
		IsRetro       bool
		SignLord      string
		SignLordship  string
	}

	table := map[string]expectedBala{
		"Sun": {
			UdayBala: 100.000000, UchchaBala: 70.734612, VakraBala: 0.000000,
			KshetraBala: 0.293985, NavamshaBala: 2.645866,
			Longitude: 270.176391, Latitude: -0.000147, Distance: 0.983651,
			SpeedLong: 1.018798, SpeedLat: -0.000002, SpeedDist: 0.000059,
			SpeedCategory: "ati-sheeghra", Vedha: "left-vedha",
			Sign: "Capricorn", NakshatraName: "Uttara Ashadha", NakshatraPada: 2,
			IsRetro: false, SignLord: "Saturn", SignLordship: "Enemy",
		},
		"Moon": {
			UdayBala: 21.054089, UchchaBala: 50.365213, VakraBala: 0.000000,
			KshetraBala: 42.685071, NavamshaBala: 15.834358,
			Longitude: 222.805521, Latitude: -5.117864, Distance: 0.002708,
			SpeedLong: 11.885197, SpeedLat: -0.056361, SpeedDist: -0.000005,
			SpeedCategory: "ati-mand", Vedha: "front-vedha",
			Sign: "Scorpio", NakshatraName: "Anuradha", NakshatraPada: 3,
			IsRetro: false, SignLord: "Mars", SignLordship: "Neutral",
		},
		"Mercury": {
			UdayBala: 0.000000, UchchaBala: 70.350375, VakraBala: 0.000000,
			KshetraBala: 14.271452, NavamshaBala: 28.443066,
			Longitude: 265.718564, Latitude: -1.758376, Distance: 1.430360,
			SpeedLong: 1.616813, SpeedLat: -0.060251, SpeedDist: 0.000222,
			SpeedCategory: "sheeghra", Vedha: "front-vedha",
			Sign: "Sagittarius", NakshatraName: "Purva Ashadha", NakshatraPada: 4,
			IsRetro: false, SignLord: "Jupiter", SignLordship: "Neutral",
		},
		"Venus": {
			UdayBala: 0.000000, UchchaBala: 77.204950, VakraBala: 0.000000,
			KshetraBala: 10.299744, NavamshaBala: 57.302308,
			Longitude: 272.059949, Latitude: -0.966153, Distance: 1.710186,
			SpeedLong: 1.257624, SpeedLat: -0.029822, SpeedDist: -0.000244,
			SpeedCategory: "ati-sheeghra", Vedha: "left-vedha",
			Sign: "Capricorn", NakshatraName: "Uttara Ashadha", NakshatraPada: 3,
			IsRetro: false, SignLord: "Saturn", SignLordship: "Friend",
		},
		"Mars": {
			UdayBala: 0.000000, UchchaBala: 96.850147, VakraBala: 0.000000,
			KshetraBala: 5.362313, NavamshaBala: 48.260821,
			Longitude: 268.927537, Latitude: -0.968436, Distance: 2.398485,
			SpeedLong: 0.774910, SpeedLat: -0.005088, SpeedDist: -0.001005,
			SpeedCategory: "ati-sheeghra", Vedha: "left-vedha",
			Sign: "Sagittarius", NakshatraName: "Uttara Ashadha", NakshatraPada: 1,
			IsRetro: false, SignLord: "Jupiter", SignLordship: "Friend",
		},
		"Jupiter": {
			UdayBala: 97.125283, UchchaBala: 99.643918, VakraBala: 98.124219,
			KshetraBala: 7.803136, NavamshaBala: 20.228221,
			Longitude: 85.318119, Latitude: 0.270673, Distance: 4.236052,
			SpeedLong: -0.134102, SpeedLat: 0.002190, SpeedDist: 0.001661,
			SpeedCategory: "kutil", Vedha: "right-vedha",
			Sign: "Gemini", NakshatraName: "Punarvasu", NakshatraPada: 2,
			IsRetro: true, SignLord: "Mercury", SignLordship: "Enemy",
		},
		"Saturn": {
			UdayBala: 28.906611, UchchaBala: 92.009166, VakraBala: 0.000000,
			KshetraBala: 9.574330, NavamshaBala: 6.915513,
			Longitude: 332.872299, Latitude: -2.217048, Distance: 9.925046,
			SpeedLong: 0.077886, SpeedLat: 0.002854, SpeedDist: 0.014745,
			SpeedCategory: "sama", Vedha: "front-vedha",
			Sign: "Pisces", NakshatraName: "Purva Bhadrapada", NakshatraPada: 4,
			IsRetro: false, SignLord: "Jupiter", SignLordship: "Neutral",
		},
		"Rahu": {
			UdayBala: 0.000000, UchchaBala: 61.468149, VakraBala: 100.000000,
			KshetraBala: 21.284016, NavamshaBala: 8.443853,
			Longitude: 317.229590, Latitude: 0.000000, Distance: 0.002570,
			SpeedLong: -0.052992, SpeedLat: 0.000000, SpeedDist: -0.000000,
			SpeedCategory: "vakra", Vedha: "left-vedha",
			Sign: "Aquarius", NakshatraName: "Shatabhisha", NakshatraPada: 4,
			IsRetro: true, SignLord: "Saturn", SignLordship: "Enemy",
		},
		"Ketu": {
			UdayBala: 0.000000, UchchaBala: 61.468149, VakraBala: 100.000000,
			KshetraBala: 21.284016, NavamshaBala: 8.443853,
			Longitude: 137.229590, Latitude: 0.000000, Distance: 0.002570,
			SpeedLong: -0.052992, SpeedLat: 0.000000, SpeedDist: -0.000000,
			SpeedCategory: "vakra", Vedha: "left-vedha",
			Sign: "Leo", NakshatraName: "Purva Phalguni", NakshatraPada: 2,
			IsRetro: true, SignLord: "Sun", SignLordship: "Enemy",
		},
	}

	for name, expected := range table {
		T.Run(name, func(t *testing.T) {
			planet, exists := planetsMap[name]
			if !assert.True(t, exists, "Planet %s should be present in result", name) {
				return
			}

			// Verify Cords is not nil
			if !assert.NotNil(t, planet.Cords, "%s: Cords should not be nil", name) {
				return
			}

			// -- Bala values --
			assert.InDelta(t, expected.UdayBala, planet.UdayBala, 0.001,
				"%s: UdayBala expected %f, got %f", name, expected.UdayBala, planet.UdayBala)
			assert.InDelta(t, expected.UchchaBala, planet.UchchaBala, 0.001,
				"%s: UchchaBala expected %f, got %f", name, expected.UchchaBala, planet.UchchaBala)
			assert.InDelta(t, expected.VakraBala, planet.VakraBala, 0.001,
				"%s: VakraBala expected %f, got %f", name, expected.VakraBala, planet.VakraBala)
			assert.InDelta(t, expected.KshetraBala, planet.KshetraBala, 0.001,
				"%s: KshetraBala expected %f, got %f", name, expected.KshetraBala, planet.KshetraBala)
			assert.InDelta(t, expected.NavamshaBala, planet.NavamshaBala, 0.001,
				"%s: NavamshaBala expected %f, got %f", name, expected.NavamshaBala, planet.NavamshaBala)

			// -- PlanetCord coordinate values --
			assert.Equal(t, name, planet.Cords.Name,
				"%s: Name expected %s, got %s", name, name, planet.Cords.Name)
			assert.InDelta(t, expected.Longitude, planet.Cords.Longitude, 0.0001,
				"%s: Longitude expected %f, got %f", name, expected.Longitude, planet.Cords.Longitude)
			assert.InDelta(t, expected.Latitude, planet.Cords.Latitude, 0.0001,
				"%s: Latitude expected %f, got %f", name, expected.Latitude, planet.Cords.Latitude)
			assert.InDelta(t, expected.Distance, planet.Cords.Distance, 0.0001,
				"%s: Distance expected %f, got %f", name, expected.Distance, planet.Cords.Distance)
			assert.InDelta(t, expected.SpeedLong, planet.Cords.SpeedLong, 0.0001,
				"%s: SpeedLong expected %f, got %f", name, expected.SpeedLong, planet.Cords.SpeedLong)
			assert.InDelta(t, expected.SpeedLat, planet.Cords.SpeedLat, 0.0001,
				"%s: SpeedLat expected %f, got %f", name, expected.SpeedLat, planet.Cords.SpeedLat)
			assert.InDelta(t, expected.SpeedDist, planet.Cords.SpeedDist, 0.0001,
				"%s: SpeedDist expected %f, got %f", name, expected.SpeedDist, planet.Cords.SpeedDist)

			// -- Derived PlanetCord fields --
			assert.Equal(t, expected.SpeedCategory, planet.Cords.SpeedCategory,
				"%s: SpeedCategory expected %s, got %s", name, expected.SpeedCategory, planet.Cords.SpeedCategory)
			assert.Equal(t, expected.Vedha, planet.Cords.Vedha,
				"%s: Vedha expected %s, got %s", name, expected.Vedha, planet.Cords.Vedha)
			assert.Equal(t, expected.Sign, planet.Cords.Sign,
				"%s: Sign expected %s, got %s", name, expected.Sign, planet.Cords.Sign)
			assert.Equal(t, expected.NakshatraName, planet.Cords.Nakshatra.Name,
				"%s: Nakshatra.Name expected %s, got %s", name, expected.NakshatraName, planet.Cords.Nakshatra.Name)
			assert.Equal(t, expected.NakshatraPada, planet.Cords.Nakshatra.Pada,
				"%s: Nakshatra.Pada expected %d, got %d", name, expected.NakshatraPada, planet.Cords.Nakshatra.Pada)
			assert.Equal(t, expected.IsRetro, planet.Cords.IsRetro,
				"%s: IsRetro expected %v, got %v", name, expected.IsRetro, planet.Cords.IsRetro)
			assert.Equal(t, expected.SignLord, planet.Cords.SignLord,
				"%s: SignLord expected %s, got %s", name, expected.SignLord, planet.Cords.SignLord)
			assert.Equal(t, expected.SignLordship, planet.Cords.SignLordship,
				"%s: SignLordship expected %s, got %s", name, expected.SignLordship, planet.Cords.SignLordship)

			// -- DMS fields should be populated (non-zero struct) --
			emptyDMS := baselib.DMS{}
			assert.NotEqual(t, emptyDMS, planet.Cords.LongitudeDMS,
				"%s: LongitudeDMS should not be zero value", name)
			assert.NotEqual(t, emptyDMS, planet.Cords.SpeedLongDMS,
				"%s: SpeedLongDMS should not be zero value", name)
		})
	}
}
