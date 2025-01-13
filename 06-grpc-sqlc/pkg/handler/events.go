package handler

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thylong/go-templates/06-grpc-sqlc/pkg/db" // Import sqlc-generated code
	eventpb "github.com/thylong/go-templates/06-grpc-sqlc/pkg/proto/events"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventServiceServer struct {
	eventpb.UnimplementedEventServiceServer
	queries *db.Queries
}

// NewEventServiceServer initializes the server with sqlc queries
func NewEventServiceServer(queries *db.Queries) *EventServiceServer {
	return &EventServiceServer{queries: queries}
}

// GetEvents handles fetching paginated and filtered events
func (s *EventServiceServer) GetEvents(ctx context.Context, req *eventpb.GetEventsRequest) (*eventpb.GetEventsResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	events, err := s.queries.GetEvents(ctx, db.GetEventsParams{
		Column1: pgtype.Text{String: req.Search, Valid: true},
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch events: %w", err)
	}

	totalCount, err := s.queries.GetEventsCount(ctx, pgtype.Text{String: req.Search, Valid: true})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch event count: %w", err)
	}

	var pbEvents []*eventpb.Event
	for _, e := range events {
		startAt, err := convertToProtoDateTime(e.StartAt)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to convert start time: %w", err)
		}
		pbEvents = append(pbEvents, &eventpb.Event{
			EventId:      e.EventID.String(),
			EventPrivacy: eventpb.EventPrivacy(e.EventPrivacy),
			Name:         e.Name,
			Type:         e.Type,
			Department:   e.Department,
			Regions:      e.Regions,
			Tags:         e.Tags,
			StartAt:      startAt,
		})
	}

	return &eventpb.GetEventsResponse{
		Events:     pbEvents,
		TotalCount: int32(totalCount),
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetEvent handles fetching a single event by ID
func (s *EventServiceServer) GetEvent(ctx context.Context, req *eventpb.GetEventRequest) (*eventpb.GetEventResponse, error) {
	// Parse and validate the UUID
	parsedUUID, err := uuid.Parse(req.EventId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID: %w", err)
	}

	// Create a pgtype.UUID instance
	eventID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	event, err := s.queries.GetEventByID(ctx, eventID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "event not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch event: %w", err)
	}

	startAt, err := convertToProtoDateTime(event.StartAt)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to convert start time: %w", err)
	}

	return &eventpb.GetEventResponse{
		Event: &eventpb.Event{
			EventId:      event.EventID.String(),
			EventPrivacy: eventpb.EventPrivacy(event.EventPrivacy),
			Name:         event.Name,
			Type:         event.Type,
			Department:   event.Department,
			Regions:      event.Regions,
			Tags:         event.Tags,
			StartAt:      startAt,
		},
	}, nil
}

// PutEvent handles inserting a new event
func (s *EventServiceServer) PutEvent(ctx context.Context, req *eventpb.PutEventRequest) (*eventpb.PutEventResponse, error) {
	if req.StartAt == nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to insert event: invalid or missing StartAt")
	}
	startAt := convertFromProtoDateTime(req.StartAt)

	event, err := s.queries.InsertEvent(ctx, db.InsertEventParams{
		EventPrivacy: int32(req.EventPrivacy),
		Name:         req.Name,
		Type:         req.Type,
		Department:   req.Department,
		Regions:      req.Regions,
		Tags:         req.Tags,
		StartAt:      pgtype.Timestamp{Time: startAt, Valid: true},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert event: %w", err)
	}

	startAtProto, _ := convertToProtoDateTime(event.StartAt)

	return &eventpb.PutEventResponse{
		Event: &eventpb.Event{
			EventId:      event.EventID.String(),
			EventPrivacy: eventpb.EventPrivacy(event.EventPrivacy),
			Name:         event.Name,
			Type:         event.Type,
			Department:   event.Department,
			Regions:      event.Regions,
			Tags:         event.Tags,
			StartAt:      startAtProto,
		},
	}, nil
}

// DeleteEvent handles deleting an event by ID
func (s *EventServiceServer) DeleteEvent(ctx context.Context, req *eventpb.DeleteEventRequest) (*eventpb.DeleteEventResponse, error) {
	// Parse and validate the UUID
	parsedUUID, err := uuid.Parse(req.EventID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID: %w", err)
	}

	// Create a pgtype.UUID instance
	eventID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	err = s.queries.DeleteEvent(ctx, eventID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete event: %w", err)
	}
	return &eventpb.DeleteEventResponse{}, nil
}
