var SearchBar = React.createClass({
	handleChange: function(event) {
        this.props.onUserInput(
        	event.target.value
        );
    },

	render: function() {
		return (
			<div className="search-container">
				<input type="text" placeholder="search" className="u-full-width" value={this.props.searchText} ref="filterTextInput" onChange={this.handleChange} />
			</div>
		);
	}

});
