package grpc

import (
	"context"
	"time"

	lib "github.com/jenujari/go-swe-api/lib"
	pb "github.com/jenujari/go-swe-api/proto"
	baselib "github.com/jenujari/planets-lib"
)

type Server struct {
	pb.UnimplementedEphServiceServer
}

func (s *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   lib.GetVersion(),
	}, nil
}

func (s *Server) GetPos(ctx context.Context, req *pb.PosRequest) (*pb.PosResponse, error) {
	t, err := time.Parse(time.RFC3339, req.Time)
	if err != nil {
		return nil, err
	}

	siderealTime, err := lib.UTCToSiderealTime(t)
	if err != nil {
		return nil, err
	}

	results := make(map[string]*pb.PlanetCord)

	if req.PlanetName == "" {
		for planet := range baselib.PLANET_LIB_MAP {
			planetCord, err := lib.GetPlanetCalculation(siderealTime, planet)
			if err != nil {
				return nil, err
			}
			planetCord.CalculateDerivedValues()
			results[planet] = mapToProtoPlanetCord(planetCord)
		}
	} else {
		planetCord, err := lib.GetPlanetCalculation(siderealTime, req.PlanetName)
		if err != nil {
			return nil, err
		}
		planetCord.CalculateDerivedValues()
		results[req.PlanetName] = mapToProtoPlanetCord(planetCord)
	}

	return &pb.PosResponse{Results: results}, nil
}

func (s *Server) FindConjunction(ctx context.Context, req *pb.ConjunctionRequest) (*pb.ConjunctionResponse, error) {
	startTime, err := time.Parse(time.RFC3339, req.Start)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(time.RFC3339, req.End)
	if err != nil {
		return nil, err
	}

	orb := float64(req.Orb)
	if orb <= 0 {
		orb = 1.0 // default orb value if not provided or invalid
	}

	step := req.Step
	if step <= 0 {
		step = 1.0 / 24.0 // default step value (1 hour) if not provided or invalid
	}

	startConj, endConj, found, err := lib.FindConjunctionRange(
		startTime,
		endTime,
		orb,
		step,
		baselib.PLANET_LIB_MAP[req.Planet1],
		baselib.PLANET_LIB_MAP[req.Planet2],
	)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, context.DeadlineExceeded
	}

	return &pb.ConjunctionResponse{
		Start: startConj.Format(time.RFC3339),
		End:   endConj.Format(time.RFC3339),
	}, nil
}

func (s *Server) Tithy(ctx context.Context, req *pb.TithyRequest) (*pb.TithyResponse, error) {
	timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		return nil, err
	}

	tithy, week, nakshatra, err := lib.CalcTithy(timestamp)
	if err != nil {
		return nil, err
	}

	return &pb.TithyResponse{
		Tithy:     tithy,
		Nakshatra: nakshatra,
		Weekday:   week,
	}, nil
}

func (s *Server) GetBalas(ctx context.Context, req *pb.BalasRequest) (*pb.BalasResponse, error) {
	t, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		return nil, err
	}

	balasMap, err := lib.GetAllPlanetsBalas(t)
	if err != nil {
		return nil, err
	}

	results := make(map[string]*pb.PlanetBalas)
	for name, pbBalas := range balasMap {
		results[name] = &pb.PlanetBalas{
			Cords:        mapToProtoPlanetCord(pbBalas.Cords),
			UdayBala:     pbBalas.UdayBala,
			UchchaBala:   pbBalas.UchchaBala,
			VakraBala:    pbBalas.VakraBala,
			KshetraBala:  pbBalas.KshetraBala,
			NavamshaBala: pbBalas.NavamshaBala,
		}
	}

	return &pb.BalasResponse{Results: results}, nil
}

func mapToProtoPlanetCord(pc *baselib.PlanetCord) *pb.PlanetCord {
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
		LongitudeDms: &pb.DMS{
			IsNegative: pc.LongitudeDMS.IsNegative,
			D:          int32(pc.LongitudeDMS.D),
			M:          int32(pc.LongitudeDMS.M),
			S:          pc.LongitudeDMS.S,
		},
		LatitudeDms: &pb.DMS{
			IsNegative: pc.LatitudeDMS.IsNegative,
			D:          int32(pc.LatitudeDMS.D),
			M:          int32(pc.LatitudeDMS.M),
			S:          pc.LatitudeDMS.S,
		},
		SpeedLongDms: &pb.DMS{
			IsNegative: pc.SpeedLongDMS.IsNegative,
			D:          int32(pc.SpeedLongDMS.D),
			M:          int32(pc.SpeedLongDMS.M),
			S:          pc.SpeedLongDMS.S,
		},
		Sign: pc.Sign,
		Nakshatra: &pb.NakshatraPada{
			Name: pc.Nakshatra.Name,
			Pada: int32(pc.Nakshatra.Pada),
		},
		IsRetro:      pc.IsRetro,
		SignLord:     pc.SignLord,
		SignLordship: pc.SignLordship,
		NavamsaSign:  pc.NavamsaSign,
		Vargottama:   pc.Vargottama,
		VedhaTarget:  pc.VedhaTarget,
	}
}

