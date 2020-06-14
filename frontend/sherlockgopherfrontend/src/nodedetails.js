/*
Will represent the details about a node.
*/

import 'bootstrap/dist/css/bootstrap.css';
import React from "react";
import axios from 'axios';

import SearchBar from './searchbar.js';

import './assets/css/nodedetails.css'

export default class Nodedetails extends React.Component {

    state = {
        value: "",
        reponsegoterror: false,
        items: [],
        showComponent: false,
        searchurl: undefined,
    }
    
    constructor(props) {
        super(props);
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.serverRequest = this.serverRequest.bind(this);
        this.getPrettyResponse = this.getPrettyResponse.bind(this);
        this.makeClickableURL = this.makeClickableURL.bind(this);
        this.isVerified = this.isVerified.bind(this);
        this.getAdvancedNodeInformation = this.getAdvancedNodeInformation.bind(this);
    }
    
    /*
    Will make the Request to the webserver.
    */
    serverRequest(submiturl) {
        axios.post("http://0.0.0.0:8081/graph/v1/detailsofnode",
          JSON.stringify({
            url: submiturl
        })).then(response => {
            const chuncks = response.data
            if (chuncks.Message !== undefined) {
                this.setState({
                    showComponent: true,
                    reponsegoterror: true,
                    items: chuncks,
                    searchurl: undefined
                })
            } else {
                this.setState({
                    searchurl: submiturl, 
                    showComponent: true,
                    reponsegoterror: false,
                    items: chuncks
                })
            }
        }).catch(response => { //TODO
            const chuncks = response.data
            console.log(response)
            this.setState({
                searchurl: undefined,
                showComponent: true,
                reponsegoterror: true,
                items: chuncks
            })
        });
    }
    
    /*
    handleChange will handle the Change of the input-field
    */
    handleChange(event) {   
        this.setState(
          { 
            value: event.target.value,
            showComponent: false,
            reponsegoterror: false,
            items: undefined
          }
        );  
    }
    
    /*
    handleSubmit will handle the submitted url.
    */
    handleSubmit(event) {
        event.preventDefault();
        this.serverRequest(this.state.value)
    }

    /*
    isVerified will check whether or not a node is verified and return the color for the badge.
    */
    isVerified(status) {
        //TODO unused
        if (status !== 'verified') {
            return "color: red"
        } else {
            return "color: green" 
        }
    }

    /*
    makeClickableURL will turn a URL in a clickable URL.
    */
    makeClickableURL(lurl) {
        if (lurl.startsWith("http://") || lurl.startsWith("https://")) {
          return lurl;
        } else {
          return "http://" + lurl
        }
      }

    /*
    getAdvancedNodeInformation will turn the given details about a node into a fitting-format to render it.
    */
    getAdvancedNodeInformation(instance) {
        try {
            let arr = [] 
            for (let key in instance) {
                if (key !== "Responsetime" && key !== 'Status' && key !== 'Statuscode' && key !== 'Filetype' && key !== "Address") {
                    arr.push([key, instance[key]])
                }   
            }
            return arr
        } catch (error) {
            return null
        }
    }

    /*
    getPrettyResponse will prettify the response.
    */
    getPrettyResponse() {
        if (this.state.reponsegoterror) {         
            return (
                    <div class="alert alert-warning">
                       {this.state.items.Message}
                    </div>
                )
        } else { 
            if (this.state.searchurl !== undefined) {
                const { items, searchurl } = this.state
                var entries = this.getAdvancedNodeInformation(items[String(searchurl)])
                console.log(entries)
                try {
                    return (
                        <div class="alert alert-secondary" role="alert">
                            <h4 class="alert-heading">{items[String(searchurl)].Address}</h4>
                            <hr></hr>
                                Metainformation<br></br>
                                <ul>
                                    {items[String(searchurl)].Status !== undefined ? <li> <b>Status</b>: {items[String(searchurl)].Status}</li> : null}
                                    {items[String(searchurl)].Responsetime !== undefined ? <li> <b>RTT</b>: {items[String(searchurl)].Responsetime}</li> : null}
                                    {items[String(searchurl)].Statuscode !== undefined ? <li> <b>HTTP/s Responsecode:</b> {items[String(searchurl)].Statuscode}</li> : null}
                                    {items[String(searchurl)].Filetype !== undefined ? <li> <b>Resourcetyp</b>: {items[String(searchurl)].Filetype}</li> : null}
                                </ul>
                            <hr></hr>
                                Advanced information
                                    <ul>
                                        {
                                            entries.length > 0 ? entries.map(item => (
                                                <li><b>{item[0]}</b>: {item[1]}</li>
                                            )): null
                                        }
                                    </ul>
                            <hr></hr>
                                <p class="mb-0">
                                    <a type="button" href={this.makeClickableURL(searchurl)} class="btn btn-outline-dark">Visit this page</a>
                                    <span style={{color: "#d5d8dc"}}>   </span>
                                </p>
                        </div>
                    )
                } catch (error) {
                    console.log(error)
                    return (
                        <div class="alert alert-danger" role="alert">
                            A problem occured with the following address: <b>{searchurl}</b>.
                        </div>
                    )
                }   
            } else {
                return (
                    <div class="alert alert-danger" role="alert">
                        An error occured! Please try to search for a different node!
                    </div>
                )
            }
        }
    }
    
    render() {
        return (
            <div>
                <SearchBar></SearchBar>
                <div class="body">
                    <p>Node</p>  
                    <hr></hr>
                    <form onSubmit={this.handleSubmit}>
                        <div  class="form-group">
                            <label for="exampleInputEmail1">Enter a address and find the right node</label><br></br>
                            <p class="informationForInput">The search is case-sensitive so a space or wrong character can lead to an empty result!</p>
                            <input class="form-control" value={this.state.value} onChange={this.handleChange} name="url" placeholder="Enter the address" required></input>
                        </div>
                        <button type="submit" class="btn btn-outline-dark btn-sm">search</button>
                    </form>
                    <hr></hr>
                    {this.state.showComponent ? this.getPrettyResponse() : null}
                </div>
         </div>
        )
    }
}