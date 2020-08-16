export const addPoemAction = (author, title, text, date) => {
  return {
    type: "ADD_POEM",
    poem: {
      author: author,
      title: title,
      text: text,
      date: date,
    },
  };
};

export const extendPoemAction = {
  type: "EXTEND_POEM",
  poem: {
    coauthor: "Nelson Lurgy",
    text: `Take that to the bank, You ungrateful tosser.`,

    // Will use this to colour code the extension
    // based on the literary device it provides
    literaryDevice: "emphasis",
  },
};

export const deletePoemAction = {
  type: "DELETE_POEM",
};

// Used to test page-filling dynamic
// and pagination
export const addDummyPoem = {
  type: "ADD_DUMMY_POEM",
  poem: {
    author: "Lorem Ipsum",
    title: "Lorem Ipusm",
    text: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. In sed libero non sem accumsan iaculis. Mauris feugiat sem vel velit auctor, vel porta ex vulputate. Cras ac est id diam porttitor aliquam. Vestibulum accumsan cursus interdum. Pellentesque elit est, commodo.`,
  },
};
