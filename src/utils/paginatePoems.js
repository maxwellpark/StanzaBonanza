import _ from "lodash";

export const paginatePoems = (poems, pageNumber, poemsPerPage) => {
  const pageIndex = (pageNumber - 1) * poemsPerPage;
  const paginatedPoems = _(poems).slice(pageIndex).take(poemsPerPage).value();

  console.log("paginated method return value: ", paginatedPoems);
  return paginatedPoems;
};
