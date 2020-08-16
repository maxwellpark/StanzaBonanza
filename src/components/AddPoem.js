import React from "react";
import { addPoemAction } from "../actions/poemActions";
import { useDispatch } from "react-redux";
import { Card } from "react-bootstrap";
import moment from "moment";

const AddPoem = (props) => {
  let author, title, text;
  const dispatch = useDispatch();

  return (
    <Card className="add-poem">
      <textarea
        form="poem_form"
        onChange={(e) => (text = e.target.value)}
        value={text}
        placeholder="Create a poem!"
      />
      <form
        id="poem_form"
        onSubmit={(e) => {
          e.preventDefault();
          console.log("Text: ", text);
          dispatch(
            addPoemAction(author, title, text, moment().format("DD-MM-YYYY")),
          );
        }}
      >
        {/* Looking to replace with ReactQuill in the future */}
        <input
          type="text"
          value={title}
          onChange={(e) => (title = e.target.value)}
          name="title"
          placeholder="Title"
          required
        />
        <input
          type="text"
          value={author}
          onChange={(e) => (author = e.target.value)}
          name="author"
          placeholder="Your name"
          required
        />
        <button type="submit">Publish!</button>
      </form>
    </Card>
  );
};
export default AddPoem;
