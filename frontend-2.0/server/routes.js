var request = require('request');
var bodyParser = require('body-parser');

module.exports = app => {
	app.use( bodyParser.json() );       // to support JSON-encoded bodies

	app.use(bodyParser.urlencoded({     // to support URL-encoded bodies
	  extended: true
	}));

	app.get('/', (req, res) => {
		res.sendfile('./client/index.html');
	});

};