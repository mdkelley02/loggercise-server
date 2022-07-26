package store

import (
	loggerciseProto "loggercise/gen/go/service"
)

type MongoMuscle struct {
	Type string `bson:"_id"`
	Name string `bson:"name"`
}

type MongoLift struct {
	LiftId       string                    `bson:"_id"`
	Name         string                    `bson:"name"`
	MuscleGroups []*loggerciseProto.Muscle `bson:"muscleGroups"`
}

type MongoSet struct {
	SetId  string `bson:"_id"`
	Weight int32  `bson:"weight"`
	Reps   int32  `bson:"reps"`
	Rest   int32  `bson:"rest"`
	Rpe    int32  `bson:"rpe"`
}

type MongoExercise struct {
	ExerciseId string      `bson:"_id"`
	WorkoutId  string      `bson:"workoutId"`
	Lift       *MongoLift  `bson:"lift"`
	Sets       []*MongoSet `bson:"sets"`
}

type MongoWorkout struct {
	WorkoutId string `bson:"_id"`
	UserId    string `bson:"userId"`
	Title     string `bson:"title"`
	Notes     string `bson:"notes"`
	Date      string `bson:"date"`
	CreatedAt string `bson:"createdAt"`
	Exercises []*MongoExercise
}
