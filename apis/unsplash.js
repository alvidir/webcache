const Unsplash = require('unsplash-js').default;
const unsplash = new Unsplash({ accessKey: process.env.API_KEY });

function HandleDefaultRequest(query) {
    return query[0];
}

function HandleRollRequest(query) {
    return query[0];
}

function HandleSupplyRequest(query) {
    return query[0];
}

module.exports.HandleDefaultRequest = HandleDefaultRequest;
module.exports.HandleRollRequest = HandleRollRequest;
module.exports.HandleSupplyRequest = HandleSupplyRequest;