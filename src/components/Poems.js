import React from "react";
// import { bindActionCreators } from "redux";
import { useSelector, connect } from "react-redux";
import Poem from "./Poem";
import AddPoem from "./AddPoem";
import Pagination from "./Pagination";
import { paginatePoems } from "../utils/paginatePoems";
import { useState } from "react";
import { Container, Row, Col, Card } from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";

const Poems = ({ collection }) => {
  const poems = useSelector((state) => state.poems);

  console.log("selector poems: ", poems);
  console.log("poems per page: ", poems.poemsPerPage);
  const poemCount = poems.length;
  console.log("poem count: ", poemCount);
  const [currentPage, setCurrentPage] = useState(1);

  let poemsPerPage = 5;
  // See utils for implementation details
  const paginatedPoems = paginatePoems(poems, currentPage, poemsPerPage);
  console.log("paginatedPoems type: ", typeof paginatedPoems);
  console.log("paginated poems: ", poems);

  const handlePageChange = (page) => {
    console.log("page no. :", page);
    setCurrentPage(page);
    console.log("current page: ", currentPage);
  };

  return (
    <Container fluid>
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
