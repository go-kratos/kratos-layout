package data

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	v1 "github.com/go-kratos/kratos-layout/api/routeguide/v1"
	"github.com/go-kratos/kratos-layout/internal/biz"
)

type routeGuideRepo struct {
	data *Data
	log  *slog.Logger

	mu         sync.RWMutex
	features   []*v1.Feature
	routeNotes map[string][]*v1.RouteNote
}

// NewRouteGuideRepo new a route guide repo.
func NewRouteGuideRepo(data *Data, logger *slog.Logger) biz.RouteGuideRepo {
	return &routeGuideRepo{
		data: data,
		log:  logger,
		features: []*v1.Feature{
			{Name: "Berkshire Valley Management Area Trail", Location: &v1.Point{Latitude: 409146138, Longitude: -746188906}},
			{Name: "Patriots Path", Location: &v1.Point{Latitude: 407838351, Longitude: -746143763}},
			{Name: "Great Piece Meadows", Location: &v1.Point{Latitude: 407235788, Longitude: -747160458}},
			{Name: "Farny Highlands Trail", Location: &v1.Point{Latitude: 410248224, Longitude: -747127767}},
		},
		routeNotes: make(map[string][]*v1.RouteNote),
	}
}

func (r *routeGuideRepo) GetFeature(_ context.Context, point *v1.Point) (*v1.Feature, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, feature := range r.features {
		if samePoint(feature.GetLocation(), point) {
			return cloneFeature(feature), nil
		}
	}
	return &v1.Feature{Location: clonePoint(point)}, nil
}

func (r *routeGuideRepo) ListFeatures(context.Context) ([]*v1.Feature, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	features := make([]*v1.Feature, 0, len(r.features))
	for _, feature := range r.features {
		features = append(features, cloneFeature(feature))
	}
	return features, nil
}

func (r *routeGuideRepo) AppendRouteNote(_ context.Context, note *v1.RouteNote) ([]*v1.RouteNote, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := pointKey(note.GetLocation())
	notes := cloneRouteNotes(r.routeNotes[key])
	r.routeNotes[key] = append(r.routeNotes[key], cloneRouteNote(note))
	return notes, nil
}

func samePoint(a, b *v1.Point) bool {
	return a.GetLatitude() == b.GetLatitude() && a.GetLongitude() == b.GetLongitude()
}

func pointKey(point *v1.Point) string {
	return fmt.Sprintf("%d:%d", point.GetLatitude(), point.GetLongitude())
}

func clonePoint(point *v1.Point) *v1.Point {
	if point == nil {
		return nil
	}
	return &v1.Point{
		Latitude:  point.GetLatitude(),
		Longitude: point.GetLongitude(),
	}
}

func cloneFeature(feature *v1.Feature) *v1.Feature {
	if feature == nil {
		return nil
	}
	return &v1.Feature{
		Name:     feature.GetName(),
		Location: clonePoint(feature.GetLocation()),
	}
}

func cloneRouteNote(note *v1.RouteNote) *v1.RouteNote {
	if note == nil {
		return nil
	}
	return &v1.RouteNote{
		Location: clonePoint(note.GetLocation()),
		Message:  note.GetMessage(),
	}
}

func cloneRouteNotes(notes []*v1.RouteNote) []*v1.RouteNote {
	out := make([]*v1.RouteNote, 0, len(notes))
	for _, note := range notes {
		out = append(out, cloneRouteNote(note))
	}
	return out
}
