import React from "react";
import "font-awesome/css/font-awesome.min.css";
import "bootstrap/dist/css/bootstrap.min.css";

const Like = ({ liked, handleLike }) => {
  return !liked ? (
    <i onClick={() => handleLike()} className="fas fa-feather-alt"></i>
  ) : (
    <i onClick={() => handleLike()} className="far fa-feather-alt"></i>
  );
};
export default Like;
