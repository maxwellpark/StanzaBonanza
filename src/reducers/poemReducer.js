import Moment from "moment";

const initialState = {
  poemsPerPage: 5,

  // Pre-composed poems
  poems: [
    {
      id: 0,
      author: "Benevolent Pricktator",
      title: "The Semblance of Liberty",

      // Backticks allow multi-line string,
      // which allows poem lines to be represented
      text: `Come, queue up to kiss my feet
        And hear the starlings forced to tweet.
        All shall blossom under my reign.
        Even the heretics will feel no pain.`,
      date: "06-09-1969",
    },
    {
      id: 1,
      author: "Membrane",
      title: "Puncture My Substrate",
      text: `Oh how I loath their empty spaces
        Heckling me without divison.
        Such immeasurable solidarity
        Ripples through my weary veneer.
        When will the pressure rise?
        Sufficiently to lift my disguise
        This is far too much torment for a boundary to endure.`,
      date: "28-11-2018",
    },
    {
      id: 2,
      author: "A Donkey Illustrator",
      title: "Pinned Down",
      text: `Deafened by the quotidian vortex
        Orientation less than impossible
        No notion of a posterior appendage whatsoever.`,
      date: "05-01-2019",
    },
    {
      id: 3,
      author: "The Disgrace",
      title: "If Only They Knew",
      text: `Dealt a diabolical hand
        Groped by the system.
        Congenital heartbreak and an ash tray
        Were the only respite in sight.`,
      date: "16-03-1999",
    },
  ],
  // These cannot be extended by other poets
  hallOfFamers: [
    {
      id: 0,
      author: "Mended By Wendy",
      title: "Fat Cow",
      text: `That bitch stole my dignity, 
        Can't fathom what she did to me.
        I'd be less insulted,
        If she wasn't so damn bloated. 
        `,
      date: "09-05-2003",
    },
    {
      id: 1,
      author: "Laslowe",
      title: "Transition",
      text: `Bouncing back induced a splatter,
      Invalidating their claims to freedom.
      Now less significant than the pitter patter,
      Of a tepid deluge.`,
      date: "18-12-2010",
    },
  ],
  extensions: [
    {
      // Id corresponds to a poem id
      id: 3,
      newStanza: true,
      text: "Witnesses weren't even contrite...",
      literaryDevice: "emphasis",
    },
  ],
};

const poemReducer = (state = initialState, action) => {
  console.log(action, state, "reducer body!");
  switch (action.type) {
    case "ADD_POEM":
      return { ...state, poems: [action.poem, ...state.poems] };

    case "EXTEND_POEM":
      return {
        ...state,
        extensions: [action.extensions, ...state.extensions],
      };

    case "DELETE_POEM":
      return {
        ...state,
        poems: [...state.poems, (action.poem = null)],
      };

    case "ADD_DUMMY_POEM":
      return {
        ...state,
        poems: [...state.poems, action.poem],
      };

    default:
      return state;
  }
};
export default poemReducer;
