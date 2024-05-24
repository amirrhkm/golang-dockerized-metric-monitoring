const axios = require("axios");

const names = [
  "a","b","c","d","e","f"
];

const types = [
  "1","2","3","4","5","6"
];

const fuel = [
  "z","y","x","w","v","u"
];

function callApi(name, type, fuel) {
  axios
    .get(`http://localhost:8008/POS/${name}`)
    .then((response) => {
      console.log(`Request for ${name} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${name}: ${error.message}`);
    });

  axios
    .get(`http://localhost:8008/HUB/${type}`)
    .then((response) => {
      console.log(`Request for ${type} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${type}: ${error.message}`);
    });

  axios
    .get(`http://localhost:8008/CLOUD/${fuel}`)
    .then((response) => {
      console.log(`Request for ${fuel} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${fuel}: ${error.message}`);
    });
}

let index = 0;

setInterval(() => {
  console.log(`${index}  ${names[index]}`);
  callApi(names[index], types[index], fuel[index]);

  if (index === names.length - 1) {
    index = 0;
  } else {
    index++;
  }
}, 10);