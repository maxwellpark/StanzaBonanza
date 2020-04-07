import React from "react";
import { Link } from "react-router-dom";
import "bootstrap/dist/css/bootstrap.min.css";

const Navbar = () => {
  return (
    <nav className="navbar navbar-light bg-light static-top">
      <div className="container">
        <Link to="/" className="navbar-brand">
          Home
        </Link>
        <Link to="/hall-of-fame" className="btn btn-primary">
          Hall Of Fame
        </Link>
      </div>
    </nav>
  );
};
export default Navbar;
