package routes

import (
	"backend/controllers"
	"backend/models"
	"backend/tools"
	"backend/tools/db"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/Jeffail/gabs"
)

// Returns the exercise list.
// Returns 401 if the user is not logged in.
// Returns 500 HTTP Status Code if tan error occurs when
// querying the exercises.
func GetExerciseList(c *fiber.Ctx) error {
	dbase := db.Orm()
	logged, _ := tools.GetCurrentUserId(c)

	if !logged {
		return controllers.ApiError(c, "Login required", 401)
	}

	var exercises *[]models.Exercise

	exercises, err := controllers.GetExercises(dbase)

	if err != nil {
		return controllers.ApiError(c, "Failed to retrieve exercises", 500)
	}

	return c.JSON(fiber.Map{
		"exercises": exercises,
	})
}

// Returns the Tag list.
// Returns 401 if the user is not logged in.
// Returns 500 HTTP Status Code if tan error occurs when
// querying the tags.
func GetTagList(c *fiber.Ctx) error {
	dbase := db.Orm()
	logged, _ := tools.GetCurrentUserId(c)

	if !logged {
		return controllers.ApiError(c, "Login required", 401)
	}

	var tags *[]models.Tag

	tags, err := controllers.GetTags(dbase)

	if err != nil {
		return controllers.ApiError(c, "Failed to retrieve exercises", 500)
	}

	return c.JSON(fiber.Map{
		"tags": tags,
	})
}

// Returns the workout list of a trainer.
// Returns 401 if the user is not logged in.
// Returns 500 HTTP Status Code if tan error occurs when
// querying the workouts.
func GetWorkoutList(c *fiber.Ctx) error {
	dbase := db.Orm()
	logged, id := tools.GetCurrentUserId(c)

	if !logged {
		return controllers.ApiError(c, "Login required", 401)
	}

	idTrainer := uint64(id)
 
	type wk_tag struct {
		Workout models.Workout 
		Tags 	[]models.Tag   
	}

	var workouts_tag *[]wk_tag
	workouts_tag = new([]wk_tag)
	var workouts *[]models.Workout
	var wk_tags *[]models.WorkoutTag
	var tags *[]models.Tag
	tags = new([]models.Tag)

	workouts, err := controllers.GetWorkoutsByTrainer(dbase, idTrainer)
	if err != nil {
		return controllers.ApiError(c, "Failed to retrieve trainer workouts", 500)
	}

	for _, w := range *workouts {
		wk_tags, err = controllers.GetWorkoutTagByWkId(dbase, uint64(w.ID))
		if err != nil {
			return controllers.ApiError(c, "Failed to retrieve workout tags", 500)
		}
		for _, wt := range *wk_tags {
			tag, err := controllers.GetTagById(dbase, uint64(wt.IdTag))
			if err != nil {
				return controllers.ApiError(c, "Failed to retrieve tags", 500)
			}
			*tags = append(*tags, *tag)
		}
		workout_with_tag := wk_tag{}
		workout_with_tag.Workout = w
		workout_with_tag.Tags = *tags

		*workouts_tag = append(*workouts_tag, workout_with_tag)
	}
	log.Print(*tags)
	log.Print(*workouts_tag)

	return c.JSON(fiber.Map{
		"workouts": *workouts_tag,
	})
}

// Returns the detailed workout (exercises sets and reps).
// Returns 401 if the user is not logged in.
// Returns 500 HTTP Status Code if tan error occurs when
// querying the workouts.
func GetWorkoutDetail(c *fiber.Ctx) error {
	dbase := db.Orm()
	logged, _ := tools.GetCurrentUserId(c)

	if !logged {
		return controllers.ApiError(c, "Login required", 401)
	}

	idWk := c.Query("idWorkout", "")

	if idWk == "" {
		log.Print("The workout id must not be empty")
		return controllers.ApiError(c, "idWorkout must not be empty", 404)
	}

	idWorkout, err := strconv.ParseUint(idWk, 10, 64)

	if err != nil {
		log.Print("The workout id could not be parsed to int")
		log.Print(err)
		return controllers.ApiError(c, "idWorkout must be an integer", 400)
	}

	type wk_exercises struct {
		IdExercise  uint64
		Name        string
		Description string
		Sets        uint64
		Reps        uint64
		Ordinal     uint64
	}

	var wk_e []wk_exercises

	result := dbase.Order("ordinal asc").Table("workout_exercises").Select("workout_exercises.id_exercise "+
		"as id_exercise, "+
		"exercises.name as name, "+
		"exercises.description as description, "+
		"workout_exercises.sets as sets, "+
		"workout_exercises.reps as reps, "+
		"workout_exercises.ordinal as ordinal").Joins("join exercises on "+
		"exercises.id = "+
		"workout_exercises.id_exercise").Where("id_workout = ?",
		idWorkout).Find(&wk_e)

	if result.Error != nil {
		return controllers.ApiError(c, "Failed to retrieve workout exercises", 500)
	}

	var wk_tags *[]models.WorkoutTag
	var tags *[]models.Tag
	tags = new([]models.Tag)

	wk_tags, err = controllers.GetWorkoutTagByWkId(dbase, idWorkout)

	for _, wt := range *wk_tags {
		tag, err := controllers.GetTagById(dbase, uint64(wt.IdTag))
		if err != nil {
			return controllers.ApiError(c, "Failed to retrieve tags", 500)
		}
		*tags = append(*tags, *tag)
	}

	return c.JSON(fiber.Map{
		"exercises": wk_e,
		"Tags": *tags,
	})
}

// Updates or creates the workout and returns the update.
// (exercises sets and reps).
// Returns 201 when the workout is created and 200 when it is
// updated.
// Returns 401 if the user is not logged in.
// Returns 500 HTTP Status Code if tan error occurs when
// querying the workouts.
func UpdateCreateWorkout(c *fiber.Ctx) error {
	dbase := db.Orm()
	logged, id := tools.GetCurrentUserId(c)

	if !logged {
		return controllers.ApiError(c, "Login required", 401)
	}

	idTrainer := uint64(id)

	body, err := gabs.ParseJSON([]byte(c.Body()))

	if err != nil {
		return controllers.ApiError(c, "Something went wrong when "+
									"parsing request body", 500)
	}

	name, exists_name := body.Path("Name").Data().(string)
		if !exists_name {
			return controllers.ApiError(c, "Name must not be empty", 404)
		}

	idWk, exists_id := body.Path("Id").Data().(uint)

	if err != nil {
		log.Print("The workout id could not be parsed to int")
		log.Print(err)
		return controllers.ApiError(c, "idWorkout must be an integer", 400)
	}
	
	newWk := models.NewWorkout(idTrainer, name)

	if !exists_id {
		log.Print("Creating new workout")

		controllers.CreateWorkout(dbase, &newWk)
	} else {
		log.Print("Updating workout")

		newWk.ID = idWk
		controllers.UpdateWorkout(dbase, &newWk)
	}

	exercises, err := body.Search("Exercises").Children()

	if err != nil {
		log.Print("The exercise array is empty")
		return c.JSON(fiber.Map{
			"message": "The workout exercises could not be updated",
		})
	}

	_idWk := uint64(idWk)
	controllers.DeleteWorkoutExercises(dbase, _idWk)

	for _, exercise := range exercises {
		idE, exists := exercise.Path("id").Data().(uint64)
		if !exists {continue}
		order, exists := exercise.Path("order").Data().(uint8)
		if !exists {continue}
		reps, exists := exercise.Path("reps").Data().(uint16)
		if !exists {continue}
		sets, exists := exercise.Path("sets").Data().(uint8)
		if !exists {continue}

		newWkEx := models.NewWorkoutExercise(idE, _idWk, order, sets, reps)
		controllers.UpdateWorkoutExercise(dbase, &newWkEx)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
