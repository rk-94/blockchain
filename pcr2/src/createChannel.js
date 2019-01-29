import React, { Component } from 'react';
import Button from 'terra-button/lib/Button';
import netConfig from './config';

class Channel extends Component {
    constructor(props){
        super(props);
        this.state ={
            message: "", auth1: '', auth2: '', chaincodeName: 'mycc', channelName: 'mychannel'

        }
    }
    registerUser1(){
        let config = {
            method: 'POST',
            headers: {
              'content-Type': 'application/x-www-form-urlencoded'
              
            },
            body: 'username=Cerner&orgName=Org1'
          }
         
      
          fetch('http://'+netConfig.hostIP+':'+netConfig.port+'' + '/users', config)
            .then(response =>  response.json() )
            .then((response) => {
                if (response.success === true){
                    this.setState({ auth1: response.token})
                        
                          
                    }
                this.registerUser2();
                }
            );
    }
    registerUser2(){
        let config = {
            method: 'POST',
            headers: {
              'content-Type': 'application/x-www-form-urlencoded'
              
            },
            body: 'username=Cerner&orgName=Org2'
          }
         
      
          fetch('http://'+netConfig.hostIP+':'+netConfig.port+'' + '/users', config)
            .then(response =>  response.json() )
            .then((response) => {
                if (response.success === true){
                    this.setState({ auth2: response.token})
                        
                          
                    }
                this.craeteChannel();

                }
            );
    }
    /*
    craeteChannel() {
        let config = {
            method: 'POST',
            headers: {
              'authorization': 'Bearer '+this.state.auth1,
              'content-Type': 'application/json'
              
            },
            body: '{ "channelName": "'+netConfig.channelName+'", "channelConfigPath":"../artifacts/channel/mychannel.tx" }'
                
                
           
          }
         
      
          fetch('http://localhost:4000' + '/channels', config)
            .then(response =>  response.json() )
            .then((response) => {
                if (response.success === true){
                    return true;
                }

    })

    }*/

    
    craeteChannel (){
        //crearte channel
        let config = {
            method: 'POST',
            headers: {
              'authorization': 'Bearer '+this.state.auth1,
              'content-Type': 'application/json'
              
            },
            body: '{ "channelName": "'+netConfig.channelName+'", "channelConfigPath":"../artifacts/channel/mychannel.tx" }'
                
                
           
          }
         
      
          fetch('http://'+netConfig.hostIP+':'+netConfig.port+'' + '/channels', config)
            .then(response =>  response.json() )
            .then((response) => {
                if (response.success === true){
                    //join channel
                    let config = {
                        method: 'POST',
                        headers: {
                          'authorization': 'Bearer '+this.state.auth1,
                          'content-Type': 'application/json'
                          
                        },
                        body: '{ "peers": ["peer0.org1.example.com","peer1.org1.example.com"] }'
                      }
                     
                  
                      fetch('http://'+netConfig.hostIP+':'+netConfig.port+'' + '/channels/'+netConfig.channelName+'/peers', config)
                        .then(response =>  response.json() )
                        .then((response) => {
                            if (response.success === true){
                                //Join Channel org 2
                                let config = {
                                    method: 'POST',
                                    headers: {
                                      'authorization': 'Bearer '+this.state.auth2,
                                      'content-Type': 'application/json'
                                      
                                    },
                                    body: '{ "peers": ["peer0.org2.example.com","peer1.org2.example.com"] }'
                                  }
                                 
                              
                                  fetch('http://'+netConfig.hostIP+':'+netConfig.port+'' + '/channels/'+netConfig.channelName+'/peers', config)
                                    .then(response =>  response.json() )
                                    .then((response) => {
                                        
                                        if (response.success === true){
                                            let config = {
                                                method: 'POST',
                                                headers: {
                                                  'authorization': 'Bearer '+this.state.auth2,
                                                  'content-Type': 'application/json'
                                                  
                                                },
                                                body: '{ "peers": ["peer0.org2.example.com","peer1.org2.example.com"], "chaincodeName":"'+netConfig.chaincodeName+'", "chaincodePath":"github.com/pcr/blockchain/payerchart/go","chaincodeType": "golang","chaincodeVersion":"v1" }'
                                              }
                                             
                                          //install chaincode org1
                                              fetch('http://'+netConfig.hostIP+':'+netConfig.port+'' + '/chaincodes', config)
                                                .then(response =>  response.json() )
                                                .then((response) => {
                                                    if (response.success === true){
                                                        let config = {
                                                            method: 'POST',
                                                            headers: {
                                                              'authorization': 'Bearer '+this.state.auth1,
                                                              'content-Type': 'application/json'
                                                              
                                                            },
                                                            body: '{ "peers": ["peer0.org1.example.com","peer1.org1.example.com"], "chaincodeName":"'+netConfig.chaincodeName+'", "chaincodePath":"github.com/pcr/blockchain/payerchart/go","chaincodeType": "golang","chaincodeVersion":"v1" }'
                                              }
                                                         
                                                            fetch('http://'+netConfig.hostIP+':'+netConfig.port+'' + '/chaincodes', config)
                                                            .then(response =>  response.json() )
                                                            .then((response) => {
                                                                if (response.success === true){
                                                                    let config = {
                                                                        method: 'POST',
                                                                        headers: {
                                                                          'authorization': 'Bearer '+this.state.auth1,
                                                                          'content-Type': 'application/json'
                                                                          
                                                                        },
                                                                        body: '{"chaincodeName": "'+netConfig.chaincodeName+'", "chaincodeVersion":"v1", "chaincodeType": "golang", "args":[""] }'
                                                                            

                                                                            
                                                                           
                        
                                                                      }
                                                                     
                                                                  
                                                                      fetch('http://'+netConfig.hostIP+':'+netConfig.port+'' + '/channels/'+netConfig.channelName+'/chaincodes', config)
                                                                        .then(response =>  response.json() )
                                                                        .then((response) => {
                                                                            if (response.success === true){
                                                                                this.setState({message: "Successfully instantiate chaingcode in organization Org1 to the channel 'mychannel'"})
                                                                                this.setState({ message: response.message})
                                                                                
                                                                                
                                                            
                                                                            }
                                                                        
                                                                        });
                                                                    
                                                                    
                                                
                                                                }
                                                            });
                                                        
                                                        
                                    
                                                    }
                                                });
                                            
                                            
                        
                                        }
                                    });

                                
            
                            }
                        });
            
                        
                          
                    }

                }
            );
          
    }
    instantiate (){
        let config = {
            method: 'POST',
            headers: {
              'authorization': 'Bearer '+this.state.auth1,
              'content-Type': 'application/json'
              
            },
            body: '{"chaincodeName": "'+netConfig.chaincodeName+'", "chaincodeVersion":"v1", "chaincodeType": "golang", "args":[] }'
                

                
               

          }
         
      
          fetch('http://localhost:4000' + '/chaincodes', config)
            .then(response =>  response.json() )
            .then((response) => {
                if (response.success === true){
                    this.setState({message: response.message})
                    
                    

                }
            });
    }
    /*
    joinchannel() {
        let config = {
            method: 'POST',
            headers: {
              'authorization': 'Bearer '+this.state.auth1,
              'content-Type': 'application/json'
              
            },
            body: {
                "peers": ["peer0.org1.example.com","peer1.org1.example.com"]
            }
          }
         
      
          fetch('http://localhost:4000' + '/channels/'+netConfig.channelName+'/peers', config)
            .then(response =>  response.json() )
            .then((response) => {
                if (response.success === true){
                    //Join Channel org 2
                    let config = {
                        method: 'POST',
                        headers: {
                          'authorization': 'Bearer '+this.state.auth2,
                          'content-Type': 'application/json'
                          
                        },
                        body: {
                            "peers": ["peer0.org2.example.com","peer1.org2.example.com"]
                        }
                      }
                     
                  
                      fetch('http://localhost:4000' + '/channels/'+netConfig.channelName+'/peers', config)
                        .then(response =>  response.json() )
                        .then((response) => {
                            if (response.success === true){
                                installChaincode()
                                
                                
            
                            }
                        });
                    

                }
            });

    }
    installChaincode(){
        let config = {
            method: 'POST',
            headers: {
              'authorization': 'Bearer '+this.state.auth2,
              'content-Type': 'application/json'
              
            },
            body: {
                "peers": ["peer0.org2.example.com","peer1.org2.example.com"],
                "chaincodeName": netConfig.chaincodeName,
                "chaincodePath":"github.com/chaincode",
                "chaincodeType": "golang",
                "chaincodeVersion":"v1"
            }
          }
         
      
          fetch('http://localhost:4000' + '/chaincodes', config)
            .then(response =>  response.json() )
            .then((response) => {
                if (response.success === true){
                    let config = {
                        method: 'POST',
                        headers: {
                          'authorization': 'Bearer '+this.state.auth1,
                          'content-Type': 'application/json'
                          
                        },
                        body: {
                            "peers": ["peer0.org1.example.com","peer1.org1.example.com"],
                            "chaincodeName": netConfig.chaincodeName,
                            "chaincodePath":"github.com/chaincode",
                            "chaincodeType": "golang",
                            "chaincodeVersion":"v1"
                        }
                      }
                     
                  
                      fetch('http://localhost:4000' + '/chaincodes', config)
                        .then(response =>  response.json() )
                        .then((response) => {
                            if (response.success === true){
                                instantiateChaincode()
                                
                                
            
                            }
                        });
                    
                    

                }
            });

    }
    instantiateChaincode(){
        let config = {
            method: 'POST',
            headers: {
              'authorization': 'Bearer '+this.state.auth1,
              'content-Type': 'application/json'
              
            },
            body: {
                "chaincodeName": netConfig.chaincodeName,
                "chaincodeVersion":"v1",
                "chaincodeType": "golang",
                "args":[""]
            }
          }
         
      
          fetch('http://localhost:4000' + '/chaincodes', config)
            .then(response =>  response.json() )
            .then((response) => {
                if (response.success === true){
                    this.setState({message: response.message})
                    
                    

                }
            });

    }*/

    render(){
        return(
            <div>
                     <Button color="success" size="lg" onClick={() => { this.registerUser1()}} text="CreateChannel" variant="action" style={{ margin: '5px'}} />
           <ul>
            AUTH 1: 
            {this.state.auth1}
            </ul>
            <ul>
            AUTH 2:
            {this.state.auth2}
            </ul>
            
            <ul>
            Message: 
            {this.state.message}
            </ul>  
            </div>
        );
    }
}

export default Channel;