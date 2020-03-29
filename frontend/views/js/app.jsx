//React app
const RESTENDPOINT = "localhost:8081"
const SEARCHENDPOINT = RESTENDPOINT + "/search/"

// Just the Searchbar Component with the logo, title and the searchbar to submit a url.
class SearchBar extends React.Component {

  state = {
    value: ""
  }

  // CTor for the Searchbar Component.
  constructor(props) {
    super(props);
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  // Making the POST REQUEST
  serverRequest(url) {
    $.post(
      SEARCHENDPOINT,
      {
        url: this.state.value
      },
      response => {
        // TODO
      }
    );
  }

  // Handling the Change of the Inputfield. 
  handleChange(event) {    
    this.setState(
      { 
        value: event.target.value
      }
    );  
  }

  // Handle the submit of the Form.
  handleSubmit(event) {
    alert('A Website was submitted: ' + this.state.value);
  }

  render() {
    return (
      <div class="searchbar">
        <img class="logo" alt="sherlock gopher" src="img/sherlockgopher.png"></img>
        <h4 class="title">Sherlock Gopher</h4>
        <form onSubmit={this.handleSubmit} class="searchbarform" method="POST">
          <input class="searchbarinput" value={this.state.value} onChange={this.handleChange} name="url" placeholder="Hier könnte Ihre Werbung stehen..."></input>
          <button class="searchbarbutton" type="submit">Schnüffel</button>
        </form>
      </div>
    )
  }
}

// The actual APP Component which will be rendered by calling the Website.
class App extends React.Component {
    render() {
      return (
        <div>
            <SearchBar></SearchBar>
        </div>
      )
    }
}

ReactDOM.render(<App />, document.getElementById('app'));
