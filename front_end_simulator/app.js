var express = require("express");
var app = express();
var path = require("path");

app.use(express.static(__dirname));

//app.use(express.static(__dirname + '/public'));
//Store all HTML files in view folder.
//app.use(express.static(__dirname + '/public/css'));
//Store all JS and CSS in Scripts folder.

// Add headers
app.use(function (req, res, next) {

    // Website you wish to allow to connect
    res.setHeader('Access-Control-Allow-Origin', 'http://localhost:18842');

    // Request methods you wish to allow
    res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS, PUT, PATCH, DELETE');

    // Request headers you wish to allow
    res.setHeader('Access-Control-Allow-Headers', 'X-Requested-With,content-type');

    // Set to true if you need the website to include cookies in the requests sent
    // to the API (e.g. in case you use sessions)
    res.setHeader('Access-Control-Allow-Credentials', true);

    // Pass to next layer of middleware
    next();
});


app.get('/',function(req,res){
res.sendFile(path.join(__dirname+'/keyboardInputExample.html'));
  //It will find and locate index.html from View or Scripts
});

app.listen(9001);

console.log("Running at Port 9001");