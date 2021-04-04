import React from 'react';
import SearchBar from './searchbar.js';
import ReactPaginate from 'react-paginate';
import axios from 'axios'
import "./imagemetadataservice.css";

export default class Imagemetadataservice extends React.Component {
    METADATA = "http://0.0.0.0:8081/graph/v1/getmetadata"

    state = {
        postdata: [],
        currentPage: 0,
        maxPage: 1,
        pageRange: 0,
        sdatamessage: "",
        hasSdataError: false
    }

    constructor(props) {
        super(props);
        this.handlePageClick = this
            .handlePageClick
            .bind(this);
    }

    receivedData() {
        axios.get(this.METADATA + "/" + this.state.currentPage).then(res => {
            try {
                const data = res.data;
                const postData = data.map(pd => <React.Fragment>
                        <table>
                            <tr><th>Zustand</th><td>{pd.condition}</td></tr>
                            <tr><th>Datum</th><td>{pd.datetime_original}</td></tr>
                            <tr><th>Gerät</th><td>{pd.model}</td></tr>
                            <tr><th></th><td>{pd.make}</td></tr>
                            <tr><th>Beschreibung</th><td>{pd.maker_note}</td></tr>
                            <tr><th>Software</th><td>{pd.software}</td></tr>
                            <tr><th>GPS-Längengrad</th><td>{pd.gps_latitude}</td></tr>
                            <tr><th>GPS-Breitengrad</th><td>{pd.gps_longitude}</td></tr>
                        </table>
                        <br></br>
                        <br></br>
                    </React.Fragment>)
                this.setState({
                    sdatamessage: "No error",
                    hasSdataError: false,
                    postData,
                    currentPage: res.data.currentpage,
                    maxPage: res.data.maxpage,
                    pageRange: res.data.pagerange,
                })
            } catch (error) {
                console.log(error)
                this.setState({
                    sdatamessage: "For some Reason an Error occured. Cannot process response.",
                    hasMetaError: true,
                })  
            }
            }).catch(error => {
                console.log(error)
                this.setState({
                    sdatamessage: "For some Reason an Error occured. Is the Webserver up?",
                    hasSdataError: true,
            })
        })
    }

    handlePageClick = (e) => {
        const selectedPage = e.selected;
        const offset = selectedPage * this.state.perPage;
        this.setState({
            currentPage: selectedPage,
        }, () => {
            this.receivedData()
        });
    };

    componentDidMount(){
        this.interval = setInterval(() => {
            try {
                this.receivedData()
            } catch(error) {
                console.log(error)
            }
        }, 5000)
    }

    componentWillUnmount() {
        clearInterval(this.interval)
    }

    render() {
        const {
            sdatamessage,
            hasSdataError
        } = this.state
        return(
            <div>
                <SearchBar></SearchBar>
                <div class="body">
                    <p>Metadata of all Images</p>
                        { 
                            hasSdataError ? 
                                <div class="alert alert-danger">
                                    {sdatamessage}
                                </div>
                                :
                                <div>
                                    <hr></hr>
                                    <br></br>
                                    {this.state.postData}
                                    <ReactPaginate
                                        previousLabel={"<<"}
                                        nextLabel={">>"}
                                        breakLabel={"..."}
                                        breakClassName={"break-me"}
                                        pageCount={this.state.maxPage}
                                        marginPagesDisplayed={1}
                                        pageRangeDisplayed={this.state.pageRange}
                                        onPageChange={this.handlePageClick}
                                        containerClassName={"pagination"}
                                        subContainerClassName={"pages pagination"}
                                        activeClassName={"active"} 
                                    /> 
                                </div>     
                        }            
                </div>
            </div>
        )
    }
}