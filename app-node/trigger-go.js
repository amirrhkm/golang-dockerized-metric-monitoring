const axios = require("axios");

function getRandomValue() {
  return Math.floor(Math.random() * 2); // Generate random number between 0 and 1
}

function callApi() {
  const a = getRandomValue();
  const b = getRandomValue();
  const c = getRandomValue();

  axios
    .get(`http://localhost:8008/a/${a}`)
    .then((response) => {
      console.log(`Request for ${a} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${a}: ${error.message}`);
    });

  axios
    .get(`http://localhost:8008/b/${b}`)
    .then((response) => {
      console.log(`Request for ${b} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${b}: ${error.message}`);
    });

  axios
    .get(`http://localhost:8008/c/${c}`)
    .then((response) => {
      console.log(`Request for ${c} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${c}: ${error.message}`);
    });
}

callApi();
