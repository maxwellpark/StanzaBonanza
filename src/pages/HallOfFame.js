import React from "react";
import Navbar from "../components/Navbar";
import Poems from "../components/Poems";

const HallOfFame = () => {
  return (
    <div id="hall-of-fame">
      <Navbar />
      <div className="jumbotron">
        <h1>Hall Of Fame</h1>
        <p>Where superlative verses go to die</p>
      </div>
      <Poems collection="hallOfFame" />
    </div>
  );
};

export default HallOfFame;
