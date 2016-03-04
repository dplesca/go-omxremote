var Videos = React.createClass({

	componentDidMount: function() {
		var that = this;
		AjaxRequest("/files", "GET", function(response) {
			that.setState({
				"Files" : JSON.parse(response),
				"SearchText": that.state.SearchText
			});
		});
	},

	getInitialState : function() {
		return {
			"Files": [],
			"SearchText": ""
		}
	},

	handleUserInput: function(filterText) {
        this.setState({
            "Files": this.state.Files,
            "SearchText": filterText
        });
    },

	render: function() {
		var results = this.state.Files;
		return (
			<div className="row">
				<h3>omxremote</h3>
				<SearchBar searchText={this.state.SearchText} onUserInput={this.handleUserInput} />
				{results.map(function(result) {
					if ( fuzzy.test(this.SearchText.toLowerCase(), result.file.toLowerCase()) ){
					    return (
			        		<Video hash={result.hash} file={result.file}/>
		          		);
					}
		        }, this.state)}
			</div>	
		);
	}

});


var AjaxRequest = function(url, method, callback) {
	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (this.readyState === 4 && this.status == 200) {
			if (callback !== undefined)
				callback(this.responseText)
		}
	}
	xmlhttp.open(method, url, true);
	xmlhttp.send()
}
