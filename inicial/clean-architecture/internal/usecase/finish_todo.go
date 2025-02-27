package usecase

// DTO (Data Transfer obeject)
type InputFinishTodo struct {
	ID string
}

type OutputFinishTodo struct {
	ID string
}

type InputCompensateFinishTodo struct {
	ID     string
	Reason string
}

type FinishTodoUseCase struct {
	TodoRepository  TodoRepository  // mundo externo
	CompensateEvent CompensateEvent // mundo externo / interno
}

func (f *FinishTodoUseCase) Execute(input any) (any, error) {
	inputFinishTodo := input.(InputFinishTodo)
	todo, err := f.TodoRepository.FindByID(inputFinishTodo.ID)
	if err != nil {
		return nil, err
	}
	todo.Done()
	err = f.TodoRepository.Save(todo)
	if err != nil {
		return nil, err
	}

	return OutputFinishTodo{
		ID: todo.ID,
	}, nil
}

func (f *FinishTodoUseCase) Compensate(input any) (any, error) {
	inputCompensateFinishTodo := input.(InputCompensateFinishTodo)
	todo, err := f.TodoRepository.FindByID(inputCompensateFinishTodo.ID)
	if err != nil {
		return nil, err
	}
	todo.Undone()
	err = f.TodoRepository.Save(todo)
	if err != nil {
		return nil, err
	}

	err = f.CompensateEvent.Publish("todo_undone", todo)
	if err != nil {
		return err, nil
	}

	return OutputFinishTodo{
		ID: todo.ID,
	}, nil

}
