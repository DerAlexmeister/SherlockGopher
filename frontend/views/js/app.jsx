// Basic RESTendpoint.
const RESTENDPOINT = "http://0.0.0.0:8081";

// Basicaddress of the graph endpoint.
const GRAPHGROUP = "/graph/v1";

// Address to post a url to search for.
const SEARCHENDPOINT = RESTENDPOINT + GRAPHGROUP + "/search";

const PERFORMENCE = RESTENDPOINT + GRAPHGROUP + "/performenceofsites"

//Site will render a table 
class SitePerformence extends React.Component {

  state = {
    items: [],
  }

  constructor(props) {
    super(props);
    this.getCardStyle = this.getCardStyle.bind(this);
  }

  componentDidMount() {
    try {
      setInterval(async() => {
        const res = await fetch(PERFORMENCE)
        const chuncks = await res.json()
        //this.makeRows(chuncks) //await or something and then turn it into a table
        this.setState({items: chuncks})
      }, 1000)
    } catch (exception) {
      console.log(exception)
    }
  }

  getCardStyle(status) {
    if (status > 199 && status < 300) { //200
      return "alert alert-success"
    } else if (status > 299 && status < 400) { //300
      return "alert alert-warning"
    } else if (status > 399 && status < 500) { //400
      return "alert alert-danger"
    } else if (status > 499 && status < 600) { //500
      return "alert alert-danger"
    } else if (status > 99 && status < 200) { //100
      return "alert alert-success"
    }  else {  // Errors
      return "alert alert-info"
    }
  }

  render() {
    const {items} = this.state
    return ( 
      <div id="tableElement" >
        {items.map(item => (
          <div class={this.getCardStyle(item.Status)} role="alert">
          <h4 class="alert-heading">#<a class="alert-link" href={item.Address}> {item.Address} </a></h4>
          <p class="mb-0">
            SherlockGopher measured for {item.Address} the following information: <hr></hr>
            <p class="font-weight-bold"><b>Responsetime:</b> {item.ResponseTime} ms,</p> 
            <p class="font-weight-bold"><b>Responsecode:</b> {item.Status} (HTTP/s Standartcode)</p>
          </p>
        </div>
          ))}
      </div>
    )
  }
}

// Just the Searchbar Component with the logo, title and the searchbar to submit a url.
class SearchBar extends React.Component {

  state = {
    value: "",
    lresponse: undefined
  }

  // CTor for the Searchbar Component.
  constructor(props) {
    super(props);
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.serverRequest = this.serverRequest.bind(this);
  }

  // Making the POST REQUEST
  serverRequest(submiturl) {
    $.post(
      SEARCHENDPOINT,
      JSON.stringify({
        url: submiturl
      }),
      response => {
        console.log(response)
      }, 'json'
    );
      console.log(this.state.lresponse)
    return
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
    this.serverRequest(this.state.value)
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
            <div style = {{position: 'absolute', top:80,left: 50, right:50}}>
              <SitePerformence></SitePerformence>
            </div>
        </div>
      )
    }
}

ReactDOM.render(<App />, document.getElementById('app'));
