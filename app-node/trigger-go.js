const axios = require("axios");

const pos = [
  "a","b","c","d","e","f"
];

const hub = [
  "1","2","3","4","5","6"
];

const cloud = [
  "z","y","x","w","v","u"
];

function callApi(pos, hub, cloud) {
  axios
    .get(`http://localhost:8008/pos/${pos}`)
    .then((response) => {
      console.log(`Request for ${pos} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${pos}: ${error.message}`);
    });

  axios
    .get(`http://localhost:8008/hub/${hub}`)
    .then((response) => {
      console.log(`Request for ${hub} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${hub}: ${error.message}`);
    });

  axios
    .get(`http://localhost:8008/cloud/${cloud}`)
    .then((response) => {
      console.log(`Request for ${cloud} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${cloud}: ${error.message}`);
    });
}

let index = 0;

setInterval(() => {
  console.log(`--- request interval ---`);
  callApi(pos[index], hub[index], cloud[index]);

  if (index === pos.length - 1) {
    index = 0;
  } else {
    index++;
  }
}, 10);