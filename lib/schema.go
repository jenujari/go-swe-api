package lib

import (
	baselib "github.com/jenujari/planets-lib"
)

type PlanetBalas struct {
	Cords        *baselib.PlanetCord `json:"cords"`
	UdayBala     float64             `json:"uday_bala"`
	UchchaBala   float64             `json:"uchcha_bala"`
	VakraBala    float64             `json:"vakra_bala"`
	KshetraBala  float64             `json:"kshetra_bala"`
	NavamshaBala float64             `json:"navamsha_bala"`
}
