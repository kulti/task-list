import { DefaultApi } from "./openapi_cli/index"
import * as models from "./openapi_cli/model/models"

const api = new DefaultApi("http://127.0.0.1:8090/api/v1")

window.onload = load_task_lists;

class DropdownMenu {
    dropdownMenu: HTMLDivElement;
    listId: models.ListId;
    taskId: string;

    constructor(parent: Node, listId: models.ListId, taskId: string, taskState: models.RespTask.StateEnum) {
        this.listId = listId;
        this.taskId = taskId;

        this.doneTask = this.doneTask.bind(this)
        this.todoTask = this.todoTask.bind(this)
        this.cancelTask = this.cancelTask.bind(this)
        this.deleteTask = this.deleteTask.bind(this)

        this.dropdownMenu = document.createElement('div') as HTMLDivElement;
        this.dropdownMenu.className = "dropdown-menu";

        this.appendItem("Done", this.doneTask)

        if (this.listId == models.ListId.Sprint) {
            const action = this.appendItem("Todo", this.todoTask)
            if (taskState == models.RespTask.StateEnum.Todo) {
                action.className += " disabled"
            }
        }

        this.appendItem("Cancel", this.cancelTask)
        this.appendItem("Delete", this.deleteTask)

        parent.appendChild(this.dropdownMenu)
    }

    doneTask() {
        api.doneTask(this.taskId).done(function () {
            load_task_lists();
        }).fail(function () {
            showErrorAlert("failed to done task")
        })
    }

    todoTask() {
        api.takeTask(models.ListId.Todo, this.taskId).done(function () {
            load_task_lists();
        }).fail(function () {
            showErrorAlert("failed to todo task")
        })
    }

    cancelTask() {
        api.cancelTask(this.taskId).done(function () {
            load_task_lists();
        }).fail(function () {
            showErrorAlert("failed to cancel task")
        })
    }

    deleteTask() {
        api.deleteTask(this.listId, this.taskId).done(function () {
            load_task_lists();
        }).fail(function () {
            showErrorAlert("failed to delete task")
        })
    }

    appendItem(text: string, handler: EventListener): HTMLElement {
        const action = document.createElement('div') as HTMLDivElement;
        action.className = "dropdown-item"
        action.innerText = text
        action.addEventListener("click", handler)
        this.dropdownMenu.appendChild(action)
        return action
    }
}

$("#new_sprint_btn")[0].addEventListener("click", function () {
    const input = <HTMLInputElement>$("#new_sprint_title")[0]
    const sprintTitle = input.value || input.placeholder
    const sprintOpts: models.SprintOpts = {
        title: sprintTitle
    }
    api.createTaskList(sprintOpts).done(function () {
        showSuccessAlert("sprint created")
        load_task_lists();
    }).fail(function (body) {
        showErrorAlert("failed to create sprint")
    });
});

function showSuccessAlert(text: string) {
    showAlert("success", text)
}

function showErrorAlert(text: string) {
    showAlert("danger", text)
}

function showAlert(type: string, text: string) {
    const alertCloseBtn = document.createElement('button') as HTMLButtonElement
    alertCloseBtn.type = "button"
    alertCloseBtn.className = "close"
    alertCloseBtn.setAttribute("data-dismiss", "alert")
    alertCloseBtn.setAttribute("aria-label", "Close")
    alertCloseBtn.innerHTML = '<span aria-hidden="true">&times;</span>'

    const alertDiv = document.createElement('div') as HTMLDivElement;
    alertDiv.className = "alert alert-dismissible fade show alert-" + type;
    alertDiv.setAttribute("role", "alert")
    alertDiv.innerText = text
    alertDiv.appendChild(alertCloseBtn)

    $('#alerts').append(alertDiv)
    setTimeout(() => {
        alertDiv.remove()
    }, 2000);
}

function load_task_lists() {
    api.getTaskList(models.ListId.Sprint).fail(function (body) {
        showErrorAlert("failed to load sprint tasks")
    }).done(function (data) {
        const sprint_task_list = data.body
        api.getTaskList(models.ListId.Todo).fail(function (body) {
            showErrorAlert("failed to load todo tasks")
        }).done(function (data) {
            const todo_task_list = data.body
            draw_task_lists(todo_task_list, sprint_task_list);
        });
    });
}

enum TaskProperty {
    Completed = 'completed',
    Canceled = 'canceled',
}

function draw_task_lists(todo_task_list: models.TaskList, sprint_task_list: models.TaskList) {
    prepare_task_list(todo_task_list)
    prepare_task_list(sprint_task_list)

    update_task_list_header(todo_task_list, "#todo_list")
    update_task_list_header(sprint_task_list, "#sprint_list")

    fill_task_list(models.ListId.Todo, todo_task_list.tasks, "#todo_list")
    fill_task_list(models.ListId.Sprint, sprint_task_list.tasks, "#sprint_list")
}

function update_task_list_header(task_list: models.TaskList, id: string) {
    const points = sum_points(task_list.tasks)
    const burnt = sum_burnt_points(task_list.tasks)

    const task_list_html = $(id + " .list_header")[0]
    task_list_html.getElementsByClassName("title")[0].innerHTML = task_list.title
    task_list_html.getElementsByClassName("points")[0].innerHTML = burnt + "/" + points
}

function fill_task_list(listId: models.ListId, tasks: Array<models.RespTask>, id: string) {
    const task_list_html = $(id + " .tasks")[0]

    task_list_html.innerHTML = ""
    tasks.forEach(function (task) {
        task_list_html.append(build_task_html(listId, task))
    })

    task_list_html.append(build_new_task_input_html(listId))
}

function build_task_html(listId: models.ListId, task: models.RespTask): HTMLElement {
    let points = ''
    let percent = 0
    if (task.state != models.Task.StateEnum.Canceled) {
        points = task.burnt + "/" + task.points
        percent = 100 * task.burnt / task.points
    }

    let task_properties = ""
    if (task.state == models.RespTask.StateEnum.Done) {
        task_properties = TaskProperty.Completed
    } else if (task.state == models.RespTask.StateEnum.Canceled) {
        task_properties = TaskProperty.Canceled
    }

    const taskIdDiv = document.createElement('div') as HTMLDivElement;
    taskIdDiv.className = "task__id";
    taskIdDiv.innerText = task.id;

    const taskTextDiv = document.createElement('div') as HTMLDivElement;
    taskTextDiv.className = "text";
    taskTextDiv.innerText = task.text;

    const taskPointsDiv = document.createElement('div') as HTMLDivElement;
    taskPointsDiv.className = "points";
    taskPointsDiv.innerText = points;
    if (percent > 0) {
        taskPointsDiv.style.background = "-webkit-linear-gradient(left, #f8f8f8 " + percent + "%, white " + percent + "%)";
    }

    const taskDiv = document.createElement('div') as HTMLDivElement;
    taskDiv.className = "task " + task_properties;
    taskDiv.setAttribute("type", "button");
    taskDiv.setAttribute("data-toggle", "dropdown");
    taskDiv.append(taskIdDiv, taskTextDiv, taskPointsDiv);

    const dropdown = document.createElement('div') as HTMLDivElement;
    dropdown.className = "dropdown show";
    dropdown.append(taskDiv);

    taskDiv.onclick = function () {
        new DropdownMenu(dropdown, listId, task.id, task.state)
    };

    dropdown.ondblclick = function (): any {
        dropdown.replaceWith(build_task_input_html(listId, task))
        return false
    }

    return dropdown;
}


function build_task_input_html(listId: models.ListId, task: models.RespTask): HTMLElement {
    const taskTextInput = document.createElement('input') as HTMLInputElement;
    taskTextInput.className = "text form-control";
    taskTextInput.type = "text";
    taskTextInput.value = task.text;

    const taskPointsInput = document.createElement('input') as HTMLInputElement;
    taskPointsInput.className = "points form-control";
    taskPointsInput.type = "text";
    taskPointsInput.value = task.burnt + "/" + task.points;

    const autofocusPoints = ($(".text:hover").length === 0)
    const autofocusInput = autofocusPoints ? taskPointsInput : taskTextInput;
    setTimeout(function () {
        autofocusInput.focus()
    }, 0);

    const taskDiv = document.createElement('div') as HTMLDivElement;
    taskDiv.className = "form-group task";
    taskDiv.append(taskTextInput, taskPointsInput);

    const handleKeyPress = function (ev: KeyboardEvent) {
        switch (ev.keyCode) {
            case 27:
                taskDiv.replaceWith(build_task_html(listId, task));
                break;
            case 13:
                const points = taskPointsInput.value.split("/")
                const opts: models.UpdateOptions = {
                    text: taskTextInput.value,
                    points: parseInt(points[1]),
                    burnt: parseInt(points[0]),
                }
                api.updateTask(task.id, opts).done(function () {
                    load_task_lists();
                }).fail(function (body) {
                    showErrorAlert("failed to update task")
                });
                load_task_lists();
        }
    }

    taskTextInput.onkeyup = handleKeyPress;
    taskPointsInput.onkeyup = handleKeyPress;

    return taskDiv;
}

function build_new_task_input_html(listId: models.ListId): HTMLElement {
    const taskTextInput = document.createElement('input') as HTMLInputElement;
    taskTextInput.className = "text form-control";
    taskTextInput.type = "text";
    taskTextInput.placeholder = "Do new task";

    const taskPointsInput = document.createElement('input') as HTMLInputElement;
    taskPointsInput.className = "points form-control";
    taskPointsInput.type = "text";
    taskPointsInput.placeholder = "0";

    const taskDiv = document.createElement('div') as HTMLDivElement;
    taskDiv.className = "form-group task";
    taskDiv.append(taskTextInput, taskPointsInput);


    const handleKeyPress = function (ev: KeyboardEvent) {
        switch (ev.keyCode) {
            case 27:
                taskDiv.replaceWith(build_new_task_input_html(listId));
                break;
            case 13:
                const task: models.Task = {
                    text: taskTextInput.value,
                    points: parseInt(taskPointsInput.value),
                }
                api.createTask(listId, task).done(function () {
                    load_task_lists();
                }).fail(function (body) {
                    showErrorAlert("failed to create task")
                });
        }
    }
    taskTextInput.onkeyup = handleKeyPress;
    taskPointsInput.onkeyup = handleKeyPress;

    return taskDiv;
}

function sum_points(tasks: Array<models.RespTask>): number {
    return tasks.reduce((sum, current) => {
        return sum + current.points;
    }, 0)
}

function sum_burnt_points(tasks: Array<models.RespTask>): number {
    return tasks.reduce((sum, current) => {
        return sum + current.burnt;
    }, 0)
}

function prepare_task_list(task_list: models.TaskList): void {
    const fix_points = (value: models.RespTask) => {
        switch (value.state) {
            case models.RespTask.StateEnum.Done:
                value.burnt = value.points
                break

            case models.RespTask.StateEnum.Canceled:
                value.points = 0
                value.burnt = 0
                break

            default:
                if (!value.burnt) {
                    value.burnt = 0
                }
                break
        }
    }

    task_list.tasks.forEach(fix_points)
}
