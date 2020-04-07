import React from "react";
import Home from "../pages/Home";
import HallOfFame from "../pages/HallOfFame";
import { BrowserRouter, Route, Switch } from "react-router-dom";

const App = () => {
  return (
    <div id="app">
      <BrowserRouter>
        <Switch>
          <Route exact path="/" component={Home} />
          <Route exact path="/hall-of-fame" component={HallOfFame} />
          <Route path="/" render={() => <div>404</div>} />
        </Switch>
      </BrowserRouter>
    </div>
  );
};
export default App;
