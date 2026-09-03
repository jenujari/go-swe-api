package grpc

import (
	lib "github.com/jenujari/go-swe-api/lib"
	pb "github.com/jenujari/go-swe-api/proto"
	baselib "github.com/jenujari/planets-lib"
)

func mapToProtoDMS(d baselib.DMS) *pb.DMS {
	return &pb.DMS{
		IsNegative: d.IsNegative,
		D:          int32(d.D),
		M:          int32(d.M),
		S:          d.S,
	}
}

func mapToProtoNakshatra(n baselib.NakshatraPada) *pb.NakshatraPada {
	return &pb.NakshatraPada{
		Name: n.Name,
		Pada: int32(n.Pada),
	}
}

func mapToProtoPlanetCord(pc *baselib.PlanetCord) *pb.PlanetCord {
	if pc == nil {
		return nil
	}
	return &pb.PlanetCord{
		Name:          pc.Name,
		Longitude:     pc.Longitude,
		Latitude:      pc.Latitude,
		Distance:      pc.Distance,
		SpeedLong:     pc.SpeedLong,
		SpeedLat:      pc.SpeedLat,
		SpeedDist:     pc.SpeedDist,
		SpeedCategory: pc.SpeedCategory,
		Vedha:         pc.Vedha,
		LongitudeDms:  mapToProtoDMS(pc.LongitudeDMS),
		LatitudeDms:   mapToProtoDMS(pc.LatitudeDMS),
		SpeedLongDms:  mapToProtoDMS(pc.SpeedLongDMS),
		Sign:          pc.Sign,
		Nakshatra:     mapToProtoNakshatra(pc.Nakshatra),
		IsRetro:       pc.IsRetro,
		SignLord:      pc.SignLord,
		SignLordship:  pc.SignLordship,
		NavamsaSign:   pc.NavamsaSign,
		Vargottama:    pc.Vargottama,
		VedhaTarget:   pc.VedhaTarget,
	}
}

func mapToProtoPlanetBalas(b lib.PlanetBalas) *pb.PlanetBalas {
	return &pb.PlanetBalas{
		Cords:        mapToProtoPlanetCord(b.Cords),
		UdayBala:     b.UdayBala,
		UchchaBala:   b.UchchaBala,
		VakraBala:    b.VakraBala,
		KshetraBala:  b.KshetraBala,
		NavamshaBala: b.NavamshaBala,
	}
}

func mapToProtoPlanetBalasMap(src map[string]lib.PlanetBalas) map[string]*pb.PlanetBalas {
	results := make(map[string]*pb.PlanetBalas, len(src))
	for name, balas := range src {
		results[name] = mapToProtoPlanetBalas(balas)
	}
	return results
}
