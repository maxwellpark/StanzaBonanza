import React, { useState, useEffect } from "react";
import Home from "../pages/Home";
import HallOfFame from "../pages/HallOfFame";
import { BrowserRouter, Route, Switch } from "react-router-dom";



const App = () => {
  const [apiResponse, setApiResponse] = useState(); 

  const callAPI = () => {
    fetch("https://localhost:9000/demoAPI")
      .then(res => res.text())
      .then(res => setApiResponse(res))
      .catch(err => err);
  }

  useEffect(() => {
    callAPI(); 
  })

  return (
    <div id="app">
      <p className="response-text">Response: {apiResponse}</p>
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