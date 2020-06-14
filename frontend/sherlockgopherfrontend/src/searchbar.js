/*
Searchbar Represents the menubanner at the top with a sidebar menu.
*/

//React Standard
import 'bootstrap/dist/css/bootstrap.css';
import React from "react";

//Stylesheets
import './assets/css/App.css'

// Javascript
import SidebarMenu from './sidebar.js';

//Images
import logo from './assets/img/sherlockgopher.png'

export default class SearchBar extends React.Component {
 
    render() {
      return (       
          <div style={{width:"100%",height:'100%', position: "absolute", top:"0px", left: "0px"}}>
            <SidebarMenu></SidebarMenu>
            <div class="searchbar" style={{zIndex: 10}}>
              <h4 class="title">SherlockGopher</h4>
            </div>
            <a href="/"><img class="logo" alt="sherlock gopher" src={logo}></img></a>
          </div>
      )
    }
  }