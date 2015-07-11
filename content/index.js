
var ContextSelector = React.createClass({
    propTypes: {
        url: React.PropTypes.string.isRequired,
        onChoseTimeline: React.PropTypes.func
    },
    getInitialState: function() {
        return {
            contexts: [],
            timelines: [],
            contextURL: "",
            currentContext: "",
            currentTimeline: ""
        }
    },
    componentDidMount: function() {
        // get your data
        $.ajax({
            url: this.props.url,
            success: this.onLoadedContexts
        })
    },
    onLoadedContexts: function(data) {
        var items = data._links.items;
        this.state.contexts= [];

        for (var i = 0; i < items.length; i++) {
            var option = items[i];
            this.state.contexts.push(
                <option key={i} value={option.href}>{option.href}</option>
            );
        }

        this.state.contextURL = data._links.self.href;
        if (items.length > 0) {
            this.state.currentContext = items[0].href;
            this.loadContext(this.state.currentContext);
        } else {
            this.state.currentContext = "";
        }
        this.forceUpdate();
    },
    onChoseContext: function(event) {
        this.state.currentContext = event.target.value;
        this.loadContext(event.target.value);
    },
    loadContext: function(context) {
        var p = urljoin(this.state.contextURL, context);
        $.ajax({
            url: p,
            success: this.onLoadedTimelines
        })
    },
    onLoadedTimelines: function(data) {
        var items = data._links.items;
        this.state.timelines = [];
        for (var i = 0; i < items.length; i++) {
            var option = items[i];
            this.state.timelines.push(
                <option key={i} value={option.href}>{option.href}</option>
            );
        }

        this.forceUpdate();

        if (items.length > 0) {
            this.state.currentTimeline = items[0].href;

            if( this.props.onChoseTimeline !== undefined ) {
                this.props.onChoseTimeline(
                    this.state.currentContext,
                    this.state.currentTimeline
                );
            }
        } else {
            this.state.currentTimeline = "";
        }
    },
    onChooseTimeline: function(event) {
        this.state.currentTimeline = event.target.value;
        if( this.props.onChoseTimeline !== undefined ) {
            this.props.onChoseTimeline(
                this.state.currentContext,
                this.state.currentTimeline
            );
        }
    },
    render: function() {
        return (
            <div>
                <h3>Context</h3>
                <select {...this.props} onChange={this.onChoseContext} >{this.state.contexts}</select>
                <h4>Timeline</h4>
                <select {...this.props} onChange={this.onChooseTimeline} >{this.state.timelines}</select>
            </div>
        );
    }
});
