import React, { useLayoutEffect } from "react";
import _ from "lodash";
import "bootstrap/dist/css/bootstrap.min.css";

const Pagination = (props) => {
  const pageCount = Math.ceil(props.poemCount / props.poemsPerPage);
  console.log("page count: ", pageCount);
  const pages = _.range(1, pageCount + 1);

  if (pageCount === 1) return null;

  return (
    <nav>
      <ul className="pagination">
        {pages.map((page) => {
          return (
            <li
              key={page}
              className={
                page === props.currentPage ? "page-item active" : "page-item"
              }
            >
              <a
                onClick={() => props.handlePageChange(page)}
                className="page-link"
              >
                {page}
              </a>
            </li>
          );
        })}
      </ul>
    </nav>
  );
};

export default Pagination;
