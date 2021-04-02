
//React Standard
import React from "react";
import Sidebar from "react-sidebar";
import 'bootstrap/dist/css/bootstrap.css';
import { Link } from 'react-router-dom'

//Stylesheets
import './assets/css/App.css'
import './assets/css/Sidebar.css'

//Images
import logo from './assets/img/sherlockgopher.png'

export default class SidebarMenu extends React.Component {
  
  constructor(props) {
    super(props);
      this.state = {
        sidebarOpen: false
      };
      this.onSetSidebarOpen = this.onSetSidebarOpen.bind(this);
  }

  /*
  handle the click to open the sidebar.
  */
  onSetSidebarOpen(open) {
    this.setState({ sidebarOpen: open });
  }
    
    render() {
        return (
          <div>
             <Sidebar
                sidebar={<div>
                  <div class="innerdiv">
                  <img class="menulogo" alt="sherlock gopher" src={logo}></img>
                  <hr class="line"></hr>
                  <p class="menubartitle">Sherlock Gopher</p>
                  <hr class="line"></hr>
                    <ul class="innerlist">
                      <li><Link class="innerlink" to="/">Home</Link></li>
                      <li><Link class="innerlink" to="/Graph">Graph</Link></li>
                      <li><Link class="innerlink" to="/sitesperformence">PoS</Link></li>
                      <li><Link class="innerlink" to="/nodedetails">Find a Node</Link></li>
                      <li><Link class="innerlink" to="/controls">Controls</Link></li>
                      <li><Link class="innerlink" to="/screenshots">Screenshots</Link></li>
                      <li><Link class="innerlink" to="/imagemetadata">Image Metadata</Link></li>
                    </ul>
                    <hr class="line"></hr>
                    <ul class="innerlist">
                      <li><Link class="innerlink" to="/scriptingapi">Scripting API</Link></li>
                    </ul>
                    <hr class="line"></hr>
                    <ul class="innerlist">
                      <li><a class="innerlink" href="https://github.com/DerAlexx/SherlockGopher"> Github </a></li>
                    </ul>
                  </div>
                </div>}
                open={this.state.sidebarOpen}
                onSetOpen={this.onSetSidebarOpen}
                styles={{ sidebar: { 
                  background: "#212f3d", 
                  zIndex:200, 
                  height:"100%", 
                  width:250,
                  color: "white",
                  position:"fixed",
                }
              }}
              >
              </Sidebar>
              <button class="btn btn-outline-light openbutton" onClick={() => this.onSetSidebarOpen(true)}>&#x2630; Menu </button>
          </div>
        );
      }
     
}