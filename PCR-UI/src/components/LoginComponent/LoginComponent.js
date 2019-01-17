import React, { Component } from 'react';
import {
    Link,
    BrowserRouter
  } from "react-router-dom";
import './LoginComponent.css';
import config  from '../../config.json';
import img from '../../images/logo.png';

class LoginComponent extends Component{

    constructor(){
        super();
        this.state = {
            username : '',
            password : '',
            config : {hostIP : config.hostIP, 
                      port : config.port, 
                      channelName : config.channelName, 
                      chaincodeName : config.chaincodeName,
                      peerName : config.peerName,
                      authToken : config.authToken
                    },
                    isSuccessfulLogin : false
        }
        
        this.handleChange = this.handleChange.bind(this);
        this.handlePasswordChange = this.handlePasswordChange.bind(this);
        this.fetchUser = this.fetchUser.bind(this);
        this.navigateChainCode = this.navigateChainCode.bind(this);
    }

    componentDidMount = () => {
    }

    handleChange(event) {
        this.setState({username: event.target.value});
    }

    handlePasswordChange(event) {
        this.setState({password: event.target.value});
    }

    onClick =(event) =>{
        event.preventDefault();
        this.fetchUser(); 
        
    }

    
    fetchUser = () => {
        let configParams = {
            method: 'GET',
            headers: {
                'authorization': 'Bearer ' + config.authToken,
                'content-Type': 'application/json'
            },
        }
        
        
        let url = 'http://'+config.hostIP+':'+config.port+'' + '/channels/'+config.channelName+'/chaincodes/'+config.chaincodeName+'?peer='+config.peerName+'&fcn=isValid&args=%5B%22payer%22,%22'+this.state.username+'%22%5D'
        
        fetch(url, configParams)
            .then(response => response.text())
            .then(response => {
               
                if (response.length === 0 && response[0] !== 'E') {
                    this.navigateChainCode();
                } else {
                    
                    alert("User did not registerd for Blockchain Service")
                }
            })
            .catch(err => console.log(err))        
    }

    navigateChainCode = () =>{
        this.props.history.push({ 
            pathname : '/responseData',state : {username :this.state.username, setFields : []}
        }); 
    }
    
    render(){

        return(
            <BrowserRouter>
                <div className = "login-main">
                    <div className = "login-header">
                        <img src= {img} alt="Cerner Logo" style={{width: '200px' , float : 'left', height: '65px'}}/>
                        <span className = "header-text"> PAYER CHART REVIEW </span>
                    </div>
                    <div className = "login-input-container">
                        
                        <form>  
                            <label className = "required-label-space">
                                <span> LOGIN </span> 
                            </label>                         
                            <label className = "required-spacing">
                                <input type = "text" className = "username-box" value = {this.state.username} onChange={this.handleChange}  placeholder = "Enter Username" required/> 
                            </label>
                            
                            <label className = "required-spacing">
                                <input type = "password" className = "username-box" value = {this.state.password} onChange={this.handlePasswordChange}  placeholder = "Enter Password" required/>                                                                                     
                            </label>

                            <label className = "required-spacing">
                                <button type = "submit"  className = "submit-Box"  onClick = {(event) => this.onClick(event)} disabled = {!this.state.password && !this.state.username}>LOGIN</button>
                            </label>
                        </form>
                        
                        
                    </div>
                
                </div>
            
            </BrowserRouter>
        )
    }  
};

export default LoginComponent;