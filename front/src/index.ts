import { DefaultApi } from "./openapi_cli/index";
import * as models from "./openapi_cli/model/models";
import { BuildDropdownMenu } from "./DropdownMenu";
import { BuildTaskEditor, TaskEditorFocus, TaskEditorTask } from "./TaskEditor";
import { buildNewSprintTitle, getNewSprintOpts } from "./SprintTitle";
import {
  showErrorAlertWithRefresh,
  showSuccessAlert,
  showErrorAlert,
} from "./Alerts";

const api = new DefaultApi(window.location.origin + "/api/v1");
const sprintId = "current";

let sprintTemplate: models.SprintTemplate;

window.onload = () => {
  const socket = new WebSocket("ws://" + window.location.host + "/ws");
  socket.onclose = () => {
    showErrorAlertWithRefresh("lost connection to server");
  };

  const input = $("#new_sprint_title")[0] as HTMLInputElement;
  input.placeholder = buildNewSprintTitle();
  setInterval(() => {
    input.placeholder = buildNewSprintTitle();
  }, 60 * 60 * 1000);

  load_task_lists();
};

$("#new_sprint_btn")[0].addEventListener("click", () => {
  const sprintOpts = getNewSprintOpts();
  void api
    .createTaskList(sprintOpts)
    .done((data) => {
      sprintTemplate = data.body;
      load_task_lists();
      showSuccessAlert("sprint created");
    })
    .fail(() => {
      showErrorAlert("failed to create sprint");
    });
});

function load_task_lists() {
  void api
    .getTaskList(sprintId)
    .fail(() => {
      showErrorAlertWithRefresh("failed to load sprint tasks");
    })
    .done((data) => {
      draw_task_lists(data.body);
    });
}

enum TaskProperty {
  Completed = "completed",
  Canceled = "canceled",
  Todo = "todo",
}

function draw_task_lists(sprintTaskList: models.TaskList) {
  prepare_task_list(sprintTaskList);

  update_task_list_header(sprintTaskList);

  fill_task_list(sprintTaskList.tasks);
}

function update_task_list_header(taskList: models.TaskList) {
  const points = sum_points(taskList.tasks);
  const burnt = sum_burnt_points(taskList.tasks);

  const taskListHtml = $(listHtmlId() + " .list_header")[0];
  taskListHtml.getElementsByClassName("title")[0].innerHTML = taskList.title;
  taskListHtml.getElementsByClassName("points")[0].innerHTML =
    burnt.toString() + "/" + points.toString();
}

function fill_task_list(tasks: models.RespTask[]) {
  const taskListHtml = $(listHtmlId() + " .tasks")[0];

  taskListHtml.innerHTML = "";
  tasks.forEach((task) => {
    taskListHtml.append(build_task_html(task));
  });

  if (
    sprintTemplate &&
    sprintTemplate.tasks &&
    sprintTemplate.tasks.length > 0
  ) {
    const task: TaskEditorTask = {
      text: sprintTemplate.tasks[0].text,
      points: sprintTemplate.tasks[0].points,
    };
    sprintTemplate.tasks.splice(0, 1);
    taskListHtml.append(build_template_task_input_html(task));
    focus_new_task_input();
  } else {
    taskListHtml.append(build_new_task_input_html());
  }
}

function build_task_html(task: models.RespTask): HTMLElement {
  let points = task.burnt.toString() + "/" + task.points.toString();
  let percent = (100 * task.burnt) / task.points;
  if (task.state === models.RespTask.StateEnum.Canceled && task.burnt === 0) {
    points = "";
    percent = 0;
  }

  let taskProperties = "";
  if (task.state === models.RespTask.StateEnum.Done) {
    taskProperties = TaskProperty.Completed;
  } else if (task.state === models.RespTask.StateEnum.Canceled) {
    taskProperties = TaskProperty.Canceled;
  } else if (task.state === models.RespTask.StateEnum.Todo) {
    taskProperties = TaskProperty.Todo;
  }

  const taskIdDiv = document.createElement("div");
  taskIdDiv.className = "task__id";
  taskIdDiv.innerText = task.id;

  const taskTextDiv = document.createElement("div");
  taskTextDiv.className = "text";
  taskTextDiv.innerText = task.text;

  const taskPointsDiv = document.createElement("div");
  taskPointsDiv.className = "points";
  taskPointsDiv.innerText = points;
  if (percent > 0) {
    taskPointsDiv.style.background =
      "-webkit-linear-gradient(left, #f8f8f8 " +
      percent.toString() +
      "%, white " +
      percent.toString() +
      "%)";
  }

  const taskDiv = document.createElement("div");
  taskDiv.className = "task " + taskProperties;
  taskDiv.setAttribute("type", "button");
  taskDiv.setAttribute("data-toggle", "dropdown");
  taskDiv.append(taskIdDiv, taskTextDiv, taskPointsDiv);

  const dropdown = document.createElement("div");
  dropdown.className = "dropdown show";
  dropdown.append(taskDiv);

  taskDiv.onclick = () => {
    dropdown.append(build_dropdown_menu(task));
  };

  dropdown.ondblclick = (): boolean => {
    dropdown.replaceWith(build_task_input_html(task, dropdown));
    return false;
  };

  return dropdown;
}

function build_dropdown_menu(task: models.RespTask): HTMLDivElement {
  return BuildDropdownMenu(
    task.state,
    () => {
      void api
        .todoTask(task.id)
        .done(() => {
          load_task_lists();
        })
        .fail(() => {
          showErrorAlert("failed to todo task");
        });
    },
    () => {
      void api
        .doneTask(task.id)
        .done(() => {
          load_task_lists();
        })
        .fail(() => {
          showErrorAlert("failed to done task");
        });
    },
    () => {
      void api
        .cancelTask(task.id)
        .done(() => {
          load_task_lists();
        })
        .fail(() => {
          showErrorAlert("failed to cancel task");
        });
    },
    () => {
      void api
        .toworkTask(task.id)
        .done(() => {
          load_task_lists();
        })
        .fail(() => {
          showErrorAlert("failed to back task to work");
        });
    },
    () => {
      void api
        .postponeTask(task.id)
        .done(() => {
          load_task_lists();
        })
        .fail(() => {
          showErrorAlert("failed to delete task");
        });
    },
    () => {
      void api
        .deleteTask(task.id)
        .done(() => {
          load_task_lists();
        })
        .fail(() => {
          showErrorAlert("failed to delete task");
        });
    }
  );
}

function build_task_input_html(
  task: models.RespTask,
  resetDiv: HTMLElement
): HTMLElement {
  const editorTask: TaskEditorTask = {
    text: task.text,
    points: task.points,
    burnt: task.burnt,
  };

  const autofocusPoints = $(".text:hover").length === 0;
  const focus = autofocusPoints ? TaskEditorFocus.Points : TaskEditorFocus.Text;

  return BuildTaskEditor(
    (text: string, points: string) => {
      const pointsArr = points.split("/");
      const opts: models.UpdateOptions = {
        text,
        points: parseInt(pointsArr[1], 10),
        burnt: parseInt(pointsArr[0], 10),
      };
      void api
        .updateTask(task.id, opts)
        .done(() => {
          load_task_lists();
        })
        .fail(() => {
          showErrorAlert("failed to update task");
        });
      load_task_lists();
    },
    undefined,
    resetDiv,
    editorTask,
    focus
  );
}

function build_new_task_input_html(): HTMLElement {
  return BuildTaskEditor((text: string, points: string) => {
    const newTask: models.Task = {
      text,
      points: parseInt(points, 10),
    };
    void api
      .createTask(sprintId, newTask)
      .done(() => {
        load_task_lists();
        focus_new_task_input();
      })
      .fail(() => {
        showErrorAlert("failed to create task");
      });
  });
}

function build_template_task_input_html(task: TaskEditorTask): HTMLElement {
  return BuildTaskEditor(
    (text: string, points: string) => {
      const newTask: models.Task = {
        text,
        points: parseInt(points, 10),
      };
      void api
        .createTask(sprintId, newTask)
        .done(() => {
          load_task_lists();
          focus_new_task_input();
        })
        .fail(() => {
          showErrorAlert("failed to create task");
        });
    },
    () => {
      load_task_lists();
      focus_new_task_input();
    },
    undefined,
    task
  );
}

function focus_new_task_input() {
  setTimeout(() => {
    $(listHtmlId() + " .text.form-control")[0].focus();
  }, 100);
}

function sum_points(tasks: models.RespTask[]): number {
  return tasks.reduce((sum, current) => {
    if (current.state == models.RespTask.StateEnum.Canceled) {
      return sum + current.burnt;
    }
    return sum + current.points;
  }, 0);
}

function sum_burnt_points(tasks: models.RespTask[]): number {
  return tasks.reduce((sum, current) => {
    return sum + current.burnt;
  }, 0);
}

function prepare_task_list(taskList: models.TaskList): void {
  const fixPoints = (value: models.RespTask) => {
    if (!value.burnt) {
      value.burnt = 0;
    }

    switch (value.state) {
      case models.RespTask.StateEnum.Done:
        value.burnt = value.points;
        break;
    }
  };

  taskList.tasks.forEach(fixPoints);
}

function listHtmlId(): string {
  return "#sprint_list";
}
