import React from "react";
import ReactDOM from "react-dom";
import { Provider } from "react-redux";
import { createStore } from "redux";
import { addDummyPoem } from "./actions/poemActions";
import rootReducer from "./reducers/rootReducer";
import poemReducer from "./reducers/poemReducer";
import App from "./components/App";

let store = createStore(poemReducer);

store.subscribe(() => {
  console.log("redux store updated!");
  console.log(store.getState());
});

const dummies = 20;

for (let i = 0; i < dummies; i++) {
  store.dispatch(addDummyPoem);
}

ReactDOM.render(
  <Provider store={store}>
    <App />,
  </Provider>,
  document.getElementById("root"),
);
