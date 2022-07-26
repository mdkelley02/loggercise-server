package store

import (
	"context"

	loggerciseProto "loggercise/gen/go/service"

	"go.mongodb.org/mongo-driver/mongo"
)

func mapCursorToWorkouts(ctx context.Context, curr *mongo.Cursor, dest *[]*loggerciseProto.Workout) error {
	for curr.Next(ctx) {
		var workout loggerciseProto.Workout
		err := curr.Decode(&workout)
		if err != nil {
			return err
		}
		workout.WorkoutId = curr.Current.Lookup("_id").StringValue()
		*dest = append(*dest, &workout)
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
