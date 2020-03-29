//React app
import axios from 'axios';

const RESTENDPOINT = "localhost:8081"

class SearchBar extends React.Component {

  constructor(props) {
    super(props);
    this.state = { url: '' };
    this.handleChange = this.handleChange.bind(this)
    this.handleSubmit = this.handleSubmit.bind(this)

  }

  handleChange(event) {
    this.setState({
      url: event.target.url
    })
    console.log(this.state.url)
  }

  handleSubmit(event) {
    alert("This will be sent of " + this.state.url)
  }

  render() {
    return (
      <div class="searchbar">
        <img class="logo" alt="sherlock gopher" src="img/sherlockgopher.png"></img>
        <h4 class="title">Sherlock Gopher</h4>
        <form onSubmit={this.handleSubmit} class="searchbarform" method="POST">
          <input class="searchbarinput" value={this.state.url} onChange={this.handleChange} name="url" placeholder="Hier könnte Ihre Werbung stehen..."></input>
          <button class="searchbarbutton" type="submit">Schnüffel</button>
        </form>
      </div>
    )
  }
}

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
