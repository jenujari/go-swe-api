package lib

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	baselib "github.com/jenujari/planets-lib"
	swelib "github.com/mshafiee/swephgo"
)

const (
	iFlag = swelib.SeflgSwieph | swelib.SeflgSpeed | swelib.SeflgSidereal
)

var (
	SWISSEPH_PATH string = ""
)

func init() {
	InitSwelib()

	sweVer := make([]byte, 12)
	swelib.Version(sweVer)

	fmt.Printf("Library used: Swiss Ephemeris v%s\n", string(bytes.Trim(sweVer, "\x00")))
}

func GetVersion() string {
	sweVer := make([]byte, 12)
	swelib.Version(sweVer)

	return string(sweVer)
}

func InitSwelib() {
	setEphePath()
	setSidMode()
	// swelib.SetTopo(-0.118092, 51.509865, 0) // London coordinates
}

func setEphePath() {
	SWISSEPH_PATH := os.Getenv("SWISSEPH_PATH")

	if SWISSEPH_PATH == "" {
		SWISSEPH_PATH = "/usr/local/lib/ephe"
		log.Println("SWISSEPH_PATH not set, using default path:", SWISSEPH_PATH)
	}

	swelib.SetEphePath([]byte(SWISSEPH_PATH))
}

func setSidMode() {
	// swelib.SetSidMode(swelib.SeSidmKrishnamurti, 0, 0)
	swelib.SetSidMode(swelib.SeSidmLahiri, 0, 0)
}

func SweClear() {
	swelib.Close()
}

func LongDiff(jd float64, p1, p2 int) (float64, error) {
	defer SweClear()

	var a = make([]float64, 6)
	var b = make([]float64, 6)
	var serr1 = make([]byte, 1000)
	var serr2 = make([]byte, 1000)

	setSidMode()
	ret := swelib.CalcUt(jd, p1, int(iFlag), a, serr1)

	serr1 = bytes.Trim(serr1, "\x00")
	errStr := string(serr1)
	if len(errStr) > 0 {
		return 0, fmt.Errorf("error calculating planet position: %s", errStr)
	}

	if ret < 0 {
		return 0, fmt.Errorf("error calculating planet position: %d", ret)
	}

	setSidMode()
	ret = swelib.CalcUt(jd, p2, int(iFlag), b, serr2)

	serr2 = bytes.Trim(serr2, "\x00")
	errStr = string(serr2)
	if len(errStr) > 0 {
		return 0, fmt.Errorf("error calculating planet position: %s", errStr)
	}

	if ret < 0 {
		return 0, fmt.Errorf("error calculating planet position: %d", ret)
	}

	d := math.Abs(a[0] - b[0])
	if d > 180 {
		d = 360 - d
	}

	return d, nil
}

func FindConjunctionRange(startTime, endTime time.Time, conjDeg, stepDays float64, p1, p2 int) (time.Time, time.Time, bool) {
	// orbDeg := 1.0          // conjunction orb in degrees
	// stepDays := 1.0 / 24.0 // 1 hour step
	inConj := false
	var startJD, endJD float64
	var emptyTime time.Time

	jdStart, _ := UTCToSiderealTime(startTime)
	jdEnd, _ := UTCToSiderealTime(endTime)

	for jd := jdStart; jd <= jdEnd; jd += stepDays {
		diff, er := LongDiff(jd, p1, p2)
		if er != nil {
			return emptyTime, emptyTime, false
		}
		if diff <= conjDeg {
			if !inConj {
				startJD = jd
				inConj = true
			}
			endJD = jd
		} else if inConj {
			retStart, _ := SiderealTimeToUTC(startJD)
			retEnd, _ := SiderealTimeToUTC(endJD)
			return retStart, retEnd, true
		}
	}

	return emptyTime, emptyTime, false
}

func UTCToSiderealTime(utcTime time.Time) (float64, error) {
	defer SweClear()

	var jd float64
	var tjdArr = make([]float64, 1)
	var serr = make([]byte, 1000)

	y, m, d := utcTime.Year(), int(utcTime.Month()), utcTime.Day()
	h, min, sec := utcTime.Hour(), utcTime.Minute(), utcTime.Second()

	setSidMode()
	swelib.UtcToJd(y, m, d, h, min, float64(sec), swelib.SeGregCal, tjdArr, serr)

	serr = bytes.Trim(serr, "\x00")
	errStr := string(serr)
	if len(errStr) > 0 {
		return 0, fmt.Errorf("error converting UTC to JD: %s", errStr)
	}

	jd = tjdArr[0]

	return jd, nil
}

func SiderealTimeToUTC(siderealTime float64) (time.Time, error) {
	defer SweClear()

	yArr := make([]int, 1)
	mArr := make([]int, 1)
	dArr := make([]int, 1)
	utArr := make([]float64, 1)

	deltaT := swelib.Deltat(siderealTime)
	swelib.Revjul((siderealTime - deltaT), swelib.SeGregCal, yArr, mArr, dArr, utArr)

	if len(yArr) == 0 || len(mArr) == 0 || len(dArr) == 0 || len(utArr) == 0 {
		return time.Time{}, fmt.Errorf("error converting sidereal time to UTC: %v %v %v %v", yArr, mArr, dArr, utArr)
	}

	y := yArr[0]
	m := mArr[0]
	d := dArr[0]
	ut := utArr[0]

	hours := int(ut)
	minF := (ut - float64(hours)) * 60
	minutes := int(minF)
	seconds := int(math.Round((minF - float64(minutes)) * 60))

	// normalize rounding overflow
	if seconds == 60 {
		seconds = 0
		minutes++
	}

	if minutes == 60 {
		minutes = 0
		hours++
	}

	return time.Date(y, time.Month(m), d, hours, minutes, seconds, 0, time.UTC), nil
}

func GetPlanetCalculation(siderealTime float64, planet string) (*baselib.PlanetCord, error) {
	defer SweClear()
	var xp = make([]float64, 6)
	var serr = make([]byte, 1000)
	var isKetu bool = false

	if planet == "Ketu" {
		isKetu = true
	}

	libPlanet := baselib.PLANET_LIB_MAP[planet]

	setSidMode()
	ret := swelib.CalcUt(siderealTime, libPlanet, int(iFlag), xp, serr)

	serr = bytes.Trim(serr, "\x00")
	errStr := string(serr)
	if len(errStr) > 0 {
		return nil, fmt.Errorf("error calculating planet position: %s", errStr)
	}

	if ret < 0 {
		return nil, fmt.Errorf("error calculating planet position: %d", ret)
	}

	// Ensure all requested flag bits are present in the returned flags.
	// The library may set additional flags; accept those as long as
	// the requested bits were honored.
	if (ret & int32(iFlag)) != int32(iFlag) {
		return nil, fmt.Errorf("calculation did not include requested flags, Requested: %x, Returned: %x", iFlag, ret)
	}

	if isKetu {
		xp[0] = xp[0] + 180.0
		if xp[0] > 360.0 {
			xp[0] = xp[0] - 360.0
		}
	}

	planetCord := new(baselib.PlanetCord)
	planetCord.Longitude = xp[0]
	planetCord.Latitude = xp[1]
	planetCord.Distance = xp[2]
	planetCord.SpeedLong = xp[3]
	planetCord.SpeedLat = xp[4]
	planetCord.SpeedDist = xp[5]

	return planetCord, nil
}
