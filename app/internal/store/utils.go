package store

import (
	"context"

	loggerciseProto "loggercise/gen/go/service"

	"go.mongodb.org/mongo-driver/mongo"
)

func mapResultToWorkout(ctx context.Context, result *mongo.SingleResult, dest *loggerciseProto.Workout) error {
	var workout *MongoWorkout
	err := result.Decode(&workout)
	if err != nil {
		return err
	}
	mapMongoWorkoutToProtoWorkout(workout, dest)
	return nil
}

func mapMongoWorkoutToProtoWorkout(mongo *MongoWorkout, dest *loggerciseProto.Workout) {
	*dest = loggerciseProto.Workout{
		WorkoutId: mongo.WorkoutId,
		Title:     mongo.Title,
		Notes:     mongo.Notes,
		CreatedAt: mongo.CreatedAt,
		Date:      mongo.Date,
		Exercises: []*loggerciseProto.Exercise{},
		UserId:    mongo.UserId,
	}
	for _, exercise := range mongo.Exercises {
		pbExercise := &loggerciseProto.Exercise{
			WorkoutId:  exercise.WorkoutId,
			ExerciseId: exercise.ExerciseId,
			Lift: &loggerciseProto.Lift{
				LiftId:       exercise.Lift.LiftId,
				Name:         exercise.Lift.Name,
				MuscleGroups: exercise.Lift.MuscleGroups,
			},
			Sets: []*loggerciseProto.ExerciseSet{},
		}
		for _, set := range exercise.Sets {
			pbSet := &loggerciseProto.ExerciseSet{
				SetId:  set.SetId,
				Weight: set.Weight,
				Reps:   set.Reps,
				Rest:   set.Rest,
				Rpe:    set.Rpe,
			}
			pbExercise.Sets = append(pbExercise.Sets, pbSet)
		}
		dest.Exercises = append(dest.Exercises, pbExercise)
	}
}

func mapCursorToWorkouts(ctx context.Context, curr *mongo.Cursor, dest *[]*loggerciseProto.Workout) error {
	for curr.Next(ctx) {
		var workout MongoWorkout
		err := curr.Decode(&workout)
		if err != nil {
			return err
		}
		var destWorkout loggerciseProto.Workout
		mapMongoWorkoutToProtoWorkout(&workout, &destWorkout)
		*dest = append(*dest, &destWorkout)
	}
	return nil
}

func mapCursorToLifts(ctx context.Context, cursor *mongo.Cursor, dest *[]*loggerciseProto.Lift) error {
	for cursor.Next(ctx) {
		var lift MongoLift
		err := cursor.Decode(&lift)
		if err != nil {
			return err
		}
		lift.LiftId = lift.LiftId[:len(lift.LiftId)-1]
		*dest = append(*dest, &loggerciseProto.Lift{
			LiftId:       lift.LiftId,
			Name:         lift.Name,
			MuscleGroups: lift.MuscleGroups,
		})
	}
	return nil
}

func mapCursorToMuscles(ctx context.Context, cursor *mongo.Cursor, dest *[]*loggerciseProto.Muscle) error {
	for cursor.Next(ctx) {
		var muscle MongoMuscle
		err := cursor.Decode(&muscle)
		if err != nil {
			return err
		}
		*dest = append(*dest, &loggerciseProto.Muscle{
			Type: muscle.Type,
			Name: muscle.Name,
		})
	}
	return nil
}
