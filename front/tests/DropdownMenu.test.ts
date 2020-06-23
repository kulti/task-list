import * as models from "../src/openapi_cli/model/models"

describe('list elements', () => {
    test.each`
        type          | items
        ${'todo'}     | ${Array('Done', 'Cancel', 'Delete')}
        ${''}         | ${Array('Done', 'Todo', 'Cancel', 'Delete')}
        ${'done'}     | ${Array('Delete')}
        ${'canceled'} | ${Array('Delete')}
    `(
        'list "$type" should contains: $items',
        ({ type, items }) => {
            const DropdownMenu = require('../src/DropdownMenu').DropdownMenu

            const parent = document.createElement('div') as HTMLDivElement;
            // tslint:disable-next-line: no-unused-expression
            new DropdownMenu(parent, type, {}, {}, {}, {})

            expect(parent.childNodes).toHaveLength(1)
            expect(parent.firstElementChild.className).toEqual('dropdown-menu')

            const menuDiv = parent.firstElementChild
            expect(menuDiv.childNodes).toHaveLength(items.length)
            items.forEach((item: string, i: number) => {
                const menuItem = menuDiv.children[i] as HTMLDivElement
                expect(menuItem.innerText).toEqual(item)
            });
        }
    )
});

describe('list item actions', () => {
    const fnMap = new Map<string, jest.Mock<any, any>>();
    fnMap.set('Done', jest.fn())
    fnMap.set('Todo', jest.fn())
    fnMap.set('Cancel', jest.fn())
    fnMap.set('Delete', jest.fn())

    const DropdownMenu = require('../src/DropdownMenu').DropdownMenu
    const parent = document.createElement('div') as HTMLDivElement;
    // tslint:disable-next-line: no-unused-expression
    new DropdownMenu(parent, '', fnMap.get('Todo'), fnMap.get('Done'), fnMap.get('Cancel'), fnMap.get('Delete'))

    const menuDiv = parent.firstElementChild
    expect(menuDiv.childNodes).toHaveLength(4)

    menuDiv.childNodes.forEach((menuItem: HTMLDivElement) => {
        it('click on ' + menuItem.innerText.toString(), () => {
            menuItem.click()
            expect(fnMap.get(menuItem.innerText)).toBeCalled()
        });
    });
});
