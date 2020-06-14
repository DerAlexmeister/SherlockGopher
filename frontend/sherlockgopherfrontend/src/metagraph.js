//React standard
import React from 'react';
import Axios from 'axios';
import { PieChart } from 'react-chartkick'
import 'chart.js'

export default class ControlGraph extends React.Component {

    MONITOR = "http://localhost:8081/monitor/v1/meta"

    state = {
        // /monitor/v1/meta
        meta: [],
        metamessage: "",
        hasMetaError: false,

    }

    constructor(props) {
        super(props);

        this.getMetaGraph = this.getMetaGraph.bind(this);
        this.prepareMeta = this.prepareMeta.bind(this);
        this.getColorByStatus = this.getColorByStatus.bind(this);
        this.interval = 0;

    }

    componentDidMount(){
        this.interval = setInterval(() => {
            try {
                this.getMetaGraph()
            } catch(error) {
                console.log("cannot clear state.")
            }
        }, 1500)
    }

    componentWillUnmount() {
        clearInterval(this.interval)
    }

    /*
    Will fetch the meta information for the crawler and anaylser -> Status /monitor/v1/meta
    */
    getMetaGraph() {
        Axios.get(this.MONITOR).then( res => {
            const response = res.data;
            try {
                this.setState({
                    meta: response,
                    metamessage: "No error",
                    hasMetaError: false,
                })
            } catch (error) {
                this.setState({
                    metamessage: "For some Reason an Error occured. Cannot process response.",
                    hasMetaError: true,
                })  
            }
        }).catch(error => {
            this.setState({
                metamessage: "For some Reason an Error occured. Is the Webserver up?",
                hasMetaError: true,
            })
        })
    }

    /*
    perpareMeta will turn the json response into a format so that the pieChart can be rendered.
    */
    prepareMeta(meta) {
        try {
            var prep = []
            for (let key in meta) {
                if (key !== "Website") {
                    prep.push([key, meta[key]])
                }
            }
            return prep
        } catch (error) {

        }
    }

    /*
    Get the color of the pill for the status by the is actuall status. 
    */
    getColorByStatus(lstatus) {
        console.log("here", lstatus)
        if (lstatus === "Finished") {
            return "badge badge-pill badge-success"
        } else if (lstatus === "CrawlerError" || lstatus === "Failed") {
            return "badge badge-pill badge-danger"
        } else if (lstatus === "Undone") {
            return "badge badge-pill badge-warning"
        } else if (lstatus === "Processing") {
            return "badge badge-pill badge-info"
        } else if (lstatus === "SendToCrawler") {
            return "badge badge-pill badge-dark"
        } else if (lstatus === "Saving") {
            return "badge badge-pill badge-primary"
        } else {
            return "badge badge-pill badge-secondary"
        }
    }

    render() {
        const {
            meta,
            metamessage,
            hasMetaError
        } = this.state
        var readymetaAnalyser = this.prepareMeta(meta["Analyser"])
        var readymetaCrawler = this.prepareMeta(meta["Crawler"])
        console.log(readymetaAnalyser)
        return (
            <div>
                <hr></hr>
                <p>Analyser Tasks:</p>
                {
                    readymetaAnalyser.length > 0 ? readymetaAnalyser.map(item => (
                        <span class={this.getColorByStatus(item[0])} style={{marginRight: 10}}> {item[0]}:  {item[1]} </span>
                    )): null
                }
                <hr></hr>
                <p>Cralwer Tasks:</p>
                {
                    readymetaCrawler.length > 0 ? readymetaCrawler.map(item => (
                        <span class={this.getColorByStatus(item[0])} style={{marginRight: 10}}> {item[0]}:  {item[1]} </span>
                    )): null
                }
                <hr></hr>
                <br></br>
                { 
                    hasMetaError ? 
                            <div class="alert alert-danger">
                                { metamessage }
                            </div>
                            : 
                            <p>
                                <div>
                                    <PieChart colors={["#dc3545", "#28a745", "#17a2b8", "#ffc107"]} download={{background: "#fff"}} width="100%" height="150px" data={ readymetaCrawler } />  
                                    <br></br>
                                    <PieChart colors={["#dc3545", "#28a745", "#17a2b8","#007bff", "#343a40", "#ffc107"]} download={{background: "#fff"}} width="100%" height="250px" data={ readymetaAnalyser } />
                                </div>
                            </p>
                }
            </div>
        )
    }

}