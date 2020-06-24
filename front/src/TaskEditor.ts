import * as models from "./openapi_cli/model/models"

export function BuildTaskEditor(applyFn: (text: string, points: string) => void, resetDiv?: HTMLElement, task?: models.RespTask): HTMLElement {
    const taskTextInput = document.createElement('input') as HTMLInputElement;
    taskTextInput.className = "text form-control";
    taskTextInput.type = "text";

    const taskPointsInput = document.createElement('input') as HTMLInputElement;
    taskPointsInput.className = "points form-control";
    taskPointsInput.type = "text";

    if (task) {
        taskTextInput.value = task.text;
        taskPointsInput.value = task.burnt + "/" + task.points;

        const autofocusPoints = ($(".text:hover").length === 0)
        const autofocusInput = autofocusPoints ? taskPointsInput : taskTextInput;
        setTimeout(() => {
            autofocusInput.focus()
        }, 0);
    } else {
        taskTextInput.placeholder = "Do new task";
        taskPointsInput.placeholder = "0";
    }

    const taskDiv = document.createElement('div') as HTMLDivElement;
    taskDiv.className = "form-group task";
    taskDiv.append(taskTextInput, taskPointsInput);

    const handleKeyPress = (ev: KeyboardEvent) => {
        switch (ev.keyCode) {
            case 27:
                if (resetDiv) {
                    taskDiv.replaceWith(resetDiv)
                } else {
                    taskTextInput.value = ""
                    taskPointsInput.value = ""
                }
                break;
            case 13:
                applyFn(taskTextInput.value, taskPointsInput.value)
        }
    }

    taskTextInput.onkeyup = handleKeyPress;
    taskPointsInput.onkeyup = handleKeyPress;

    return taskDiv;
}
