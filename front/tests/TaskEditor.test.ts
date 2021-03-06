import { BuildTaskEditor, TaskEditorTask } from "../src/TaskEditor";

describe("build task editor", () => {
  const task: TaskEditorTask = {
    text: "test task",
    points: 3,
    burnt: 1,
  };

  it("has specific html for new task", () => {
    const taskDiv = BuildTaskEditor(jest.fn());

    expect(taskDiv.childNodes).toHaveLength(2);

    const textInput = taskDiv.childNodes[0] as HTMLInputElement;
    expect(textInput.className).toContain("text");
    expect(textInput.value).toEqual("");
    expect(textInput.placeholder).toEqual("Do new task");

    const pointsInput = taskDiv.childNodes[1] as HTMLInputElement;
    expect(pointsInput.className).toContain("points");
    expect(pointsInput.value).toEqual("");
    expect(pointsInput.placeholder).toEqual("0");
  });

  it("has specific html for edit task", () => {
    const taskDiv = BuildTaskEditor(jest.fn(), undefined, undefined, task);

    expect(taskDiv.childNodes).toHaveLength(2);

    const textInput = taskDiv.childNodes[0] as HTMLInputElement;
    expect(textInput.className).toContain("text");
    expect(textInput.value).toEqual(task.text);

    const pointsInput = taskDiv.childNodes[1] as HTMLInputElement;
    expect(pointsInput.className).toContain("points");
    expect(pointsInput.value).toEqual(
      task.burnt.toString() + "/" + task.points.toString()
    );
  });

  it("has specific html for edit task with zero-burnt", () => {
    const zeroBurntTask = { ...task };
    zeroBurntTask.burnt = 0;
    const taskDiv = BuildTaskEditor(
      jest.fn(),
      undefined,
      undefined,
      zeroBurntTask
    );

    expect(taskDiv.childNodes).toHaveLength(2);

    const pointsInput = taskDiv.childNodes[1] as HTMLInputElement;
    expect(pointsInput.className).toContain("points");
    expect(pointsInput.value).toEqual("0/" + task.points.toString());
  });

  it("has specific html for edit task without burnt points", () => {
    const noBurntTask = task;
    noBurntTask.burnt = undefined;
    const taskDiv = BuildTaskEditor(
      jest.fn(),
      undefined,
      undefined,
      noBurntTask
    );

    expect(taskDiv.childNodes).toHaveLength(2);

    const pointsInput = taskDiv.childNodes[1] as HTMLInputElement;
    expect(pointsInput.className).toContain("points");
    expect(pointsInput.value).toEqual(task.points.toString());
  });

  it("handles enter key as apply", () => {
    const applyFn = jest.fn();
    const taskDiv = BuildTaskEditor(applyFn) as HTMLDivElement;
    const textInput = taskDiv.childNodes[0] as HTMLInputElement;
    const pointsInput = taskDiv.childNodes[1] as HTMLInputElement;

    const text = "test text";
    const points = "5";
    textInput.value = text;
    pointsInput.value = points;

    const ev = new KeyboardEvent("keyup", { key: "Enter" });
    textInput.dispatchEvent(ev);

    expect(applyFn).toBeCalledWith(text, points);
  });

  it("handles escape key to cleanup inputs", () => {
    const applyFn = jest.fn();
    const cancelFn = jest.fn();
    const taskDiv = BuildTaskEditor(applyFn, cancelFn) as HTMLDivElement;
    const textInput = taskDiv.childNodes[0] as HTMLInputElement;
    const pointsInput = taskDiv.childNodes[1] as HTMLInputElement;

    const text = "some text";
    const points = "1";
    textInput.value = text;
    pointsInput.value = points;

    const ev = new KeyboardEvent("keyup", { key: "Escape" });
    textInput.dispatchEvent(ev);

    expect(applyFn).not.toBeCalled();
    expect(cancelFn).toBeCalled();
    expect(textInput.value).toEqual("");
    expect(pointsInput.value).toEqual("");
  });

  it("handles escape key to apply reset div", () => {
    const applyFn = jest.fn();
    const mainDiv = document.createElement("div");
    const resetDiv = document.createElement("div");
    resetDiv.className = "reset";

    const taskDiv = BuildTaskEditor(
      applyFn,
      undefined,
      resetDiv
    ) as HTMLDivElement;
    const textInput = taskDiv.childNodes[0] as HTMLInputElement;
    const pointsInput = taskDiv.childNodes[1] as HTMLInputElement;

    mainDiv.append(taskDiv);
    const text = "some text";
    const points = "1";
    textInput.value = text;
    pointsInput.value = points;

    const ev = new KeyboardEvent("keyup", { key: "Escape" });
    textInput.dispatchEvent(ev);

    expect(applyFn).not.toBeCalled();
    expect(mainDiv.innerHTML).toEqual(resetDiv.outerHTML);
  });
});
