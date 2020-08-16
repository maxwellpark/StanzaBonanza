// Container for Poem components

import React from "react";
import { useSelector, connect } from "react-redux";
import Poem from "./Poem";
import AddPoem from "./AddPoem";
import Pagination from "./Pagination";
import { paginatePoems } from "../utils/paginatePoems";
import { useState } from "react";
import { Container, Row, Col, Card } from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";

const Poems = ({ collection }) => {
  const state = useSelector((state) => state);
  const poems = collection == "main" ? state.poems : state.hallOfFamers;
  const poemsPerPage = state.poemsPerPage;
  const poemCount = poems.length;

  const [currentPage, setCurrentPage] = useState(1);

  // See utils for implementation details
  const paginatedPoems = paginatePoems(poems, currentPage, poemsPerPage);

  const handlePageChange = (page) => {
    console.log("page no. :", page);
    setCurrentPage(page);
    console.log("current page: ", currentPage);
  };

  return (
    <Container fluid className="poem-container">
      <Row className="justify-content-md-center">
        <Col>{collection == "main" ? <AddPoem /> : null}</Col>
        <Col>
          <Pagination
            poemCount={poemCount}
            poemsPerPage={poemsPerPage}
            currentPage={currentPage}
            handlePageChange={handlePageChange}
          />
        </Col>
      </Row>
      <Row>
        {paginatedPoems.map((poem) => {
          return (
            <Col md="auto">
              <Card key={poem.id} className="poem-card">
                <Poem
                  author={poem.author}
                  title={poem.title}
                  text={poem.text}
                  date={poem.date}
                />
              </Card>
            </Col>
          );
        })}
      </Row>
    </Container>
  );
};

const mapStateToProps = (state, ownProps) => {
  switch (ownProps.collection) {
    case "main":
      return {
        poems: state.poems,
        poemsPerPage: state.poemsPerPage,
      };
    case "hallOfFame":
      return {
        poems: state.hallOfFamers,
        poemsPerPage: state.poemsPerPage,
      };
    default:
      return null;
  }
};

export default connect(mapStateToProps)(Poems);
