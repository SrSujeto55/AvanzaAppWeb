import axios from "axios"
import { excerciseOptions } from "./Exercise"

const fetchExercises = async () => {
    const response = await axios.get("api/exercises")
    const data = response.data
    const exercises = data.exercises.map( (exercise: { ID: number, Name: string}) => ({ ID: exercise.ID, value: exercise.Name, label: exercise.Name }))
    return exercises
}

const filterTags = async (inputValue: string) => {
    const tagsOptions = await fetchExercises() 

    return tagsOptions.filter((i : excerciseOptions) =>
      i.label.toLowerCase().includes(inputValue.toLowerCase())
    );
};
  
const ExercisePromiseOptions = (inputValue: string) =>
    new Promise<excerciseOptions[]>((resolve) => {
    setTimeout(() => {
        resolve(filterTags(inputValue));
    }, 1000);
});

export default ExercisePromiseOptions