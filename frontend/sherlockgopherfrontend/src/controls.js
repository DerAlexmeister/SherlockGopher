//React standard
import React from 'react';
import Axios from 'axios';

//Javascript
import SearchBar from './searchbar.js';
import ControlGraph from './metagraph.js'

export default class Controls extends React.Component {

    DROPPGENDPOINT = "http://0.0.0.0:8081/controller/v1/droppg"
    DROPMGENDPOINT = "http://0.0.0.0:8081/controller/v1/dropmg"
    DROPITENDPOINT = "http://0.0.0.0:8081/controller/v1/dropit"
    CHANGESTATUS = "http://0.0.0.0:8081/controller/v1/changestate"
    STATUSENDPOINT = "http://0.0.0.0:8081/controller/v1/status"
    AREYOUTHERE = "http://0.0.0.0:8081/areyouthere"

    state = {
        message: "",
        haserror: false,
        showmessage: false,

        //status
        status: [],
        statusmessage: "",
        hasstatuserror: false,

        //changestatus
        changestatusmessage: "",
        changestatushasMessage: false,
        changestatushasErrorMessage: false,

    }

    constructor(props) {
        super(props);

        this.dropPg = this.dropPg.bind(this);
        this.dropMg = this.dropMg.bind(this);
        this.dropIt = this.dropIt.bind(this);
        this.areyouthere = this.areyouthere.bind(this);
        this.submitStatusToWebserver = this.submitStatusToWebserver.bind(this);
        this.getColorByStatus = this.getColorByStatus.bind(this);
        this.interval = 0;
    }

    componentDidMount() {
        this.getStatus()
        this.interval = setInterval(() => {
            try {
                this.getStatus()
                this.setState({
                    message: undefined,
                    haserror: false,
                    showmessage: false,
                    changestatushasMessage: false,
                    changestatushasErrorMessage: false,
                })
            } catch(error) {
                console.log("cannot clear state.")
            }
        }, 4500)
    }

    /*
    Will unmount any component at a given point.
    */
    componentWillUnmount() {
        clearInterval(this.interval)
    }

    /*
    Will fetch the status of the services.
    */
    getStatus() {
        try{
            Axios.get(this.STATUSENDPOINT).then( res => {
                const response = res.data;
                this.setState({
                    status: [response],
                    statusmessage: "",
                    hasstatuserror: false,
                })
            }).catch(error => {
                this.setState({
                    status: [],
                    statusmessage: "Can not retrieve the status. Is the Webserver up?",
                    hasstatuserror: true,
                })
            })
        } catch (error) {

        }
    }



    /*
    dropIt will send a droptable to the webserver which will drop the table.
    */
    dropIt() {
        Axios.get(this.DROPITENDPOINT).then( res => {
            const response = res.data;
            try {
                this.setState({
                    message: response.Message,
                    haserror: false,
                    showmessage: true,
                })
            } catch (error) {
                this.setState({
                    message: "Cannot read the response for DropTable",
                    haserror: true,
                    showmessage: true,
                })
            }
        }).catch(error => {
            this.setState({
                message: "An error occured while trying to drop the table. Is the Database online aswell as the Webserver?",
                haserror: true,
                showmessage: true,
            })
            console.log(error)
        })
    }

    /*
    dropMg will send a droptable to the webserver which will drop the table.
    */
    dropMg() {
        Axios.get(this.DROPMGENDPOINT).then( res => {
            const response = res.data;
            try {
                this.setState({
                    message: response.Message,
                    haserror: false,
                    showmessage: true,
                })
            } catch (error) {
                this.setState({
                    message: "Cannot read the response for DropTable",
                    haserror: true,
                    showmessage: true,
                })
            }
        }).catch(error => {
            this.setState({
                message: "An error occured while trying to drop the table. Is the Database online aswell as the Webserver?",
                haserror: true,
                showmessage: true,
            })
            console.log(error)
        })
    }

    /*
    dropPg will send a droptable to the webserver which will drop the table.
    */
    dropPg() {
        Axios.get(this.DROPPGENDPOINT).then( res => {
            const response = res.data;
            try {
                this.setState({
                    message: response.Message,
                    haserror: false,
                    showmessage: true,
                })
            } catch (error) {
                this.setState({
                    message: "Cannot read the response for DropTable",
                    haserror: true,
                    showmessage: true,
                })
            }
        }).catch(error => {
            this.setState({
                message: "An error occured while trying to drop the table. Is the Database online aswell as the Webserver?",
                haserror: true,
                showmessage: true,
            })
            console.log(error)
        })
    }

    /*
    areyouthere will send a areyouthere to the webserver which will response incase the webserver is online.
    */
    areyouthere() {
        console.log("dropit");
        Axios.get(this.AREYOUTHERE).then(res => {
            const response = res.data;
            try {
                this.setState({
                    message: response.Message + " 🤩",
                    haserror: false,
                    showmessage: true,
                })
            } catch (error) {
                this.setState({
                    message: "Cannot read the response for are you there",
                    haserror: true,
                    showmessage: true,
                })
            }
        }).catch(error => {
            this.setState({
                message: "Seems like the webserver is down. 😭",
                haserror: true,
                showmessage: true,
            })
        })
    }

    /*
    Submit the Status changes to the Webserver.
    */
    submitStatusToWebserver(target, operation) {
        try{
            Axios.post(this.CHANGESTATUS,
                JSON.stringify({
                "target": target,
                "operation": operation
                })).then(response => {
                    const chuncks = response.data
                    if (chuncks.Status !== undefined) {
                        this.setState({
                            changestatusmessage: chuncks.Status,
                            changestatushasMessage: true,
                            changestatushasErrorMessage: false,
                        })
                    } else {
                        this.setState({
                            changestatusmessage: "Status of the Crawler or Analyser is unknown!",
                            changestatushasMessage: true,
                            changestatushasErrorMessage: true,
                        })
                    }
                }).catch(error => {
                    this.setState({
                        changestatusmessage: "An Erroc occured. Error: " + error,
                        changestatushasMessage: true,
                        changestatushasErrorMessage: true,
                    })
                });
        } catch(error) {
            console.log("An error occured: " + error)
        }
    }


    /*
    Get the color of the pill for the status by the is actuall status.
    */
    getColorByStatus(lstatus) {
        if (lstatus === "Running") {
            return "badge badge-pill badge-success"
        } else if (lstatus === "Stop") {
            return "badge badge-pill badge-danger"
        } else if (lstatus === "Paused") {
            return "badge badge-pill badge-warning"
        } else if (lstatus === "Idle") {
            return "badge badge-pill badge-info"
        } else {
            return "badge badge-pill badge-dark"
        }
    }

    render() {
        const {
            message,
            haserror,
            showmessage,
            status,
            statusmessage,
            hasstatuserror,
            changestatusmessage,
            changestatushasMessage,
            changestatushasErrorMessage,
        } = this.state
        var AnalyserStatus = status.map(item => item["Analyser"])
        var CrawlerStatus = status.map(item => item["Crawler"])
        return (
          <div>
            <SearchBar></SearchBar>
                <div class="body">
                    <p>Controls</p>
                    {hasstatuserror ?
                        <div class="alert alert-danger"> {statusmessage} </div>
                        :
                        <div>
                            <span class={this.getColorByStatus(AnalyserStatus[0])}> Analyser:  { AnalyserStatus }</span>
                            <span style={{color: "#d5d8dc"}}>   </span>
                            <span class={this.getColorByStatus(CrawlerStatus[0])}> Crawler: { CrawlerStatus }</span>
                            <span style={{color: "#d5d8dc"}}>   </span>
                        </div>
                    }
                    <hr></hr>
                    <p>
                        Current Systemload
                        <ControlGraph></ControlGraph>
                    </p>
                    <hr></hr>
                    <p>
                        Servicemanager<br></br><br></br>
                            {
                                //General Messages
                            showmessage ?
                            <div style={{textAlign: "center"}} class={haserror ? "alert alert-danger" : "alert alert-success"} role="alert">
                                {message}
                            </div>
                            : null
                            }

                            {
                                // Change Status Message.
                            changestatushasMessage ?
                            <div style={{textAlign: "center"}} class={changestatushasErrorMessage ? "alert alert-danger" : "alert alert-success"} role="alert">
                                {changestatusmessage}
                            </div>
                            : null
                            }
                        <div class="alert alert-dark" role="alert">
                            <h5 class="alert-heading">General</h5>
                            <p>Control all services of Sherlock Gopher.</p>
                            <hr></hr>
                            <button onClick={() => this.submitStatusToWebserver("All", "Stop")} type="button" class="btn btn-danger">Stop</button>
                            <span style={{color: "#d5d8dc"}}>   </span>
                            <button onClick={() => this.submitStatusToWebserver("All", "Pause")} type="button" class="btn btn-primary">Pause</button>
                            <span style={{color: "#d5d8dc"}}>   </span>
                            <button onClick={() => this.submitStatusToWebserver("All", "Resume")} type="button" class="btn btn-success">Resume</button>
                            <span style={{color: "#d5d8dc"}}>   </span>
                            <hr></hr>
                            <p>Clear the Queue of all services</p>
                            <button onClick={() => this.submitStatusToWebserver("All", "Clean")} type="button" class="btn btn-info">Clear Queue</button>
                            <span style={{color: "#d5d8dc"}}>   </span>
                            <hr></hr>
                            <p>Check with one click whether or not the webserver is alive.</p>
                            <button onClick={this.areyouthere} type="submit" class="btn btn-secondary">Are you there?</button>
                            <hr></hr>
                        </div>
                        <hr></hr>
                        <div class="alert alert-dark" role="alert">
                            <h5 class="alert-heading">Database Management</h5>
                            <p>Delete the current Neo4J-Database.</p>
                            <button onClick={this.dropIt} type="submit" class="btn btn-danger">Drop it!</button>
                            <hr></hr>
                            <p>Delete the current MongoDB-Database.</p>
                            <button onClick={this.dropMg} type="submit" class="btn btn-danger">Drop it!</button>
                            <hr></hr>
                            <p>Delete the current Postgres-Database.</p>
                            <button onClick={this.dropMg} type="submit" class="btn btn-danger">Drop it!</button>
                        </div>
                        <hr></hr>
                        <div class="alert alert-dark" role="alert">
                            <h5 class="alert-heading">Analyser</h5>
                            <p>Control the Analyser.</p>
                            <hr></hr>
                            <button onClick={() => this.submitStatusToWebserver("Analyser", "Stop")} type="button" class="btn btn-danger">Stop</button>
                            <span style={{color: "#d5d8dc"}}>   </span>
                            <button onClick={() => this.submitStatusToWebserver("Analyser", "Pause")} type="button" class="btn btn-primary">Pause</button>
                            <span style={{color: "#d5d8dc"}}>   </span>
                            <button onClick={() => this.submitStatusToWebserver("Analyser", "Resume")} type="button" class="btn btn-success">Resume</button>
                            <span style={{color: "#d5d8dc"}}>   </span>
                            <hr></hr>
                            <p>Clear the Queue of the Analyser-Service</p>
                            <button onClick={() => this.submitStatusToWebserver("Analyser", "Clean")} type="button" class="btn btn-info">Clear Queue</button>
                            <span style={{color: "#d5d8dc"}}>   </span>
                        </div>
                        <hr></hr>
                        <div class="alert alert-dark" role="alert">
                            <h5 class="alert-heading">Crawler</h5>
                            <p>Control the Crawler.</p>
                            <hr></hr>
                            <button onClick={() => this.submitStatusToWebserver("Crawler", "Stop")} type="button" class="btn btn-danger">Stop</button>
                            <span style={{color: "#d5d8dc"}}>   </span>
                            <button onClick={() => this.submitStatusToWebserver("Crawler", "Pause")} type="button" class="btn btn-primary">Pause</button>
                            <span style={{color: "#d5d8dc"}}>   </span>
                            <button onClick={() => this.submitStatusToWebserver("Crawler", "Resume")} type="button" class="btn btn-success">Resume</button>
                            <span style={{color: "#d5d8dc"}}>   </span>
                            <hr></hr>
                            <p>Clear the Queue of the Crawler-Service</p>
                            <button onClick={() => this.submitStatusToWebserver("Crawler", "Clean")} type="button" class="btn btn-info">Clear Queue</button>
                            <span style={{color: "#d5d8dc"}}>   </span>
                        </div>
                    </p>
                </div>
          </div>
        )
      }
}
