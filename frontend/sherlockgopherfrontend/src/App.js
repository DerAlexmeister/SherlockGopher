/*
App represents the entrypoint to the frontend.
*/

//React Standard
import 'bootstrap/dist/css/bootstrap.css';
import React from "react";
import axios from 'axios';

//Javascripts
import SearchBar from './searchbar.js';

//Stylesheets
import './assets/css/App.css';

//Images
import logo from './assets/img/sherlockgopher.png'

// The actual APP Component which will be rendered by calling the Website.
export default class App extends React.Component {

  SEARCHENDPOINT = "http://0.0.0.0:8081/graph/v1/search"
  
  state = {
    showComponent: 0,
    value: "", 
    message: "",
    showMessage: false,
    isErrorMessage: false,
  }
  
  constructor(props) {
    super(props);
    this.handleClick = this.handleClick.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.serverRequest = this.serverRequest.bind(this);

    this.interval = 0;
  }

  componentDidMount() {
    this.interval = setInterval(async () => {
      try {
        this.setState({
          showMessage: false,
        })
      } catch (error) {
          console.log("An error occured while trying to get the Metadata.")
          }
    }, 5000) 
  }

  componentWillUnmount() {
    clearInterval(this.interval)
  }

  /*
  handleClick will handle a click on the button to search.
  */
  handleClick(param) {
    this.setState({
      showComponent:param
    })
  }

  /*
  serverRequest will make a search request and handle the response.
  */
  serverRequest(submiturl) {
    axios.post(this.SEARCHENDPOINT,
      JSON.stringify({
        url: submiturl
    })).then(res => {
          const response = res.data 
          try {
            if (response.Message !== undefined) {
              this.setState({
                message: response.Message,
                showMessage: true,
                isErrorMessage: true,
                value: "",
              })
            } else {
              if (response.Status !== undefined) {
                this.setState({
                  message: "Submitted URL to the Crawler.",
                  showMessage: true,
                  isErrorMessage: false,
                  value: "",
                })
              } else {
                this.setState({
                  message: "Url Submitted",
                  showMessage: true,
                  isErrorMessage: false,
                  value: "",
                })
              }
            }
          } catch (error) {
            this.setState({
              message: "An error occured while trying to process the response. Reload the page and try again!",
              showMessage: true,
              isErrorMessage: true,
              value: "",
            })
          }
    }).catch(response => {
        this.setState({
          message: "An error occured while trying to process your request. Maybe the URL is malformed!",
          showMessage: true,
          isErrorMessage: true,
          value: "",
        })
    });
  }

  /*
  handleChange handles the change of the search-field.
  */
  handleChange(event) {   
    this.setState(
      { 
        value: event.target.value,
        message: undefined,
        showMessage: false,
        isErrorMessage: false,
      }
    );  
  }

  /*
  handleSubmit will work with the submitted value.
  */
  handleSubmit(event) {
    event.preventDefault();
    this.serverRequest(this.state.value)
  }

  render() {
    const {message, showMessage, isErrorMessage} = this.state
    return (
      <div>
          <SearchBar></SearchBar>
          <div class="body">
            {
              showMessage ? 
              <div style={{textAlign: "center"}} class={isErrorMessage ? "alert alert-danger" : "alert alert-success"} role="alert">
                {message}
              </div>
              : null
            }
            <img class="bigLogo" alt="sherlock gopher" src={logo}></img>
              <form onSubmit={this.handleSubmit} class="searchbarform" >
                <input class="searchbarinput" value={this.state.value} onChange={this.handleChange} name="url" placeholder="Hier könnte Ihre Werbung stehen!"></input>
                <p>
                  <br></br>
                  <button class="btn btn-dark searchbarbutton" type="submit">Schnüffel</button>                  
                </p>
              </form>
          </div>
      </div>
    )
  }
}