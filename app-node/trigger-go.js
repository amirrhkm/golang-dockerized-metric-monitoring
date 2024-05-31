const axios = require("axios");

function callApi(a, b, c) {
  axios
    .get(`http://localhost:8008/service/${a}`)
    .then((response) => {
      console.log(`Request for ${a} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${a}: ${error.message}`);
    });

  axios
    .get(`http://localhost:8008/service/${b}`)
    .then((response) => {
      console.log(`Request for ${b} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${b}: ${error.message}`);
    });

  axios
    .get(`http://localhost:8008/service/${c}`)
    .then((response) => {
      console.log(`Request for ${c} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${c}: ${error.message}`);
    });
}

let index = 0;

callApi(1, 0, 1);
