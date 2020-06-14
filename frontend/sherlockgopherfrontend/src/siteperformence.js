/*
React component showing the performence of all websites.
*/

//React
import React from "react";
import 'bootstrap/dist/css/bootstrap.css';

// Non-Standard
//import PERFORMENCE from './const.js';
import SearchBar from './searchbar.js';

//Stylesheets
import './assets/css/App.css'
import Axios from "axios";

export default class SitePerformence extends React.Component {

    PERFORMENCE = "http://0.0.0.0:8081/graph/v1/performenceofsites"

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
      isEmpty: false
    }
  
    constructor(props) {
      super(props);
      this.getCardStyle = this.getCardStyle.bind(this);
      this.getAmountOfEntrys = this.getAmountOfEntrys.bind(this);
      this.makeClickableURL = this.makeClickableURL.bind(this)
      this.interval = 0;
    }
    
    /*
    getAmountOfEntrys will count the Amount of entrys.
    */
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
        console.log(item.Statuscode)
        if (item.Statuscode >= 200 && item.Statuscode < 300) { //200
          zweihunderter++
        } else if (item.Statuscode >= 300 && item.Statuscode < 400) { //300
          dreithunderter++
        } else if (item.Statuscode >= 400 && item.Statuscode < 500) { //400
          vierhunderter++
        } else if (item.Statuscode >= 500 && item.Statuscode < 600) { //500
          fünfhunderter++
        } else if (item.Statuscode >= 100 && item.Statuscode < 200) { //100
          einhunderter++
        }  else {  // Errors
          errors++
        }
        entries++ 
        avrtt += (item.Responsetime * 1)
        return null
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
        this.interval = setInterval(async() => {
          Axios.get(this.PERFORMENCE).then(Response => {
            const chuncks = Response.data
            this.getAmountOfEntrys(chuncks)
            this.setState({
              items: chuncks,
              isEmpty: false
            })
          }).catch(errors => {
            console.log(errors)
            this.setState({
              items: [],
              isEmpty: true
            })
          })
        }, 1000)
      } catch (exception) {
        console.log(exception)
      }
    }

    componentWillUnmount() {
      clearInterval(this.interval)
    }
  
    /*
    getCardStyle will return the Style of a Card so e.g the site return 200 Ok so the card style will be alert-success.
    */
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

    /*
    makeClickableURL will turn a url in a clickable URL.
    */
    makeClickableURL(lurl) {
      if (lurl.startsWith("http://") || lurl.startsWith("https://")) {
        return lurl;
      } else {
        return "http://" + lurl
      }
    }
  
    render() {
      const {
        items, amountofentries, amountofHunderter, amountofZweiHunderter, amountofDreiHunderter, amountofVierHunderter, 
        amountofFünfHunderter, amountoferrors, averageRTT, isEmpty
      } = this.state
      return ( 
        <div>
          <SearchBar></SearchBar>
            <div class="body">
          <div id="tableElement" >
            <p>Metainformation</p>
            <span class="badge badge-pill badge-dark"> Entries: {amountofentries} </span> 
            <span style={{color: "#d5d8dc"}}>   </span>
            <span class="badge badge-pill badge-success"> 100er: {amountofHunderter}</span>
            <span style={{color: "#d5d8dc"}}>   </span>
            <span class="badge badge-pill badge-success"> 200er: {amountofZweiHunderter}</span>
            <span style={{color: "#d5d8dc"}}>   </span>
            <span class="badge badge-pill badge-warning"> 300er: {amountofDreiHunderter}</span>
            <span style={{color: "#d5d8dc"}}>   </span>
            <span class="badge badge-pill badge-danger"> 400er: {amountofVierHunderter}</span>
            <span style={{color: "#d5d8dc"}}>   </span>
            <span class="badge badge-pill badge-danger"> 500er: {amountofFünfHunderter}</span>
            <span style={{color: "#d5d8dc"}}>  </span>
            <span class="badge badge-pill badge-danger"> Errors: {amountoferrors}</span>
            <span style={{color: "#d5d8dc"}}>  </span>
            <span class="badge badge-pill badge-dark"> Ø RTT: {averageRTT}</span>
            <hr></hr>
            { isEmpty ? 
              <div class="alert alert-success">
                There a no entries in the database. 
              </div> :
                items.map(item => (
                  <div class={this.getCardStyle(item.Statuscode)} role="alert">
                  <h4 class="alert-heading"># <a class="alert-link" href={this.makeClickableURL(item.Address)}>{item.Address}</a></h4>
                  <p class="mb-0">
                    SherlockGopher gathered following information for address {item.Address}:<hr></hr>
                    <p class="font-weight-bold"><b>Responsetime:</b> {item.Responsetime} ms,</p> 
                    <p class="font-weight-bold"><b>Responsecode:</b> {item.Statuscode} (HTTP/s Standardcode)</p>
                  </p>
                </div>
                  ))
          }
          </div>
        </div>
        </div>
      )
    }
  }