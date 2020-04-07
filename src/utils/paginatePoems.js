import _ from "lodash";

export const paginatePoems = (poems, pageNumber, poemsPerPage) => {
  const pageIndex = (pageNumber - 1) * poemsPerPage;
  return _(poems).slice(pageIndex).take(poemsPerPage).value();
};
