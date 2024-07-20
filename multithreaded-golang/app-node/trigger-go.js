const axios = require("axios");

const pos = [
  "a","b","c","d","e","f"
];

const hub = [
  "1","2","3","4","5","6"
];

const cds = [
  "z","y","x","w","v","u"
];

function callApi(pos, hub, cds) {
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
    .get(`http://localhost:8008/cds/${cds}`)
    .then((response) => {
      console.log(`Request for ${cds} completed successfully`);
    })
    .catch((error) => {
      console.error(`Error making request for ${cds}: ${error.message}`);
    });
}

let index = 0;

setInterval(() => {
  console.log(`--- request interval ---`);
  callApi(pos[index], hub[index], cds[index]);

  if (index === pos.length - 1) {
    index = 0;
  } else {
    index++;
  }
}, 10);