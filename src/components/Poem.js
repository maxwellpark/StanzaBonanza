import React, { useState } from "react";
// import ReactQuill from "react-quill";
import Like from "./Like";
import { Col } from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";

const Poem = props => {
  const [liked, setLiked] = useState(false);

  const handleLike = () => {
    !liked ? setLiked(true) : setLiked(false);
  };

  return (
    <Col md="auto">
      <h3 className="poem-title">{props.title}</h3>
      <h4 className="poem-author">By: {props.author}</h4>
      <p>{props.text}</p>
      <p>Created: {props.date}</p>
      <Like liked={liked} handleLike={handleLike} />
      <button>Extend</button>
    </Col>
  );
};
export default Poem;
