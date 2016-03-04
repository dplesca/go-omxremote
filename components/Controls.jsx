var Controls = React.createClass({
    handleClick: function(e) {
        e.preventDefault();
        AjaxRequest(e.target.href, "POST");
    },

	render: function() {
        var backLink = "/file/" + this.props.hash + "/backward",
            forwardLink = "/file/" + this.props.hash + "/forward",
            pauseLink = "/file/" + this.props.hash + "/pause";
		return (
			<div className="controls">
    			<FullAction Action="Start" URL="start" Hash={this.props.hash}/>
    			<div className="row">
    				<div className="one-third column">
    					<a onClick={this.handleClick} className="button u-full-width back" href={backLink}><i className="icon-back"></i> Back</a>
    				</div>
    				<div className="one-third column">
    					<a onClick={this.handleClick} className="button u-full-width pause" href={pauseLink}><i className="icon-pause"></i> Pause</a>
    				</div>
    				<div className="one-third column">
    					<a onClick={this.handleClick} className="button u-full-width forward" href={forwardLink}>Forward <i className="icon-forward"></i></a>
    				</div>
    			</div>
    			<FullAction Action="Toggle Subs" URL="subs" Hash={this.props.hash}/>
    			<FullAction Action="Stop" URL="stop" Hash={this.props.hash}/>
    		</div>
		);
	}

});
