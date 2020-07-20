import * as models from "./openapi_cli/model/models";

export function getNewSprintOpts(): models.SprintOpts {
  const { begin, end } = getNewSprintDates();
  return {
    title: formatSrpintTitle(begin, end),
    begin: formatSrpintDate(begin),
    end: formatSrpintDate(end),
  };
}

function formatSrpintDate(d: Date): string {
  return (
    d.getUTCFullYear().toString() +
    "-" +
    numToString(d.getMonth() + 1) +
    "-" +
    numToString(d.getDate())
  );
}

export function buildNewSprintTitle(): string {
  const { begin, end } = getNewSprintDates();
  return formatSrpintTitle(begin, end);
}

export function getNewSprintDates(): { begin: Date; end: Date } {
  const date = new Date();
  const dayOfWeek = date.getDay();

  date.setDate(date.getDate() + 1 - dayOfWeek);
  const beginDate = date;
  const endDate = new Date(date);

  endDate.setDate(date.getDate() + 6);

  return { begin: beginDate, end: endDate };
}

function formatSrpintTitle(begin: Date, end: Date): string {
  return dateToString(begin) + " - " + dateToString(end);
}

function dateToString(d: Date): string {
  return numToString(d.getDate()) + "." + numToString(d.getMonth() + 1);
}

function numToString(v: number): string {
  let s = v.toString();
  if (v < 10) {
    s = "0" + s;
  }
  return s;
}
