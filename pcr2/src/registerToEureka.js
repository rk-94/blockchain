const Eureka = require('eureka-js-client').Eureka;

const serviceId = 'blockchain';
const host = 'http://RDCTSTBC003.northamerica.cerner.net:3000';

const client = new Eureka({
    // application instance information
    instance: {
      app: serviceId,
      hostName: host,
      ipAddr: '10.190.148.215',
      port: {
          '$':3000,
          '@enabled':'true',
      },
      vipAddress: serviceId,
      secureVipAddress: serviceId,
      virtualHostName: serviceId,      
      homePageUrl: `${host}/`,
      statusPageUrl: `${host}/info`,
      healthCheckUrl: `${host}/health`,      
      dataCenterInfo: {
        '@class': 'com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo',
        name: 'MyOwn'
      },
      registerWithEureka: true,
      fetchRegistry: true
    },
    eureka: {
        // eureka server host / port
	host: '10.190.150.234',
        port: 8761,
        servicePath: '/eureka/apps/'
      }
  });
  function connectToEureka() {              
    client.logger.level('debug');  
    client.stop(function(error) {
    console.log('########################################################');
    console.log(JSON.stringify(error) || 'Eureka registration complete');   });
    }
    connectToEureka();
