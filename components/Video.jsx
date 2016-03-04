var Video = React.createClass({

	getInitialState: function() {
        return { showResults: false };
    },

    onClick: function(e) {
    	e.preventDefault();
    	if (this.state.showResults){
    		this.setState({ showResults: false });
    	} else {
    		this.setState({ showResults: true });
    	}
    },

	render: function() {
		var controls;

		if (this.state.showResults){
			controls = <Controls hash={this.props.hash} />
		} else {
			controls = null;
		}

		return (
			<div key={this.props.hash} className="video">
        		<h5><a href="#" onClick={this.onClick}>{this.props.file}</a></h5>
        		{controls}
    		</div>
		);
	}

});
