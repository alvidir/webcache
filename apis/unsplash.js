const Unsplash = require('unsplash-js').default;
const unsplash = new Unsplash({ accessKey: process.env.API_KEY });

function HandleSingleRequest() {
    return 'single one';
}

function HandleRollRequest(query) {
    if (query.size === 0) {
        return '', 'Error: no parameters set for roll call'
    }

    return JSON.stringify(query);
}

function HandleSupplyRequest(query) {
    if (query.size === 0) {
        return 'Error: No parameters set for supply call'
    }

    return query[0];
}

module.exports.HandleSingleRequest = HandleSingleRequest;
module.exports.HandleRollRequest = HandleRollRequest;
module.exports.HandleSupplyRequest = HandleSupplyRequest;