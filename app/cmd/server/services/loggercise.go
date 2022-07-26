package services

import (
	"context"
	loggerciseProto "loggercise/gen/go/service"
	"loggercise/internal/loggerciseHandler"

	"github.com/sirupsen/logrus"
)

type LoggerciseServiceIF interface {
	GetWorkouts(ctx context.Context, req *loggerciseProto.GetWorkoutsRequest) (*loggerciseProto.WorkoutResponse, error)
	UpsertWorkout(ctx context.Context, req *loggerciseProto.UpsertWorkoutRequest) (*loggerciseProto.WorkoutResponse, error)
	DeleteWorkout(ctx context.Context, req *loggerciseProto.WorkoutRequest) (*loggerciseProto.WorkoutResponse, error)
	UpsertExercise(ctx context.Context, req *loggerciseProto.UpsertExerciseRequest) (*loggerciseProto.ExerciseResponse, error)
	DeleteExercise(ctx context.Context, req *loggerciseProto.ExerciseRequest) (*loggerciseProto.ExerciseResponse, error)
	UpsertSet(ctx context.Context, req *loggerciseProto.UpsertSetRequest) (*loggerciseProto.Empty, error)
	DeleteSet(ctx context.Context, req *loggerciseProto.DeleteSetRequest) (*loggerciseProto.Empty, error)
	UpsertLift(ctx context.Context, req *loggerciseProto.UpsertLiftRequest) (*loggerciseProto.LiftResponse, error)
	GetMuscles(ctx context.Context, req *loggerciseProto.Empty) (*loggerciseProto.MusclesResponse, error)
	GetLifts(ctx context.Context, req *loggerciseProto.Empty) (*loggerciseProto.LiftResponse, error)
}

type LoggerciseService struct {
	log     *logrus.Logger
	handler loggerciseHandler.LoggerciseHandlerIF
	loggerciseProto.UnimplementedLoggerciseServer
}

func NewLoggerciseService(log *logrus.Logger, handler loggerciseHandler.LoggerciseHandlerIF) *LoggerciseService {
	service := &LoggerciseService{
		log:     log,
		handler: handler,
	}
	return service
}

func (svc *LoggerciseService) UpsertWorkout(ctx context.Context, req *loggerciseProto.UpsertWorkoutRequest) (*loggerciseProto.WorkoutResponse, error) {
	return svc.handler.UpsertWorkout(ctx, req)
}

func (svc *LoggerciseService) GetWorkouts(ctx context.Context, req *loggerciseProto.GetWorkoutsRequest) (*loggerciseProto.WorkoutResponse, error) {
	return svc.handler.GetWorkouts(ctx, req)
}

func (svc *LoggerciseService) DeleteWorkout(ctx context.Context, req *loggerciseProto.WorkoutRequest) (*loggerciseProto.WorkoutResponse, error) {
	return svc.handler.DeleteWorkout(ctx, req)
}

func (svc *LoggerciseService) UpsertExercise(ctx context.Context, req *loggerciseProto.UpsertExerciseRequest) (*loggerciseProto.ExerciseResponse, error) {
	return svc.handler.UpsertExercise(ctx, req)
}

func (svc *LoggerciseService) DeleteExercise(ctx context.Context, req *loggerciseProto.ExerciseRequest) (*loggerciseProto.ExerciseResponse, error) {
	return svc.handler.DeleteExercise(ctx, req)
}

func (svc *LoggerciseService) UpsertSet(ctx context.Context, req *loggerciseProto.UpsertSetRequest) (*loggerciseProto.Empty, error) {
	return svc.handler.UpsertSet(ctx, req)
}

func (svc *LoggerciseService) DeleteSet(ctx context.Context, req *loggerciseProto.DeleteSetRequest) (*loggerciseProto.Empty, error) {
	return svc.handler.DeleteSet(ctx, req)
}

func (svc *LoggerciseService) UpsertLift(ctx context.Context, req *loggerciseProto.UpsertLiftRequest) (*loggerciseProto.LiftResponse, error) {
	return svc.handler.UpsertLift(ctx, req)
}

func (svc *LoggerciseService) GetLifts(ctx context.Context, req *loggerciseProto.Empty) (*loggerciseProto.LiftResponse, error) {
	return svc.handler.GetLifts(ctx, req)
}

func (svc *LoggerciseService) GetMuscles(ctx context.Context, req *loggerciseProto.Empty) (*loggerciseProto.MusclesResponse, error) {
	return svc.handler.GetMuscles(ctx, req)
}
