import * as models from "./openapi_cli/model/models"

export enum TaskEditoFocus {
    Text = 'text',
    Points = 'points',
    None = 'none'
}

export function BuildTaskEditor(
    applyFn: (text: string, points: string) => void,
    resetDiv?: HTMLElement,
    task?: models.RespTask,
    focus = TaskEditoFocus.None
): HTMLElement {
    const taskTextInput = document.createElement('input') as HTMLInputElement;
    taskTextInput.className = "text form-control";
    taskTextInput.type = "text";
    taskTextInput.placeholder = "Do new task";

    const taskPointsInput = document.createElement('input') as HTMLInputElement;
    taskPointsInput.className = "points form-control";
    taskPointsInput.type = "text";
    taskPointsInput.placeholder = "0";

    if (task) {
        taskTextInput.value = task.text;
        taskPointsInput.value = task.burnt + "/" + task.points;
    }

    if (focus !== TaskEditoFocus.None) {
        const autofocusPoints = (focus === TaskEditoFocus.Points)
        const autofocusInput = autofocusPoints ? taskPointsInput : taskTextInput;
        setTimeout(() => {
            autofocusInput.focus()
        }, 0);
    }

    const taskDiv = document.createElement('div') as HTMLDivElement;
    taskDiv.className = "form-group task";
    taskDiv.append(taskTextInput, taskPointsInput);

    const handleKeyPress = (ev: KeyboardEvent) => {
        switch (ev.key) {
            case 'Escape':
                if (resetDiv) {
                    taskDiv.replaceWith(resetDiv)
                } else {
                    taskTextInput.value = ""
                    taskPointsInput.value = ""
                }
                break;
            case 'Enter':
                applyFn(taskTextInput.value, taskPointsInput.value)
        }
    }

    taskTextInput.onkeyup = handleKeyPress;
    taskPointsInput.onkeyup = handleKeyPress;

    return taskDiv;
}
