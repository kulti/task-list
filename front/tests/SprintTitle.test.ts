import {
  getNewSprintDates,
  buildNewSprintTitle,
  getNewSprintOpts,
} from "../src/SprintTitle";

const sunday = 0;
const monday = 1;

const RealDate = Date;

function mockDate(d: Date) {
  (global as any).Date = class extends RealDate {
    constructor(fromDate?: Date) {
      super();
      if (fromDate) {
        return new RealDate(fromDate);
      }
      return d;
    }
  };
}

describe("new sprint dates", () => {
  afterAll(() => {
    global.Date = RealDate;
  });

  it("begins from monday and ends on sunday", () => {
    mockDate(new RealDate());

    const { begin, end } = getNewSprintDates();
    expect(begin.getDay()).toEqual(monday);
    expect(end.getDay()).toEqual(sunday);
  });

  it("rounds to next week on sunday", () => {
    const mockedDate = new RealDate(2019, 6, 7);
    expect(mockedDate.getDay()).toEqual(sunday);
    mockDate(mockedDate);

    const nextMonday = new RealDate(2019, 6, 8);
    const nextSunday = new RealDate(2019, 6, 14);

    const { begin, end } = getNewSprintDates();
    expect(begin).toEqual(nextMonday);
    expect(end).toEqual(nextSunday);
  });

  it("still in current week on monday to saturday", () => {
    const curMonday = new RealDate(2019, 6, 1);
    const curSunday = new RealDate(2019, 6, 7);

    for (let date = 1; date <= 6; date++) {
      const mockedDate = new RealDate(2019, 6, date);
      mockDate(mockedDate);

      const { begin, end } = getNewSprintDates();
      expect(begin).toEqual(curMonday);
      expect(end).toEqual(curSunday);
    }
  });
});

describe("build sprint title", () => {
  afterAll(() => {
    global.Date = RealDate;
  });

  it("with single digit dates", () => {
    const mockedDate = new RealDate(2019, 6, 3);
    mockDate(mockedDate);

    expect(buildNewSprintTitle()).toEqual("01.07 - 07.07");
  });

  it("with two digits dates", () => {
    const mockedDate = new RealDate(2019, 10, 15);
    mockDate(mockedDate);

    expect(buildNewSprintTitle()).toEqual("11.11 - 17.11");
  });
});

describe("build new sprint options", () => {
  afterAll(() => {
    global.Date = RealDate;
  });

  it("formats title and dates", () => {
    const mockedDate = new RealDate(2019, 6, 1);
    mockDate(mockedDate);

    const opts = getNewSprintOpts();
    expect(opts.title).toEqual("01.07 - 07.07");
    expect(opts.begin).toEqual("2019-07-01");
    expect(opts.end).toEqual("2019-07-07");
  });
});
