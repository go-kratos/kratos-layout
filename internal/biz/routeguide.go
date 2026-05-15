package biz

import (
	"context"
	"math"
	"time"

	v1 "github.com/go-kratos/kratos-layout/api/routeguide/v1"
)

const earthRadiusMeters = 6371000

// RouteGuideRepo is a route guide repo.
type RouteGuideRepo interface {
	GetFeature(context.Context, *v1.Point) (*v1.Feature, error)
	ListFeatures(context.Context) ([]*v1.Feature, error)
	AppendRouteNote(context.Context, *v1.RouteNote) ([]*v1.RouteNote, error)
}

// RouteGuideUsecase is a route guide usecase.
type RouteGuideUsecase struct {
	repo RouteGuideRepo
}

// NewRouteGuideUsecase new a route guide usecase.
func NewRouteGuideUsecase(repo RouteGuideRepo) *RouteGuideUsecase {
	return &RouteGuideUsecase{repo: repo}
}

// GetFeature obtains the feature at a given point.
func (uc *RouteGuideUsecase) GetFeature(ctx context.Context, point *v1.Point) (*v1.Feature, error) {
	return uc.repo.GetFeature(ctx, point)
}

// ListFeatures lists the features within a rectangle.
func (uc *RouteGuideUsecase) ListFeatures(ctx context.Context, rect *v1.Rectangle) ([]*v1.Feature, error) {
	features, err := uc.repo.ListFeatures(ctx)
	if err != nil {
		return nil, err
	}
	if rect == nil || rect.GetLo() == nil || rect.GetHi() == nil {
		return features, nil
	}

	result := make([]*v1.Feature, 0, len(features))
	for _, feature := range features {
		if feature.GetName() == "" || !inRange(feature.GetLocation(), rect) {
			continue
		}
		result = append(result, feature)
	}
	return result, nil
}

// RecordRoute summarizes a traversed route.
func (uc *RouteGuideUsecase) RecordRoute(ctx context.Context, points []*v1.Point, elapsed time.Duration) (*v1.RouteSummary, error) {
	summary := &v1.RouteSummary{
		PointCount:  int32(len(points)),
		ElapsedTime: int32(elapsed.Seconds()),
	}

	var previous *v1.Point
	for _, point := range points {
		feature, err := uc.repo.GetFeature(ctx, point)
		if err != nil {
			return nil, err
		}
		if feature.GetName() != "" {
			summary.FeatureCount++
		}
		if previous != nil {
			summary.Distance += calcDistance(previous, point)
		}
		previous = point
	}
	return summary, nil
}

// RouteChat appends a route note and returns prior notes at the same location.
func (uc *RouteGuideUsecase) RouteChat(ctx context.Context, note *v1.RouteNote) ([]*v1.RouteNote, error) {
	return uc.repo.AppendRouteNote(ctx, note)
}

func inRange(point *v1.Point, rect *v1.Rectangle) bool {
	if point == nil {
		return false
	}
	left := min(rect.GetLo().GetLongitude(), rect.GetHi().GetLongitude())
	right := max(rect.GetLo().GetLongitude(), rect.GetHi().GetLongitude())
	top := max(rect.GetLo().GetLatitude(), rect.GetHi().GetLatitude())
	bottom := min(rect.GetLo().GetLatitude(), rect.GetHi().GetLatitude())

	return point.GetLongitude() >= left &&
		point.GetLongitude() <= right &&
		point.GetLatitude() >= bottom &&
		point.GetLatitude() <= top
}

func calcDistance(p1, p2 *v1.Point) int32 {
	if p1 == nil || p2 == nil {
		return 0
	}
	lat1 := toRadians(float64(p1.GetLatitude()) / 1e7)
	lat2 := toRadians(float64(p2.GetLatitude()) / 1e7)
	lng1 := toRadians(float64(p1.GetLongitude()) / 1e7)
	lng2 := toRadians(float64(p2.GetLongitude()) / 1e7)

	dlat := lat2 - lat1
	dlng := lng2 - lng1
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return int32(earthRadiusMeters * c)
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
