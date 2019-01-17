import React, { Component } from 'react';
import {
  Link,
  BrowserRouter,
  Redirect
} from "react-router-dom";
import './URLComponent.css';
import config from '../../config.json';

class URLComponent extends Component{
   
    constructor(){
        super();
        this.state = {
            hashValue : '',
            jsonData : [],
            fields :[],
            username : '',
            shouldRedirect : true,
            totalFhirResponse : {},
            fhirUrl : ''
        }
        this.handleChange = this.handleChange.bind(this);
        this.navigatePage = this.navigatePage.bind(this);
    }
    
    componentDidMount = () => {
        if(this.props.location.state === undefined){
            this.props.history.push({ 
                pathname : '/'
            });
        } else {          
            this.setState({username: this.props.location.state.username});
        }
        
    }

    handleChange(event) {
        this.setState({hashValue: event.target.value});
    }

    onClick = (event) =>{
        event.preventDefault();
        this.fetchHashValid();
    }

    fetchURL = (fhirUrl) => {
          
        let configParams = {
            method: 'GET',
            headers: {
              'Accept': 'application/json+fhir'
            },
          }
        fetch(fhirUrl, configParams)
            .then(response =>  response.json())
            .then((response) => {
                
                this.setState({ jsonData : response.entry , totalFhirResponse: response});
                var finalKeys =[];
                this.state.jsonData.forEach((eachEntry)=>{
                    finalKeys.push({"keyName" : eachEntry.resource.code.text, "category" : "unused", "value" :eachEntry.resource.text.div });
                    this.setState({fields: finalKeys});
                });
                this.navigatePage();

            });
        
    } 

    fetchData() {
        let configParams = {
            method: 'GET',
            headers: {
            'authorization': 'Bearer '+config.authToken,
            'content-Type': 'application/json'
            },
        }
      // eslint-disable-next-line
        fetch('http://'+config.hostIP+':'+config.port+'' + '/channels/'+config.channelName+'/chaincodes/'+config.chaincodeName+'?peer='+config.peerName+'&fcn=queryCustom&args=%5B%22%7B%5C%22selector%5C%22:%7B%5C%22_rev%5C%22:%5C%22'+this.state.hashValue+'%5C%22,%5C%22payerId%5C%22:%5C%22'+this.state.username+'%5C%22%7D%7D%22%5D', configParams)
            .then(response =>  response.json() )
            .then(response => {
                if (JSON.stringify(response) === '[]'){
                  
                } else {
                   
                    const fhirUrl = response.Record.fhirUrl;
                    this.fetchURL(fhirUrl);
                    this.setState({  fhirUrl: response.Record.fhirUrl })
                }
            } )
            .catch(err => console.log(err))
 
    }


    fetchHashValid(){
        let configParams = {
            method: 'GET',
            headers: {
              'authorization': 'Bearer '+config.authToken,
              'content-Type': 'application/json'
            },
          }
        // eslint-disable-next-line
        fetch('http://'+config.hostIP+':'+config.port+'' + '/channels/'+config.channelName+'/chaincodes/'+config.chaincodeName+'?peer='+config.peerName+'&fcn=isValid&args=%5B%22hash%22,%22'+this.state.hashValue+'%22%5D', configParams)
            .then(response =>  response.text() )
            .then(response => {
                if (response.length === 0 && response[0] !== 'E'){  
                    this.fetchData();
                    
                } else {
                    alert("Invalid Hash key")
                }
            } )
            .catch(err => console.log(err))
    }

    navigatePage =() =>{
        this.props.history.push({ 
            pathname : '/responseData',state : {fields :this.state.fields, hashValue : this.state.hashValue, username : this.state.username, setFields : []}
        });
    }

    render(){

        return(
            <BrowserRouter>
                
                    <div className="App">
                        
                            <div className = "title">
                                <span className = "header-text"> PAYER CHART REVIEW </span>
                            </div>         
                            
                            <form>
                                <div>
                                    <input type = "text" className = "inputBox" value = {this.state.hashValue} onChange={this.handleChange}  placeholder = "Enter Hash" required/>                 
                                    
                                    <button type = "submit"  className = "submitBox"  onClick = {(event) => this.onClick(event)} disabled = {!this.state.hashValue}>Submit</button>
                                    
                                </div>
                            </form> 

                        
                    </div>
                    
               

        </BrowserRouter>
        )
    }
}
export default URLComponent;