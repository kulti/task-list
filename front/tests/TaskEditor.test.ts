import * as models from "../src/openapi_cli/model/models"

describe('build task editor', () => {
    const BuildTaskEditor = require('../src/TaskEditor').BuildTaskEditor

    const task: models.RespTask = {
        id: 'any id',
        text: "test task",
        points: 3,
        burnt: 1,
        state: models.RespTask.StateEnum.Todo,
    }

    it('has specific html for new task', () => {
        const taskDiv = BuildTaskEditor({},) as HTMLDivElement

        expect(taskDiv.childNodes).toHaveLength(2)

        const textInput = taskDiv.childNodes[0] as HTMLInputElement
        expect(textInput.className).toContain('text')
        expect(textInput.value).toEqual('')
        expect(textInput.placeholder).toEqual('Do new task')

        const pointsInput = taskDiv.childNodes[1] as HTMLInputElement
        expect(pointsInput.className).toContain('points')
        expect(pointsInput.value).toEqual('')
        expect(pointsInput.placeholder).toEqual('0')
    });

    it('has specific html for edit task', () => {
        const taskDiv = BuildTaskEditor({}, {}, task) as HTMLDivElement

        expect(taskDiv.childNodes).toHaveLength(2)

        const textInput = taskDiv.childNodes[0] as HTMLInputElement
        expect(textInput.className).toContain('text')
        expect(textInput.value).toEqual(task.text)

        const pointsInput = taskDiv.childNodes[1] as HTMLInputElement
        expect(pointsInput.className).toContain('points')
        expect(pointsInput.value).toEqual(task.burnt + "/" + task.points)
    });

    it('handles enter key as apply', () => {
        const applyFn = jest.fn()
        const taskDiv = BuildTaskEditor(applyFn) as HTMLDivElement
        const textInput = taskDiv.childNodes[0] as HTMLInputElement
        const pointsInput = taskDiv.childNodes[1] as HTMLInputElement

        const text = "test text"
        const points = "5"
        textInput.value = text
        pointsInput.value = points

        const ev = new KeyboardEvent('keyup', { 'key': 'Enter' });
        textInput.dispatchEvent(ev)

        expect(applyFn).toBeCalledWith(text, points)
    });

    it('handles escape key to cleanup inputs', () => {
        const applyFn = jest.fn()
        const taskDiv = BuildTaskEditor(applyFn) as HTMLDivElement
        const textInput = taskDiv.childNodes[0] as HTMLInputElement
        const pointsInput = taskDiv.childNodes[1] as HTMLInputElement

        const text = "some text"
        const points = "1"
        textInput.value = text
        pointsInput.value = points

        const ev = new KeyboardEvent('keyup', { 'key': 'Escape' });
        textInput.dispatchEvent(ev)

        expect(applyFn).not.toBeCalled()
        expect(textInput.value).toEqual('')
        expect(pointsInput.value).toEqual('')
    });

    it('handles escape key to apply reset div', () => {
        const applyFn = jest.fn()
        const mainDiv = document.createElement('div') as HTMLDivElement
        const resetDiv = document.createElement('div') as HTMLDivElement
        resetDiv.className = 'reset'

        const taskDiv = BuildTaskEditor(applyFn, resetDiv) as HTMLDivElement
        const textInput = taskDiv.childNodes[0] as HTMLInputElement
        const pointsInput = taskDiv.childNodes[1] as HTMLInputElement

        mainDiv.append(taskDiv)
        const text = "some text"
        const points = "1"
        textInput.value = text
        pointsInput.value = points

        const ev = new KeyboardEvent('keyup', { 'key': 'Escape' });
        textInput.dispatchEvent(ev)

        expect(applyFn).not.toBeCalled()
        expect(mainDiv.innerHTML).toEqual(resetDiv.outerHTML)
    });
});
