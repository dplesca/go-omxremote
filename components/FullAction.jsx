var FullAction = React.createClass({
	handleClick: function(e) {
    	e.preventDefault();
    	AjaxRequest(e.target.href, "POST");
    },

	render: function() {
		var link = "/file/" + this.props.Hash + "/" + this.props.URL,
			buttonClass = "button u-full-width " + this.props.URL,
			iconClass= "icon-" + this.props.URL;
		return (
			<div className="row">
				<a className={buttonClass} onClick={this.handleClick} href={link}><i className={iconClass}></i>{this.props.Action}</a>
			</div>
		);
	}

});
