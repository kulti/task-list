import * as models from "./openapi_cli/model/models"

export class DropdownMenu {
    dropdownMenu: HTMLDivElement;

    constructor(parent: Node, taskState: models.RespTask.StateEnum, todoEL: EventListener,
        doneEl: EventListener, cancelEL: EventListener, deleteEL: EventListener) {
        this.dropdownMenu = document.createElement('div') as HTMLDivElement;
        this.dropdownMenu.className = "dropdown-menu";

        switch (taskState) {
            case models.RespTask.StateEnum.Todo:
                this.appendItem("Done", doneEl)
                this.appendItem("Cancel", cancelEL)
                break;
            case models.RespTask.StateEnum.Done:
                break;
            case models.RespTask.StateEnum.Canceled:
                break;
            default:
                this.appendItem("Done", doneEl)
                this.appendItem("Todo", todoEL)
                this.appendItem("Cancel", cancelEL)
                break;
        }

        this.appendItem("Delete", deleteEL)

        parent.appendChild(this.dropdownMenu)
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
