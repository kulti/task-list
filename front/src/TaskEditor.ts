export enum TaskEditorFocus {
  Text = "text",
  Points = "points",
  None = "none",
}

export interface TaskEditorTask {
  text: string;
  points: number;
  burnt?: number;
}

export function BuildTaskEditor(
  applyFn: (text: string, points: string) => void,
  cancelFn?: () => void,
  resetDiv?: HTMLElement,
  task?: TaskEditorTask,
  focus = TaskEditorFocus.None
): HTMLElement {
  const taskTextInput = document.createElement("input") as HTMLInputElement;
  taskTextInput.className = "text form-control";
  taskTextInput.type = "text";
  taskTextInput.placeholder = "Do new task";

  const taskPointsInput = document.createElement("input") as HTMLInputElement;
  taskPointsInput.className = "points form-control";
  taskPointsInput.type = "text";
  taskPointsInput.placeholder = "0";

  if (task) {
    taskTextInput.value = task.text;
    if (task.burnt !== undefined) {
      taskPointsInput.value = task.burnt + "/";
    }
    taskPointsInput.value += task.points;
  }

  if (focus !== TaskEditorFocus.None) {
    const autofocusPoints = focus === TaskEditorFocus.Points;
    const autofocusInput = autofocusPoints ? taskPointsInput : taskTextInput;
    setTimeout(() => {
      autofocusInput.focus();
    }, 0);
  }

  const taskDiv = document.createElement("div") as HTMLDivElement;
  taskDiv.className = "form-group task";
  taskDiv.append(taskTextInput, taskPointsInput);

  const handleKeyPress = (ev: KeyboardEvent) => {
    switch (ev.key) {
      case "Escape":
        if (resetDiv) {
          taskDiv.replaceWith(resetDiv);
        } else {
          taskTextInput.value = "";
          taskPointsInput.value = "";
        }

        if (cancelFn) {
          cancelFn();
        }
        break;
      case "Enter":
        applyFn(taskTextInput.value, taskPointsInput.value);
    }
  };

  taskTextInput.onkeyup = handleKeyPress;
  taskPointsInput.onkeyup = handleKeyPress;

  return taskDiv;
}
