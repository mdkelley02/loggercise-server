package loggerciseHandler

import (
	"context"
	"fmt"
	store "loggercise/internal/store"

	loggerciseProto "loggercise/gen/go/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type LoggerciseHandlerIF interface {
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

type LoggerciseHandler struct {
	log   *logrus.Logger
	store store.StoreIF
}

func NewLoggerciseHandler(log *logrus.Logger, store store.StoreIF) *LoggerciseHandler {
	return &LoggerciseHandler{
		log:   log,
		store: store,
	}
}

func (svc *LoggerciseHandler) UpsertWorkout(ctx context.Context, req *loggerciseProto.UpsertWorkoutRequest) (*loggerciseProto.WorkoutResponse, error) {
	res, err := svc.store.UpsertWorkout(ctx, req)
	if err != nil {
		svc.log.Errorf("Error inserting workout: %v", err)
		return nil, err
	}
	return res, nil
}

func (svc *LoggerciseHandler) GetWorkouts(ctx context.Context, req *loggerciseProto.GetWorkoutsRequest) (*loggerciseProto.WorkoutResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	req.UserId = userId
	res, err := svc.store.GetWorkouts(ctx, req)
	if err != nil {
		svc.log.Errorf("Error getting workouts: %v", err)
		return nil, err
	}
	return res, nil
}

func (svc *LoggerciseHandler) DeleteWorkout(ctx context.Context, req *loggerciseProto.WorkoutRequest) (*loggerciseProto.WorkoutResponse, error) {
	res, err := svc.store.DeleteWorkout(ctx, req)
	if err != nil {
		svc.log.Errorf("Error deleting workout: %v", err)
		return nil, err
	}
	return res, nil
}

func (svc *LoggerciseHandler) UpsertExercise(ctx context.Context, req *loggerciseProto.UpsertExerciseRequest) (*loggerciseProto.ExerciseResponse, error) {
	res, err := svc.store.UpsertExercise(ctx, req)
	if err != nil {
		svc.log.Errorf("Error inserting exercise: %v", err)
		return nil, err
	}
	return res, nil
}

func (svc *LoggerciseHandler) UpsertSet(ctx context.Context, req *loggerciseProto.UpsertSetRequest) (*loggerciseProto.Empty, error) {
	return nil, nil
}

func (svc *LoggerciseHandler) DeleteSet(ctx context.Context, req *loggerciseProto.DeleteSetRequest) (*loggerciseProto.Empty, error) {
	return nil, nil
}

func (svc *LoggerciseHandler) DeleteExercise(ctx context.Context, req *loggerciseProto.ExerciseRequest) (*loggerciseProto.ExerciseResponse, error) {
	res, err := svc.store.DeleteExercise(ctx, req)
	if err != nil {
		svc.log.Errorf("Error deleting exercise: %v", err)
		return nil, err
	}
	return res, nil
}

func (svc *LoggerciseHandler) UpsertLift(ctx context.Context, req *loggerciseProto.UpsertLiftRequest) (*loggerciseProto.LiftResponse, error) {
	return nil, nil
}

func (svc *LoggerciseHandler) GetLifts(ctx context.Context, req *loggerciseProto.Empty) (*loggerciseProto.LiftResponse, error) {
	res, err := svc.store.GetLifts(ctx, req)
	if err != nil {
		svc.log.Errorf("Error getting lifts: %v", err)
		return nil, err
	}
	return res, nil
}

func (svc *LoggerciseHandler) GetMuscles(ctx context.Context, req *loggerciseProto.Empty) (*loggerciseProto.MusclesResponse, error) {
	res, err := svc.store.GetMuscles(ctx, req)
	if err != nil {
		svc.log.Errorf("Error getting muscles: %v", err)
		return nil, err
	}
	return res, nil
}

func getUserId(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing metadata")
	}
	userID, ok := md["user-id"]
	if !ok {
		return "", fmt.Errorf("missing user-id")
	}
	return userID[0], nil
}
