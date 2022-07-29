package utils

import (
	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

func ParsePolygon(polygon ewkb.Polygon) []*pb.Polygon {
	var polygonResponse []*pb.Polygon
	for _, e := range polygon.Coords() {
		var coordinates = []float64{}
		for _, e2 := range e {
			coordinates = append(coordinates, e2.X())
			coordinates = append(coordinates, e2.Y())
		}
		polygonResponse = append(polygonResponse, &pb.Polygon{Coordinates: coordinates})
	}
	// polygonResponse = append(polygonResponse, &pb.Polygon{Coordinates: polygon.FlatCoords()})
	return polygonResponse
}

// for (var item in parameter) {
//     list.add(Polygon(coordinates: [item[0], item[1]]));
// }
// return list;
