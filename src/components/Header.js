import React from "react";
import Navbar from "../components/Navbar";
import "bootstrap/dist/css/bootstrap.min.css";

const Header = () => {
  return (
    <>
      <Navbar />
      <div className="jumbotron">
        <h1>Stanza Bonanza</h1>
        <p>Boundless collective poetry</p>
      </div>
    </>
  );
};
export default Header;
