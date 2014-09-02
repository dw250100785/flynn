//= require ../views/main
//= require ../views/login

(function () {

"use strict";

Dashboard.routers.main = new (Marbles.Router.createClass({
	displayName: "routers.main",

	routes: [
		{ path: "", handler: "root" },
		{ path: "login", handler: "login", auth: false },
	],

	root: function () {
		React.renderComponent(
			Dashboard.Views.Main({
					githubAuthed: !!Dashboard.githubClient
				}), Dashboard.el);
	},

	login: function (params) {
		var redirectPath = params[0].redirect || null;
		if (redirectPath && redirectPath.indexOf("//") !== -1) {
			redirectPath = null;
		}
		if ( !redirectPath ) {
			redirectPath = "";
		}

		var performRedirect = function () {
			React.unmountComponentAtNode(Dashboard.config.containerEl);
			Marbles.history.navigate(decodeURIComponent(redirectPath));
		};

		if (Dashboard.config.authenticated) {
			performRedirect();
			return;
		}

		React.renderComponent(
			Dashboard.Views.Login({
					onSuccess: performRedirect
				}), Dashboard.el);
	}

}))();

})();
