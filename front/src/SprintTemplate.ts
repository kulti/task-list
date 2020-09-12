import * as models from "./openapi_cli/model/models";
import { showErrorAlert } from "./Alerts";

export function BuildSprintTemplateEditor(
  template: models.SprintTemplate,
  applyFn: (template: models.SprintTemplate) => void,
  cancelFn: () => void
): HTMLElement {
  const formDiv = document.createElement("div");
  formDiv.className = "form-group";

  const textArea = document.createElement("textarea");
  textArea.className = "form-control";

  if (template.tasks) {
    const taskStrs: string[] = [];
    template.tasks.forEach((task) => {
      taskStrs.push("[" + task.points.toString() + "] " + task.text);
    });
    textArea.setRangeText(taskStrs.join("\n"));
    textArea.setAttribute("rows", taskStrs.length.toString());
  }

  formDiv.append(textArea);

  const applyBtn = document.createElement("button");
  applyBtn.className = "btn btn-outline-success";
  applyBtn.type = "button";
  applyBtn.innerText = "Apply";
  applyBtn.onclick = () => {
    const template: models.SprintTemplate = { tasks: [] };
    const regex = new RegExp("^\\[([0-9]+)\\] (.+)$");

    for (const task of textArea.value.split("\n")) {
      const res = regex.exec(task);
      if (!res) {
        showErrorAlert("invalid task string: " + task);
        return;
      }

      const points = parseInt(res[1]);
      if (isNaN(points)) {
        showErrorAlert("invalid task points: " + task);
        return;
      }
      const tmplTask: models.TaskTemplate = {
        id: "",
        points: points,
        text: res[2],
      };
      template.tasks.push(tmplTask);
    }

    applyFn(template);
  };
  formDiv.append(applyBtn);

  const cancelBtn = document.createElement("button");
  cancelBtn.className = "btn btn-outline-danger";
  cancelBtn.type = "button";
  cancelBtn.innerText = "Cancel";
  cancelBtn.onclick = cancelFn;
  formDiv.append(cancelBtn);

  return formDiv;
}
