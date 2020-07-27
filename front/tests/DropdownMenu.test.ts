import { BuildDropdownMenu } from "../src/DropdownMenu";
import { RespTask } from "../src/openapi_cli/index";

describe("list elements", () => {
  test.each<[RespTask.StateEnum, string[]]>([
    [RespTask.StateEnum.Todo, ["Done", "Cancel", "Postpone", "Delete"]],
    [
      RespTask.StateEnum.Empty,
      ["Done", "Todo", "Cancel", "Postpone", "Delete"],
    ],
    [RespTask.StateEnum.Done, ["Delete"]],
    [RespTask.StateEnum.Canceled, ["Delete"]],
  ])('list "%s" should contains: %s', (type, items) => {
    const menuDiv = BuildDropdownMenu(
      type,
      jest.fn(),
      jest.fn(),
      jest.fn(),
      jest.fn(),
      jest.fn()
    );

    expect(menuDiv.className).toEqual("dropdown-menu");
    expect(menuDiv.childNodes).toHaveLength(items.length);

    items.forEach((item: string, i: number) => {
      const menuItem = menuDiv.children[i] as HTMLDivElement;
      expect(menuItem.innerText).toEqual(item);
    });
  });
});

describe("list item actions", () => {
  const fnMap = new Map<string, jest.Mock<unknown, unknown[]>>();
  fnMap.set("Done", jest.fn());
  fnMap.set("Todo", jest.fn());
  fnMap.set("Cancel", jest.fn());
  fnMap.set("Postpone", jest.fn());
  fnMap.set("Delete", jest.fn());

  const menuDiv = BuildDropdownMenu(
    RespTask.StateEnum.Empty,
    fnMap.get("Todo"),
    fnMap.get("Done"),
    fnMap.get("Cancel"),
    fnMap.get("Postpone"),
    fnMap.get("Delete")
  );

  expect(menuDiv.childNodes).toHaveLength(fnMap.size);

  menuDiv.childNodes.forEach((menuItem: HTMLDivElement) => {
    it("click on " + menuItem.innerText.toString(), () => {
      menuItem.click();
      expect(fnMap.get(menuItem.innerText)).toBeCalled();
    });
  });
});
