import * as models from "./openapi_cli/model/models";

export function BuildDropdownMenu(
  taskState: models.RespTask.StateEnum,
  todoEL: EventListener,
  doneEL: EventListener,
  cancelEL: EventListener,
  toworkEL: EventListener,
  postponeEL: EventListener,
  deleteEL: EventListener
): HTMLDivElement {
  const dropdownMenu = document.createElement("div");
  dropdownMenu.className = "dropdown-menu";

  switch (taskState) {
    case models.RespTask.StateEnum.Todo:
      appendItem(dropdownMenu, "Done", doneEL);
      appendItem(dropdownMenu, "Cancel", cancelEL);
      appendItem(dropdownMenu, "Postpone", postponeEL);
      break;
    case models.RespTask.StateEnum.Done:
      break;
    case models.RespTask.StateEnum.Canceled:
      appendItem(dropdownMenu, "ToWork", toworkEL);
      break;
    default:
      appendItem(dropdownMenu, "Done", doneEL);
      appendItem(dropdownMenu, "Todo", todoEL);
      appendItem(dropdownMenu, "Cancel", cancelEL);
      appendItem(dropdownMenu, "Postpone", postponeEL);
      break;
  }

  appendItem(dropdownMenu, "Delete", deleteEL);
  return dropdownMenu;
}

function appendItem(
  dropdownMenu: HTMLDivElement,
  text: string,
  handler: EventListener
) {
  const action = document.createElement("div");
  action.className = "dropdown-item";
  action.innerText = text;
  action.addEventListener("click", handler);
  dropdownMenu.appendChild(action);
}
