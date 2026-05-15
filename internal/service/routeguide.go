package service

import (
	"context"
	"errors"
	"io"
	"time"

	v1 "github.com/go-kratos/kratos-layout/api/routeguide/v1"
	"github.com/go-kratos/kratos-layout/internal/biz"
)

// RouteGuideService is a route guide service.
type RouteGuideService struct {
	v1.UnimplementedRouteGuideServer

	uc *biz.RouteGuideUsecase
}

// NewRouteGuideService new a route guide service.
func NewRouteGuideService(uc *biz.RouteGuideUsecase) *RouteGuideService {
	return &RouteGuideService{uc: uc}
}

// GetFeature implements routeguide.RouteGuideServer.
func (s *RouteGuideService) GetFeature(ctx context.Context, in *v1.Point) (*v1.Feature, error) {
	return s.uc.GetFeature(ctx, in)
}

// ListFeatures implements routeguide.RouteGuideServer.
func (s *RouteGuideService) ListFeatures(in *v1.Rectangle, stream v1.RouteGuide_ListFeaturesServer) error {
	features, err := s.uc.ListFeatures(stream.Context(), in)
	if err != nil {
		return err
	}
	for _, feature := range features {
		if err := stream.Send(feature); err != nil {
			return err
		}
	}
	return nil
}

// RecordRoute implements routeguide.RouteGuideServer.
func (s *RouteGuideService) RecordRoute(stream v1.RouteGuide_RecordRouteServer) error {
	start := time.Now()
	points := make([]*v1.Point, 0)
	for {
		point, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			summary, err := s.uc.RecordRoute(stream.Context(), points, time.Since(start))
			if err != nil {
				return err
			}
			return stream.SendAndClose(summary)
		}
		if err != nil {
			return err
		}
		points = append(points, point)
	}
}

// RouteChat implements routeguide.RouteGuideServer.
func (s *RouteGuideService) RouteChat(stream v1.RouteGuide_RouteChatServer) error {
	for {
		note, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		notes, err := s.uc.RouteChat(stream.Context(), note)
		if err != nil {
			return err
		}
		for _, note := range notes {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}
