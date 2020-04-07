import React from "react";
import Header from "../components/Header";
import Poems from "../components/Poems";

const Home = () => {
  return (
    <div id="home">
      <Header />
      <Poems collection="main" />
    </div>
  );
};
export default Home;
