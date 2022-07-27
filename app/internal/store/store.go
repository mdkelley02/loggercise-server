package store

import (
	"context"
	"fmt"

	loggerciseProto "loggercise/gen/go/service"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sirupsen/logrus"
)

const (
	LoggerciseDatabase = "Loggercise"
	WorkoutsCollection = "workouts"
	LiftsCollection    = "lifts"
	MusclesCollection  = "muscles"
)

type StoreIF interface {
	GetWorkout(ctx context.Context, req *loggerciseProto.WorkoutRequest) (*loggerciseProto.Workout, error)
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

type Store struct {
	log      *logrus.Logger
	mongo    *mongo.Client
	workouts *mongo.Collection
	lifts    *mongo.Collection
	muscles  *mongo.Collection
}

func NewMongoClient(connString string) (*mongo.Client, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(connString).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewStore(mongo *mongo.Client, log *logrus.Logger) *Store {
	return &Store{
		log:      log,
		mongo:    mongo,
		workouts: mongo.Database(LoggerciseDatabase).Collection(WorkoutsCollection),
		lifts:    mongo.Database(LoggerciseDatabase).Collection(LiftsCollection),
		muscles:  mongo.Database(LoggerciseDatabase).Collection(MusclesCollection),
	}
}

func (svc *Store) UpsertWorkout(ctx context.Context, req *loggerciseProto.UpsertWorkoutRequest) (*loggerciseProto.WorkoutResponse, error) {
	workout := &loggerciseProto.Workout{
		Title:     req.Title,
		Notes:     req.Notes,
		Date:      req.Date,
		UserId:    req.UserId,
		WorkoutId: uuid.New().String(),
		CreatedAt: time.Now().String(),
		Exercises: []*loggerciseProto.Exercise{},
	}
	filter := bson.M{"_id": req.WorkoutId}
	if req.WorkoutId == "" {
		filter["_id"] = workout.WorkoutId
	}
	update := bson.M{
		"$set": bson.M{
			"title": workout.Title,
			"notes": workout.Notes,
			"date":  workout.Date,
		},
		"$setOnInsert": bson.M{
			"userId":    req.UserId,
			"createdAt": workout.CreatedAt,
			"exercises": workout.Exercises,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := svc.workouts.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return &loggerciseProto.WorkoutResponse{}, nil
}

func (svc *Store) GetWorkout(ctx context.Context, req *loggerciseProto.WorkoutRequest) (*loggerciseProto.Workout, error) {
	filter := bson.M{"_id": req.WorkoutId}
	result := svc.workouts.FindOne(ctx, filter)
	var workout loggerciseProto.Workout
	err := mapResultToWorkout(ctx, result, &workout)
	if err != nil {
		return nil, err
	}
	return &workout, nil
}

func (svc *Store) GetWorkouts(ctx context.Context, req *loggerciseProto.GetWorkoutsRequest) (*loggerciseProto.WorkoutResponse, error) {
	res := &loggerciseProto.WorkoutResponse{}
	var workouts []*loggerciseProto.Workout
	filter := bson.M{"userId": req.UserId}
	curr, err := svc.workouts.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = mapCursorToWorkouts(ctx, curr, &workouts)
	if err != nil {
		return nil, err
	}
	res.Workouts = workouts
	return res, nil
}

func (svc *Store) DeleteWorkout(ctx context.Context, req *loggerciseProto.WorkoutRequest) (*loggerciseProto.WorkoutResponse, error) {
	res := &loggerciseProto.WorkoutResponse{}
	filter := bson.M{"_id": req.WorkoutId}
	_, err := svc.workouts.DeleteOne(ctx, filter)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (svc *Store) UpsertExercise(ctx context.Context, req *loggerciseProto.UpsertExerciseRequest) (*loggerciseProto.ExerciseResponse, error) {
	res := &loggerciseProto.ExerciseResponse{}
	isInsert := req.ExerciseId == ""
	exercise := MongoExercise{
		ExerciseId: uuid.New().String(),
		WorkoutId:  req.WorkoutId,
		Lift: &MongoLift{
			LiftId:       req.Lift.LiftId,
			Name:         req.Lift.Name,
			MuscleGroups: req.Lift.MuscleGroups,
		},
		Sets: []*MongoSet{},
	}
	if !isInsert {
		exercise.ExerciseId = req.ExerciseId
		filter := bson.M{"_id": exercise.WorkoutId, "exercises._id": exercise.ExerciseId}
		update := bson.M{
			"$set": bson.M{
				"exercises.$": exercise,
			},
		}
		opts := options.Update().SetUpsert(true)
		_, err := svc.workouts.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return nil, err
		}
		res.Exercises = []*loggerciseProto.Exercise{{
			ExerciseId: exercise.ExerciseId,
			WorkoutId:  exercise.WorkoutId,
			Lift: &loggerciseProto.Lift{
				LiftId:       exercise.Lift.LiftId,
				Name:         exercise.Lift.Name,
				MuscleGroups: exercise.Lift.MuscleGroups,
			},
			Sets: []*loggerciseProto.ExerciseSet{},
		}}

	} else {
		filter := bson.M{"_id": exercise.WorkoutId}
		update := bson.M{
			"$push": bson.M{
				"exercises": exercise,
			},
		}
		opts := options.Update().SetUpsert(true)
		_, err := svc.workouts.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (svc *Store) DeleteExercise(ctx context.Context, req *loggerciseProto.ExerciseRequest) (*loggerciseProto.ExerciseResponse, error) {
	res := &loggerciseProto.ExerciseResponse{}
	fmt.Println(req.ExerciseId, req.WorkoutId)
	filter := bson.M{"_id": req.WorkoutId, "exercises._id": req.ExerciseId}
	update := bson.M{
		"$pull": bson.M{
			"exercises": bson.M{"_id": req.ExerciseId},
		},
	}
	_, err := svc.workouts.UpdateOne(ctx, filter, update)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (svc *Store) UpsertSet(ctx context.Context, req *loggerciseProto.UpsertSetRequest) (*loggerciseProto.Empty, error) {
	res := &loggerciseProto.Empty{}
	isInsert := req.SetId == ""
	set := MongoSet{
		SetId:  uuid.New().String(),
		Weight: req.Set.Weight,
		Reps:   req.Set.Reps,
		Rpe:    req.Set.Rpe,
		Rest:   req.Set.Rest,
	}
	if !isInsert {
		svc.log.Info("Update")
		set.SetId = req.SetId
		filter := bson.M{"_id": req.WorkoutId, "exercises._id": req.ExerciseId, "exercises.sets._id": set.SetId}
		update := bson.M{
			"$set": bson.M{
				"exercises.$[].sets.$[]": set,
			},
		}
		opts := options.Update().SetUpsert(true)
		_, err := svc.workouts.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return nil, err
		}
	} else {
		svc.log.Info("Insert")
		filter := bson.M{"_id": req.WorkoutId, "exercises._id": req.ExerciseId}
		update := bson.M{
			"$push": bson.M{
				"exercises.$[].sets": set,
			},
		}
		opts := options.Update().SetUpsert(true)
		_, err := svc.workouts.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return nil, err
		}

	}
	return res, nil
}

func (svc *Store) DeleteSet(ctx context.Context, req *loggerciseProto.DeleteSetRequest) (*loggerciseProto.Empty, error) {
	filter := bson.M{"_id": req.WorkoutId, "exercises._id": req.ExerciseId}
	update := bson.M{
		"$pull": bson.M{
			"exercises.$[].sets": bson.M{"_id": req.SetId},
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := svc.workouts.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return &loggerciseProto.Empty{}, nil
}

func (svc *Store) UpsertLift(ctx context.Context, req *loggerciseProto.UpsertLiftRequest) (*loggerciseProto.LiftResponse, error) {
	return nil, nil
}

func (svc *Store) GetMuscles(ctx context.Context, req *loggerciseProto.Empty) (*loggerciseProto.MusclesResponse, error) {
	res := &loggerciseProto.MusclesResponse{}
	filter := bson.M{}
	curr, err := svc.muscles.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var muscles []*loggerciseProto.Muscle
	err = mapCursorToMuscles(ctx, curr, &muscles)
	if err != nil {
		return nil, err
	}
	res.Muscles = muscles
	return res, nil
}

func (svc *Store) GetLifts(ctx context.Context, req *loggerciseProto.Empty) (*loggerciseProto.LiftResponse, error) {
	res := &loggerciseProto.LiftResponse{}
	var lifts []*loggerciseProto.Lift
	filter := bson.M{}
	cursor, err := svc.lifts.Find(ctx, filter)
	if err != nil {
		return res, err
	}
	err = mapCursorToLifts(ctx, cursor, &lifts)
	if err != nil {
		return res, err
	}
	res.Lifts = lifts
	return res, nil
}
