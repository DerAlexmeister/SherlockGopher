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
    amountofentries: 0,
    amountofHunderter: 0,
    amountofZweiHunderter: 0,
    amountofDreiHunderter: 0,
    amountofVierHunderter:0,
    amountofFünfHunderter:0,
    amountoferrors: 0,
    averageRTT: 0,
  }

  constructor(props) {
    super(props);
    this.getCardStyle = this.getCardStyle.bind(this);
    this.getAmountOfEntrys = this.getAmountOfEntrys.bind(this);
  }
  
  getAmountOfEntrys(chuncks) {
    var entries = 0
    var einhunderter = 0
    var zweihunderter = 0
    var dreithunderter = 0
    var vierhunderter = 0
    var fünfhunderter = 0
    var errors = 0
    var avrtt = 0
    chuncks.map(item => {
      if (item.Status > 199 && item.Status < 300) { //200
        zweihunderter++
      } else if (item.Status > 299 && item.Status < 400) { //300
        dreithunderter++
      } else if (item.Status > 399 && item.Status < 500) { //400
        vierhunderter++
      } else if (item.Status > 499 && item.Status < 600) { //500
        fünfhunderter++
      } else if (item.Status > 99 && item.Status < 200) { //100
        einhunderter++
      }  else {  // Errors
        errors++
      }
      entries++ 
      avrtt += (item.ResponseTime * 1)
    },
    )
    avrtt = Math.floor(avrtt/entries)
    this.setState({
      amountofentries: entries, amountofHunderter: einhunderter, amountofZweiHunderter: zweihunderter, 
      amountofDreiHunderter: dreithunderter, amountofVierHunderter: vierhunderter, amountofFünfHunderter: fünfhunderter, 
      amountoferrors: errors, averageRTT: avrtt
    })
  }

  componentDidMount() {
    try {
      setInterval(async() => {
        const res = await fetch(PERFORMENCE)
        const chuncks = await res.json()
        this.getAmountOfEntrys(chuncks)
        this.setState({items: chuncks})
      }, 500)
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
    const {
      items, amountofentries, amountofHunderter, amountofZweiHunderter, amountofDreiHunderter, amountofVierHunderter, 
      amountofFünfHunderter, amountoferrors, averageRTT
    } = this.state
    return ( 
      <div id="tableElement" >
        <span class="badge badge-pill badge-primary"> Entries: {amountofentries} </span> 
        <span style={{color: "#d5d8dc"}}>   </span>
        <span class="badge badge-pill badge-primary"> 100er: {amountofHunderter}</span>
        <span style={{color: "#d5d8dc"}}>   </span>
        <span class="badge badge-pill badge-primary"> 200er: {amountofZweiHunderter}</span>
        <span style={{color: "#d5d8dc"}}>   </span>
        <span class="badge badge-pill badge-primary"> 300er: {amountofDreiHunderter}</span>
        <span style={{color: "#d5d8dc"}}>   </span>
        <span class="badge badge-pill badge-primary"> 400er: {amountofVierHunderter}</span>
        <span style={{color: "#d5d8dc"}}>   </span>
        <span class="badge badge-pill badge-primary"> 500er: {amountofFünfHunderter}</span>
        <span style={{color: "#d5d8dc"}}>  </span>
        <span class="badge badge-pill badge-primary"> Errors: {amountoferrors}</span>
        <span style={{color: "#d5d8dc"}}>  </span>
        <span class="badge badge-pill badge-primary"> Ø RTT: {averageRTT}</span>
        <hr></hr>
        {items.map(item => (
          <div class={this.getCardStyle(item.Status)} role="alert">
          <h4 class="alert-heading"># <a class="alert-link" href={item.Address}>{item.Address}</a></h4>
          <p class="mb-0">
            SherlockGopher gathered following information for address {item.Address}:<hr></hr>
            <p class="font-weight-bold"><b>Responsetime:</b> {item.ResponseTime} ms,</p> 
            <p class="font-weight-bold"><b>Responsecode:</b> {item.Status} (HTTP/s Standardcode)</p>
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
            <div style = {{position: 'absolute', top:65,left: 50, right:50}}>
              <SitePerformence></SitePerformence>
            </div>
        </div>
      )
    }
}

ReactDOM.render(<App />, document.getElementById('app'));
