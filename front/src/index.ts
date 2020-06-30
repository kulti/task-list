import { DefaultApi } from "./openapi_cli/index"
import * as models from "./openapi_cli/model/models"
import { BuildDropdownMenu } from "./DropdownMenu"
import { BuildTaskEditor, TaskEditorFocus, TaskEditorTask } from "./TaskEditor"

const api = new DefaultApi(window.location.origin + "/api/v1")

let sprintTemplate: models.SprintTemplate;

window.onload = () => {
    const socket = new WebSocket("ws://" + window.location.host + "/ws");
    socket.onclose = () => {
        showErrorAlertWithRefresh("lost connection to server")
    };

    const input = $("#new_sprint_title")[0] as HTMLInputElement
    input.placeholder = buildNewSprintTitle();
    setInterval(() => {
        input.placeholder = buildNewSprintTitle()
    }, 60 * 60 * 1000)

    load_task_lists();
}

function buildNewSprintTitle(): string {
    const numToString = (v: number): string => {
        let s = v.toString()
        if (v < 10) {
            s = "0" + s
        }
        return s;
    }

    const dateToString = (d: Date): string => {
        return numToString(d.getDate()) + "." + numToString(d.getMonth() + 1);
    }

    const date = new Date()
    const dayOfWeek = date.getDay()

    date.setDate(date.getDate() + 1 - dayOfWeek)
    const beginDate = dateToString(date);

    date.setDate(date.getDate() + 6)
    const endDate = dateToString(date);

    return beginDate + " - " + endDate;
}

$("#new_sprint_btn")[0].addEventListener("click", () => {
    const input = $("#new_sprint_title")[0] as HTMLInputElement
    const sprintTitle = input.value || input.placeholder
    const sprintOpts: models.SprintOpts = {
        title: sprintTitle
    }
    api.createTaskList(sprintOpts).done(() => {
        api.getSprintTemplate().done((data) => {
            sprintTemplate = data.body
            load_task_lists();
        }).fail(() => {
            load_task_lists();
        });
        showSuccessAlert("sprint created")
    }).fail((body) => {
        showErrorAlert("failed to create sprint")
    });
});

function showSuccessAlert(text: string) {
    showAlert("success", text)
}

function showErrorAlert(text: string) {
    showAlert("danger", text)
}

function showErrorAlertWithRefresh(text: string) {
    showAlert("danger", text, 5)
}

function showAlert(type: string, text: string, refreshSec?: number) {
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

    if (refreshSec) {
        alertDiv.innerText = text + " - refresh in " + refreshSec + "seconds";
        let i = refreshSec - 1;
        setInterval(() => {
            alertDiv.innerText = text + " - refresh in " + i + " seconds";
            i--;
            if (i === 0) {
                window.location.reload();
            }
        }, 1000);
    } else {
        setTimeout(() => {
            alertDiv.remove()
        }, 2000);
    }
}

function load_task_lists() {
    api.getTaskList(models.ListId.Sprint).fail((body) => {
        showErrorAlertWithRefresh("failed to load sprint tasks")
    }).done((data) => {
        draw_task_lists(data.body);
    });
}

enum TaskProperty {
    Completed = 'completed',
    Canceled = 'canceled',
    Todo = 'todo',
}

function draw_task_lists(sprintTaskList: models.TaskList) {
    prepare_task_list(sprintTaskList)

    update_task_list_header(sprintTaskList, models.ListId.Sprint)

    fill_task_list(models.ListId.Sprint, sprintTaskList.tasks)
}

function update_task_list_header(taskList: models.TaskList, listId: models.ListId) {
    const points = sum_points(taskList.tasks)
    const burnt = sum_burnt_points(taskList.tasks)

    const taskListHtml = $(listHtmlId(listId) + " .list_header")[0]
    taskListHtml.getElementsByClassName("title")[0].innerHTML = taskList.title
    taskListHtml.getElementsByClassName("points")[0].innerHTML = burnt + "/" + points
}

function fill_task_list(listId: models.ListId, tasks: models.RespTask[]) {
    const taskListHtml = $(listHtmlId(listId) + " .tasks")[0]

    taskListHtml.innerHTML = ""
    tasks.forEach((task) => {
        taskListHtml.append(build_task_html(listId, task))
    })

    if (sprintTemplate && sprintTemplate.tasks && sprintTemplate.tasks.length > 0) {
        const task: TaskEditorTask = {
            text: sprintTemplate.tasks[0].text,
            points: sprintTemplate.tasks[0].points,
        }
        sprintTemplate.tasks.splice(0, 1)
        taskListHtml.append(build_template_task_input_html(models.ListId.Sprint, task))
        focus_new_task_input(listId);
    } else {
        taskListHtml.append(build_new_task_input_html(listId))
    }
}

function build_task_html(listId: models.ListId, task: models.RespTask): HTMLElement {
    let points = ''
    let percent = 0
    if (task.state !== models.RespTask.StateEnum.Canceled) {
        points = task.burnt + "/" + task.points
        percent = 100 * task.burnt / task.points
    }

    let taskProperties = ""
    if (task.state === models.RespTask.StateEnum.Done) {
        taskProperties = TaskProperty.Completed
    } else if (task.state === models.RespTask.StateEnum.Canceled) {
        taskProperties = TaskProperty.Canceled
    } else if (task.state === models.RespTask.StateEnum.Todo) {
        taskProperties = TaskProperty.Todo
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
    taskDiv.className = "task " + taskProperties;
    taskDiv.setAttribute("type", "button");
    taskDiv.setAttribute("data-toggle", "dropdown");
    taskDiv.append(taskIdDiv, taskTextDiv, taskPointsDiv);

    const dropdown = document.createElement('div') as HTMLDivElement;
    dropdown.className = "dropdown show";
    dropdown.append(taskDiv);

    taskDiv.onclick = () => {
        dropdown.append(build_dropdown_menu(listId, task))
    };

    dropdown.ondblclick = (): any => {
        dropdown.replaceWith(build_task_input_html(task, dropdown))
        return false
    }

    return dropdown;
}

function build_dropdown_menu(listId: models.ListId, task: models.RespTask): HTMLDivElement {
    return BuildDropdownMenu(task.state,
        () => {
            api.todoTask(task.id).done(() => {
                load_task_lists();
            }).fail(() => {
                showErrorAlert("failed to todo task")
            })
        },
        () => {
            api.doneTask(task.id).done(() => {
                load_task_lists();
            }).fail(() => {
                showErrorAlert("failed to done task")
            })
        },
        () => {
            api.cancelTask(task.id).done(() => {
                load_task_lists();
            }).fail(() => {
                showErrorAlert("failed to cancel task")
            })
        },
        () => {
            api.deleteTask(listId, task.id).done(() => {
                load_task_lists();
            }).fail(() => {
                showErrorAlert("failed to delete task")
            })
        }
    )

}

function build_task_input_html(task: models.RespTask, resetDiv: HTMLElement): HTMLElement {
    const editorTask: TaskEditorTask = {
        text: task.text,
        points: task.points,
        burnt: task.burnt,
    }

    const autofocusPoints = ($(".text:hover").length === 0)
    const focus = autofocusPoints ? TaskEditorFocus.Points : TaskEditorFocus.Text;

    return BuildTaskEditor((text: string, points: string) => {
        const pointsArr = points.split("/")
        const opts: models.UpdateOptions = {
            text,
            points: parseInt(pointsArr[1], 10),
            burnt: parseInt(pointsArr[0], 10),
        }
        api.updateTask(task.id, opts).done(() => {
            load_task_lists();
        }).fail((body) => {
            showErrorAlert("failed to update task")
        });
        load_task_lists();
    }, undefined, resetDiv, editorTask, focus);
}

function build_new_task_input_html(listId: models.ListId): HTMLElement {
    return BuildTaskEditor((text: string, points: string) => {
        const newTask: models.Task = {
            text,
            points: parseInt(points, 10),
        }
        api.createTask(listId, newTask).done(() => {
            load_task_lists();
            focus_new_task_input(listId);
        }).fail((body) => {
            showErrorAlert("failed to create task")
        });
    });
}

function build_template_task_input_html(listId: models.ListId, task: TaskEditorTask): HTMLElement {
    return BuildTaskEditor((text: string, points: string) => {
        const newTask: models.Task = {
            text,
            points: parseInt(points, 10),
        }
        api.createTask(listId, newTask).done(() => {
            load_task_lists();
            focus_new_task_input(listId);
        }).fail((body) => {
            showErrorAlert("failed to create task")
        });
    }, () => {
        load_task_lists();
        focus_new_task_input(listId);
    }, undefined, task);
}

function focus_new_task_input(listId: models.ListId) {
    setTimeout(() => {
        $(listHtmlId(listId) + " .text.form-control")[0].focus()
    }, 100);
}

function sum_points(tasks: models.RespTask[]): number {
    return tasks.reduce((sum, current) => {
        return sum + current.points;
    }, 0)
}

function sum_burnt_points(tasks: models.RespTask[]): number {
    return tasks.reduce((sum, current) => {
        return sum + current.burnt;
    }, 0)
}

function prepare_task_list(taskList: models.TaskList): void {
    const fixPoints = (value: models.RespTask) => {
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

    taskList.tasks.forEach(fixPoints)
}

function listHtmlId(listId: models.ListId): string {
    return "#" + listId + "_list";
}
